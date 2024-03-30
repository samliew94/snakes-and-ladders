package main

import (
	"context"
	"log"

	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	lambdaService "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang-jwt/jwt/v5"
	"github.com/valyala/fastjson"
)

func main() {
	lambda.Start(handler)
}

// reply the websocket connectionId with whatever was sent
func handler(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	cfg, err := config.LoadDefaultConfig(ctx);

	if err  != nil {
		log.Fatal(err)
	}

	connectionID := request.RequestContext.ConnectionID

	lbd := lambdaService.NewFromConfig(cfg)

	res, err := lbd.Invoke(ctx, &lambdaService.InvokeInput{
		FunctionName: aws.String("snl-board-data"),
	})

	if err != nil {
		log.Fatal(err)
	}

	body := PayloadToBody(res.Payload)

	// determine if connectionID is p1 or p2
	p1ConnectionID := string(body.GetStringBytes("player0"))
	p2ConnectionID := string(body.GetStringBytes("player1"))
	players := []string{""} // DDB Set of String must have atleast one value, so we default to ""

	player := -1

	if connectionID != "" {

		if connectionID == p1ConnectionID {
			player = 0
		} else if connectionID == p2ConnectionID {
			player = 1 
		} else {
			log.Fatal("connectionID is neither p1 nor p2")
		}

	} else {

		log.Fatal("connectionID is empty")

	}

	if player == 0 {
		p1ConnectionID = "" // remove p1 ConnID
	} else if player == 1 {
		p2ConnectionID = ""
	}

	if p1ConnectionID != "" {
		players = append(players, p1ConnectionID)
	}

	if p2ConnectionID != "" {
		players = append(players, p2ConnectionID)
	}

	ddb := dynamodb.NewFromConfig(cfg)

	tableName := aws.String("snakes-and-ladders")
	key := map[string]types.AttributeValue {
		"roomid": &types.AttributeValueMemberS {
			Value: "0000",
		},
	}

	_, err = ddb.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: tableName,
		Key: key,
		UpdateExpression: aws.String("set player0 = :player0, player1 = :player1, players = :players"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":player0": &types.AttributeValueMemberS{Value: p1ConnectionID},
			":player1": &types.AttributeValueMemberS{Value: p2ConnectionID},
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


func DecodeJwt(request *events.APIGatewayWebsocketProxyRequest) (jwt.MapClaims) {

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
		log.Fatal("fatal error. failed to parse jwt")
	}

	claims, ok := decoded.Claims.(jwt.MapClaims)

	if !ok {
		log.Fatal("fatal error. failed to cast claims to jwt.MapClaims")
	}

	return claims

}

func PayloadToBody(payload []byte) (*fastjson.Value) {
	
	body := fastjson.GetBytes(payload, "body")

	parsed, err := fastjson.ParseBytes(body)

	if err != nil {
		log.Fatal(err)
	}

	return parsed
	
}