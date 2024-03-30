package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	cfg, err := config.LoadDefaultConfig(ctx);

	if err  != nil {
		log.Fatal(err)
	}
	
	ddb := dynamodb.NewFromConfig(cfg)

	tableName := "snakes-and-ladders"
	key := map[string]types.AttributeValue{
		"roomid": &types.AttributeValueMemberS{
			Value: "0000",
		},
	}

	res, err := ddb.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: key,
	})

	if err != nil {
		log.Fatal(err)
	}

	p1ConnectionID := res.Item["player0"].(*types.AttributeValueMemberS).Value
	p2ConnectionID := res.Item["player1"].(*types.AttributeValueMemberS).Value
	
	positions := res.Item["positions"].(*types.AttributeValueMemberL).Value

	p1PosS := positions[0].(*types.AttributeValueMemberN).Value
	p2PosS := positions[1].(*types.AttributeValueMemberN).Value

	p1Pos, err := strconv.Atoi(p1PosS)

	if err != nil {
		log.Fatal(err)
	}

	p2Pos, err := strconv.Atoi(p2PosS)

	if err != nil {
		log.Fatal(err)
	}

	turn := res.Item["turn"].(*types.AttributeValueMemberBOOL).Value

	lastMessages := res.Item["lastMessages"].(*types.AttributeValueMemberL).Value

	newLastMessages := []string{}

	for _, lastMessage := range lastMessages {
		msg := lastMessage.(*types.AttributeValueMemberS).Value
		newLastMessages = append(newLastMessages, msg)
	}

	payloadString := map[string]interface{} {
		"player0": p1ConnectionID,
		"player1": p2ConnectionID,
		"turn": turn,
		"positions": []int{p1Pos, p2Pos},
		"lastMessages": newLastMessages,
	}

	jbytes, err := json.Marshal(payloadString)

	if err != nil {
		log.Fatal(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body: string(jbytes),
	}, nil

}
