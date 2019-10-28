package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var uiTpl, _ = template.New("_").Parse(`
<!doctype html>
<html>
<head>
	<title>Happiness</title>
</head>
<body>

<div>
	<img src="{{.ImgSrc}}" alt="" />

	<hr />

{{if .Attribution}}
	<a href="{{.Attribution}}">Source</a>
{{else}}
	Source not known
{{end}}
</div>

<div>
	<a href="{{.BaseUrl}}/happy">Show me another</a>
</div>

</body>
</html>
`)

type Happiness struct {
	Id          string
	Attribution string
}

var happiness = []Happiness{
	{"10e239c4167f", "https://gizmodo.com/owls-are-weighed-wrapped-up-in-blankets-like-little-bir-1621869419"},
	{"173fbced6b8a", ""},
	{"2c66790f5801", ""},
	{"3799bc32c4c5", ""},
	{"a320b503df35", ""},
	{"d234c2dbc127", ""},
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

			responseBody := &bytes.Buffer{}
			uiTpl.Execute(responseBody, struct {
				ImgSrc      string
				Attribution string
				BaseUrl     string
			}{
				ImgSrc:      makeMediaUrl(id),
				Attribution: record.Attribution,
				BaseUrl:     createBaseUrl(req),
			})

			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Headers: map[string]string{
					"Content-Type": "text/html",
				},
				Body: responseBody.String(),
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
