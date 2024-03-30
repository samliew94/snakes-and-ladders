package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	claims := DecodeJwt(request)

	player, ok := claims["player"].(float64)

	if !ok {
		log.Fatal("fatal error. claims['player'] is not a number")
	}

	if (!(player == 0 || player == 1)) {
		log.Fatal("fatal error. claims['player'] must be either 0 or 1")
	}

	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		log.Fatal(err)
	}

	ddb := dynamodb.NewFromConfig(cfg)

	tableName := aws.String("snakes-and-ladders")
	key := map[string]types.AttributeValue{
		"roomid": &types.AttributeValueMemberS{Value: "0000"},
	}

	// update the current playerX connId (only 1 for player1 and player2 respectively)
	res, err := ddb.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: tableName,
		Key: key,
	})

	if err != nil {
		log.Fatal(err)
	}

	// get existing player1 and 2 values from DDB
	player1 := res.Item["player0"].(*types.AttributeValueMemberS).Value
	player2 := res.Item["player1"].(*types.AttributeValueMemberS).Value

	// overwrite as needed
	if (player == 0) {
		player1 = request.RequestContext.ConnectionID
	} else if player == 1 {
		player2 = request.RequestContext.ConnectionID
	}

	// get all existing connections
	players := res.Item["players"].(*types.AttributeValueMemberSS).Value

	invalidConnectionIds := []string{}

	for _, connectionId := range players {

		if !(player1 == connectionId || player2 == connectionId) {
			invalidConnectionIds = append(invalidConnectionIds, connectionId)
		}

	}

	// disconnect the invalid connection ids (if err is ok)
	api := apigatewaymanagementapi.NewFromConfig(cfg, func(o *apigatewaymanagementapi.Options) {
		o.BaseEndpoint = aws.String("https://6mr0m656c6.execute-api.ap-southeast-1.amazonaws.com/dev")
	})

	for _, connectionId := range invalidConnectionIds {
		_, err := api.DeleteConnection(ctx, &apigatewaymanagementapi.DeleteConnectionInput{
			ConnectionId: aws.String(connectionId),
		})

		if err != nil {
			log.Printf(`Error when attempting to delete connectionId %v %v`, connectionId, err)
		} else {
			log.Printf(`DeletedConnection: %v`, connectionId)
		}
	}

	// reset the Set `players`, populated with the player1 and player2's connectionId`
	players = []string{player1, player2}

	// update player1|player2 and players
	_, err = ddb.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: tableName,
		Key: key,
		UpdateExpression: aws.String("set player0 = :player0, player1 = :player1, players = :players"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":player0": &types.AttributeValueMemberS{Value: player1},
			":player1": &types.AttributeValueMemberS{Value: player2},
			":players": &types.AttributeValueMemberSS{Value: players},
		},

	})

	if err != nil {
		log.Fatal(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil

}

func DecodeJwt(request events.APIGatewayWebsocketProxyRequest) (jwt.MapClaims) {

	authorizer, ok := request.RequestContext.Authorizer.(map[string]interface{})

	if !ok {
		log.Fatal("fatal error. failed to cast authorizer to map[string]interface")
	}

	token, ok := authorizer["token"].(string)

	if !ok {
		log.Fatal("fatal error. authorizer property 'token' is not a string")
	}

	decoded, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		log.Fatalf("failed to parse jwt: %v", err)
	}

	claims, ok := decoded.Claims.(jwt.MapClaims)

	if !ok {
		log.Fatal("fatal error. failed to cast claims to jwt.MapClaims")
	}

	return claims

}