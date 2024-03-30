package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayCustomAuthorizerResponse, error) {

	token := request.QueryStringParameters["token"]

	if token == "" {
		log.Fatal("token is empty")
	}

	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		log.Fatalf("jwt parse failed. %v", err)
	}

	return events.APIGatewayCustomAuthorizerResponse{
		// PrincipalID: "foo-my-bar",
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version:   "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   "Allow",
					Resource: []string{"arn:aws:execute-api:ap-southeast-1:168591133936:6mr0m656c6/*/$connect"},
				},
			},			
		},
		Context:map[string]interface{}{
			"token": token,
		},
	}, nil

}


