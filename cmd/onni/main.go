package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const template = `
<!doctype html>
<html>
<head>
	<title>Happiness</title>
</head>
<body>

<div>
	<img src="%s" alt="" />

	<hr />

	<a href="%s">Source</a>
</div>

<div>
	<a href="%s/happy">Show me another</a>
</div>

</body>
</html>
`

type Happiness struct {
	Id     string
	Source string
}

var happiness = []Happiness{
	{"10e239c4167f", "https://gizmodo.com/owls-are-weighed-wrapped-up-in-blankets-like-little-bir-1621869419"},
}

func onniHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	synopsis := req.HTTPMethod + " " + req.Path

	switch synopsis {
	case "GET /": // root should redirect to project homepage
		return redirect("https://github.com/function61/onni"), nil
	case "GET /happy": // root should redirect to project homepage
		id := req.QueryStringParameters["id"]

		if id != "" {
			record := findRecord(id)
			if record == nil {
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusNotFound,
					Body:       "file not found",
				}, nil
			}
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Headers: map[string]string{
					"Content-Type": "text/html",
				},
				Body: fmt.Sprintf(
					template,
					makeMediaUrl(id),
					record.Source,
					createBaseUrl(req)),
			}, nil
		} else {
			idx := randBetween(0, len(happiness)-1)

			return redirect(createBaseUrl(req) + "/happy?id=" + happiness[idx].Id), nil
		}
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       fmt.Sprintf("Unknown endpoint: %s", synopsis),
		}, nil
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(onniHandler)
}

func findRecord(id string) *Happiness {
	for _, record := range happiness {
		if record.Id == id {
			return &record
		}
	}

	return nil
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

func makeMediaUrl(id string) string {
	return "https://s3.amazonaws.com/" + os.Getenv("S3_BUCKET") + "/media/" + id
}

func createBaseUrl(req events.APIGatewayProxyRequest) string {
	return "https://" + req.RequestContext.APIID + ".execute-api.us-east-1.amazonaws.com/prod"
}
