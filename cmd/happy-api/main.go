package main

import (
	"embed"
	"html/template"
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

//go:embed item.html
var templates embed.FS

func main() {
	rand.Seed(time.Now().UnixNano())

	// AWS Lambda doesn't support giving argv, so we use an ugly hack to detect when
	// we're in Lambda
	if lambdautils.InLambda() {
		lambda.StartHandler(lambdautils.NewLambdaHttpHandlerAdapter(httpHandler()))
		return
	}

	srv := &http.Server{
		Addr:    ":80",
		Handler: httpHandler(),
	}

	osutil.ExitIfError(httputils.CancelableServer(
		osutil.CancelOnInterruptOrTerminate(logex.StandardLogger()),
		srv,
		func() error { return srv.ListenAndServe() }))
}

func httpHandler() http.Handler {
	uiTpl, err := template.ParseFS(templates, "item.html")
	if err != nil {
		panic(err)
	}

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
			ImgSrc:      "https://s3.amazonaws.com/onni.function61.com/media/" + id,
			Attribution: record.Attribution,
		})
	})

	return routes
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
