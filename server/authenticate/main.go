package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt/v5"
	"github.com/valyala/fastjson"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	body, err := fastjson.Parse(request.Body)

	if err != nil {
		log.Fatalf("error parsing %v" ,err)
	}

	player := body.GetFloat64("player")

	// player must be either 0 or 1
	if !(player == 0 || player == 1) {
		log.Fatal("property 'player' must be either 0 or 1")
	}

	// generate jwt
	key := []byte(os.Getenv("SECRET"))
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "AWS_LAMBDA_snl_authenticate",
		"player": player,
	})
	
	token, err := tokenObj.SignedString(key)

	if err != nil {
		log.Fatalf("error signing %v" ,err)
	}

	jbytes, err := json.Marshal(map[string]string{"token": token})

	if err != nil {
		log.Fatalf("error marshalling %v" ,err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		// Headers: map[string]string{
		// 	"Access-Control-Allow-Headers": "content-type,x-amz-date,authorization,x-api-key,x-amz-security-token",
		// 	"Access-Control-Allow-Methods": "*",
		// 	"Access-Control-Allow-Origin": "*",
		// },
		Body: string(jbytes),
	}, nil

}

