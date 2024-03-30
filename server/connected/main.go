package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	lambdaService "github.com/aws/aws-sdk-go-v2/service/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	cfg, err := config.LoadDefaultConfig(ctx);

	if err  != nil {
		log.Fatal(err)
	}

	lbd := lambdaService.NewFromConfig(cfg)

	payloadOne, err := json.Marshal(map[string]interface{}{"ConnectionID": request.RequestContext.ConnectionID})

	if err != nil {
		log.Fatalf("failed to marshal connectionid %v", err)
	}

	payload, err := json.Marshal(map[string]interface{}{"Body": string(payloadOne)})

	if err != nil{
		log.Fatalf(`failed to marshal payload err=%v`, payload)
	}

	_, err = lbd.Invoke(context.TODO(), &lambdaService.InvokeInput{
		FunctionName: aws.String("snl-broadcast-board"),
		Payload: payload,
	})

	if err != nil {
		log.Fatalf(`Error invoking lambda snl-broadcast-board %v`, err)
	}
	
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil

}
