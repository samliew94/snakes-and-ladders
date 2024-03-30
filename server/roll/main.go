package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	lambdaService "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/valyala/fastjson"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request *events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	cfg, err := config.LoadDefaultConfig(ctx);

	if err != nil {
		log.Fatal(err)
	}

	lbd := lambdaService.NewFromConfig(cfg)

	res, err := lbd.Invoke(ctx, &lambdaService.InvokeInput{
		FunctionName: aws.String("snl-board-data"),
	})

	if err != nil {
		log.Fatal(err)
	}

	body := PayloadToBody(res.Payload)

	connectionID := request.RequestContext.ConnectionID

	p1ConnectionID := string(body.GetStringBytes("player0"))
	p2ConnectionID := string(body.GetStringBytes("player1"))

	turn := body.GetBool("turn")

	// verify it's correct player's turn
	if !((turn && connectionID != "" && connectionID == p1ConnectionID) || (!turn && connectionID != "" && connectionID == p2ConnectionID)) {
		log.Fatal("wrong player turn!")
	} 

	turn = !turn

	player := 0

	if connectionID == p1ConnectionID && p1ConnectionID != "" {
		player = 0
	} else if connectionID == p2ConnectionID && p2ConnectionID != "" {
		player = 1
	} else {
		log.Fatalf("ConnectionID %v is neither p1 nor p2", connectionID)
	}
	
	positionsRaw := body.GetArray("positions")
	positions := []int{int(positionsRaw[0].GetFloat64()), int(positionsRaw[1].GetFloat64())}
	dice := rand.Intn(6) + 1

	positions[player] += dice

	newLogs := []string{}

	loc := time.FixedZone("GMT+8", 8*60*60)
	now := time.Now().In(loc)

	newLog := fmt.Sprintf("[%s] Player%d rolled %d. ", now.Format("2006-01-02 15:04"), player+1, dice)

	if positions[player] > 24 {
		diff := positions[player] - 24
		positions[player] = 24 - diff
		newLog += fmt.Sprintf("Went beyond 25. Backed down by %d. ", diff)		
	}

	newLog += fmt.Sprintf("Landed on Tile %d. ", positions[player]+1)

	newLog += LandedOnSnakesOrLadder(positions, player, newLogs)

	newLog += LandedOnWin(positions, player, newLogs)

	lastMessagesRaw := body.GetArray("lastMessages")

	lastMessages := []string{}

	for _, msg := range lastMessagesRaw {
		lastMessages = append(lastMessages, string(msg.GetStringBytes()))
	}

	lastMessages = append(lastMessages, newLog)
	
	for len(lastMessages) > 5 {
		lastMessages = lastMessages[1:]
	}

	ddbLastMessages := &types.AttributeValueMemberL{}

	for _, msg := range lastMessages {
		ddbLastMessages.Value = append(ddbLastMessages.Value, &types.AttributeValueMemberS{
			Value: msg,
		})
	}

	tableName := aws.String("snakes-and-ladders")
	key := map[string]types.AttributeValue {
		"roomid": &types.AttributeValueMemberS {
			Value: "0000",
		},
	}
	
	ddb := dynamodb.NewFromConfig(cfg)

	_, err = ddb.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: tableName,
		Key: key,		
		UpdateExpression: aws.String("SET positions[0] = :p1, positions[1] = :p2, lastMessages = :lastMessages, turn = :turn"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":p1": &types.AttributeValueMemberN{
				Value: strconv.Itoa(positions[0]),
			},
			":p2": &types.AttributeValueMemberN{
				Value: strconv.Itoa(positions[1]),
			},
			":lastMessages": ddbLastMessages,
			":turn": &types.AttributeValueMemberBOOL{
				Value: turn,
			},
		},		
	})

	if err != nil {
		log.Fatal(err)
	}

	_, err = lbd.Invoke(ctx, &lambdaService.InvokeInput{
		FunctionName: aws.String("snl-broadcast-board"),
		Payload: BodyToPayload(map[string]interface{}{"ConnectionID": request.RequestContext.ConnectionID}),
	})

	if err != nil {
		log.Fatal(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil

}

func PayloadToBody(payload []byte) (*fastjson.Value) {
	
	body := fastjson.GetBytes(payload, "body")

	parsed, err := fastjson.ParseBytes(body)

	if err != nil {
		log.Fatal(err)
	}

	return parsed
	
}

func BodyToPayload(input map[string]interface{}) ([]byte) {
	prePayload, err := json.Marshal(input)

	if err != nil {
		log.Fatal(err)
	}

	payload, err := json.Marshal(map[string]interface{}{"Body": string(prePayload)})

	if err != nil {
		log.Fatal(err)
	}

	return payload
}

func LandedOnSnakesOrLadder(positions []int, player int, newLogs []string) (string) {
	
	newLog := ""

	// -1 (nil), 0 (snaked), 1 (laddered)
	snakedOrLaddered := -1

	if positions[player] == 2 {
		// ladder to 8
		positions[player] = 8
		snakedOrLaddered = 1
	} else if positions[player] == 11 {
		// snakes to 5
		positions[player] = 5
		snakedOrLaddered = 0
	} else if positions[player] == 13 {
		// ladder to 20
		positions[player] = 20
		snakedOrLaddered = 1
	} else if positions[player] == 22 {
		// snakes to 1
		positions[player] = 1
		snakedOrLaddered = 0
	}

	if snakedOrLaddered == 0 {
		newLog += fmt.Sprintf("Snaked to Tile %d. ", positions[player]+1)
	} else if snakedOrLaddered == 1 {
		newLog += fmt.Sprintf("Laddered to Tile %d. ", positions[player]+1)
	}

	return newLog

}

func LandedOnWin(positions []int, player int, newLogs []string) (string) {
	
	if positions[player] == 24 {

		newLog := fmt.Sprintf("Player %d Wins!. Resetting board. ", player+1)

		positions[0] = 0
		positions[1] = 0

		return newLog

	}

	return ""

}