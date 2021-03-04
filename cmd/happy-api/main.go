package main

import (
	"context"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/function61/gokit/app/aws/lambdautils"
	"github.com/function61/gokit/log/logex"
	"github.com/function61/gokit/net/http/httputils"
	"github.com/function61/gokit/os/osutil"
	"github.com/gorilla/mux"
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
	<a href="/api/happy/">Show me another</a>
</div>

</body>
</html>
`)

func main() {
	rand.Seed(time.Now().UnixNano())

	// AWS Lambda doesn't support giving argv, so we use an ugly hack to detect when
	// we're in Lambda
	if lambdautils.InLambda() {
		lambda.StartHandler(lambdautils.NewLambdaHttpHandlerAdapter(httpHandler()))
		return
	}

	rootLogger := logex.StandardLogger()

	osutil.ExitIfError(runStandaloneRestApi(
		osutil.CancelOnInterruptOrTerminate(rootLogger),
		rootLogger))
}

func httpHandler() http.Handler {
	routes := mux.NewRouter()

	redirectToRandomItem := func(w http.ResponseWriter, r *http.Request) {
		idx := randBetween(0, len(happiness)-1)

		http.Redirect(w, r, "/api/happy/happiness/"+happiness[idx].Id, http.StatusFound)
	}

	routes.HandleFunc("/api/happy", redirectToRandomItem)
	routes.HandleFunc("/api/happy/", redirectToRandomItem)

	routes.HandleFunc("/api/happy/happiness/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		record := findRecord(id)
		if record == nil {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html")

		_ = uiTpl.Execute(w, struct {
			ImgSrc      string
			Attribution string
		}{
			ImgSrc:      makeMediaUrl(id),
			Attribution: record.Attribution,
		})
	})

	return routes
}

// for standalone use
func runStandaloneRestApi(ctx context.Context, logger *log.Logger) error {
	srv := &http.Server{
		Addr:    ":80",
		Handler: httpHandler(),
	}

	return httputils.CancelableServer(ctx, srv, func() error { return srv.ListenAndServe() })
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

func makeMediaUrl(id string) string {
	return "https://s3.amazonaws.com/onni.function61.com/media/" + id
}
