package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	lambdaService "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go/aws"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Println("[SNL-BROADCAST-BOARD] called...")

	cfg, err := config.LoadDefaultConfig(ctx);

	if err  != nil {
		log.Fatal(err)
	}

	log.Println("request.Body:")
	log.Println(request.Body)

	reqBody := map[string]interface{}{}
	err = json.Unmarshal([]byte(request.Body), &reqBody)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("marhsalled reqBody %v", reqBody)
	}
	
	// get the current data from DDB by invoking snl-board-data
	lbd := lambdaService.NewFromConfig(cfg)
	res, err := lbd.Invoke(ctx, &lambdaService.InvokeInput{
		FunctionName: aws.String("snl-board-data"),		
	})

	if err != nil {
		log.Fatal(err)
	}

	fullRes := make(map[string]interface{})
	err = json.Unmarshal([]byte(res.Payload), &fullRes)

	if err != nil {
		log.Fatal(err)
	}

	rawBody := fullRes["body"].(string)
	body := map[string]interface{}{}
	json.Unmarshal([]byte(rawBody), &body)

	p1ConnectionID := body["player0"].(string)
	p2ConnectionID := body["player1"].(string)
	playersConnectionIDs := []string{p1ConnectionID, p2ConnectionID}

	delete(body, "player0")
	delete(body, "player1")
	
	api := apigatewaymanagementapi.NewFromConfig(cfg)

	// broadcast to all players
	for _, playerConnectionID := range playersConnectionIDs {

		// don't bother sending to empty ConnectionID
		if playerConnectionID == "" {
			continue
		}

		// determine if this ConnectionID is p1 or p2
		if playerConnectionID == p1ConnectionID {
			body["player"] = 0
		} else if playerConnectionID == p2ConnectionID {
			body["player"] = 1
		}

		data, err := json.Marshal(body)

		if err != nil {
			log.Fatal(err)
		}

		_, err = api.PostToConnection(ctx, &apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: aws.String(playerConnectionID),
			Data: data,
		}, func(o *apigatewaymanagementapi.Options) {
			o.BaseEndpoint = aws.String("https://6mr0m656c6.execute-api.ap-southeast-1.amazonaws.com/dev")
		})

		if err != nil {
			log.Println(err)
		}
		
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil

}
