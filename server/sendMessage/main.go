package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/aws"
)

func main() {
	lambda.Start(handler)
}

// reply the websocket connectionId with whatever was sent
func handler(ctx context.Context, request *events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	cfg, err := config.LoadDefaultConfig(ctx);

	connectionId := request.RequestContext.ConnectionID
	log.Printf("connectionId %v", connectionId)

	if err  != nil {
		log.Fatal(err)
	}

	api := apigatewaymanagementapi.NewFromConfig(cfg)

	var reqBody map[string]interface{}

	err = json.Unmarshal([]byte(request.Body), &reqBody)

	if err != nil {
		log.Fatal(err)
	}

	delete(reqBody, "action")

	jbytes, err := json.Marshal(reqBody)

	if err != nil {
		log.Fatal(err)
	}

	_, err = api.PostToConnection(ctx, &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(connectionId),
		Data: jbytes,
	}, func(o *apigatewaymanagementapi.Options) {
		o.BaseEndpoint = aws.String("https://6mr0m656c6.execute-api.ap-southeast-1.amazonaws.com/dev")
	})

	if err != nil {
		log.Fatal(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil

}
