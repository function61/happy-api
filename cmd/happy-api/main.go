package main

import (
	"embed"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/function61/gokit/app/aws/lambdautils"
	"github.com/function61/gokit/app/dynversion"
	"github.com/function61/gokit/crypto/cryptoutil"
	"github.com/function61/gokit/log/logex"
	"github.com/function61/gokit/net/http/httputils"
	"github.com/function61/gokit/os/osutil"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

//go:embed item.html
var templates embed.FS

//go:embed images/*
var images embed.FS

func main() {
	rand.Seed(time.Now().UnixNano())

	// AWS Lambda doesn't support giving argv, so we use an ugly hack to detect when
	// we're in Lambda
	if lambdautils.InLambda() {
		lambda.StartHandler(lambdautils.NewLambdaHttpHandlerAdapter(httpHandler()))
		return
	}

	app := &cobra.Command{
		Use:     os.Args[0],
		Short:   "Happiness as a service",
		Version: dynversion.Version,
		Args:    cobra.NoArgs,
		Run: func(_ *cobra.Command, args []string) {
			srv := &http.Server{
				Addr:    ":80",
				Handler: httpHandler(),
			}

			osutil.ExitIfError(httputils.CancelableServer(
				osutil.CancelOnInterruptOrTerminate(logex.StandardLogger()),
				srv,
				func() error { return srv.ListenAndServe() }))

		},
	}

	app.AddCommand(&cobra.Command{
		Use:   "new",
		Short: "Generate ID for new file",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, args []string) {
			fmt.Println(cryptoutil.RandBase64UrlWithoutLeadingDash(3))
		},
	})

	osutil.ExitIfError(app.Execute())
}

func httpHandler() http.Handler {
	uiTpl, err := template.ParseFS(templates, "item.html")
	if err != nil {
		panic(err)
	}

	happiness, err := images.ReadDir("images")
	if err != nil {
		panic(err)
	}

	routes := mux.NewRouter()

	redirectToRandomItem := func(w http.ResponseWriter, r *http.Request) {
		idx := randBetween(0, len(happiness)-1)

		http.Redirect(w, r, "/happy/"+fileIdFromFilename(happiness[idx].Name()), http.StatusFound)
	}

	routes.PathPrefix("/happy/images/").Handler(http.StripPrefix("/happy/", http.FileServer(http.FS(images))))

	routes.HandleFunc("/happy", redirectToRandomItem)
	routes.HandleFunc("/happy/", redirectToRandomItem)

	routes.HandleFunc("/happy/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		attribution, err := findAttributionFromExifArtist(id)
		if err != nil { // assuming error is ErrNotExist
			if os.IsNotExist(err) {
				http.NotFound(w, r)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "text/html")

		_ = uiTpl.Execute(w, struct {
			ImgSrc      string
			Attribution string
		}{
			ImgSrc:      "/happy/images/" + id + ".jpg",
			Attribution: attribution,
		})
	})

	return routes
}
