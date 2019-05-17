package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"math/rand"
	"net/http"
	"time"
)

var happiness = []string{
	"https://twitter.com/respros/status/1121496846042636289",
	"http://www.awesomelycute.com/2015/04/25-of-the-cutest-kittens-ever/",
}

func onniHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	idx := randBetween(0, len(happiness)-1)

	return redirect(happiness[idx]), nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(onniHandler)
}

func randBetween(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func redirect(to string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusFound,
		Headers: map[string]string{
			"Location": to,
		},
		Body: fmt.Sprintf("Redirecting to %s", to),
	}
}
