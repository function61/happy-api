package main

import (
	"embed"
	"html/template"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/function61/gokit/app/aws/lambdautils"
	"github.com/function61/gokit/app/dynversion"
	"github.com/function61/gokit/log/logex"
	"github.com/function61/gokit/net/http/httputils"
	"github.com/function61/gokit/os/osutil"
	"github.com/function61/happy-api/pkg/turbocharger/turbochargerapp"
	"github.com/function61/happy-api/static"
	"github.com/spf13/cobra"
)

//go:embed item.html
var templates embed.FS

func main() {
	// AWS Lambda doesn't support giving argv, so we use an ugly hack to detect when
	// we're in Lambda
	if lambdautils.InLambda() {
		lambda.Start(lambdautils.NewLambdaHttpHandlerAdapter(httpHandler()))
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

				ReadHeaderTimeout: httputils.DefaultReadHeaderTimeout,
			}

			osutil.ExitIfError(httputils.CancelableServer(
				osutil.CancelOnInterruptOrTerminate(logex.StandardLogger()),
				srv,
				srv.ListenAndServe))

		},
	}

	app.AddCommand(newEntry())
	app.AddCommand(turbochargerapp.StaticFilesExportEntrypoint(static.Files))

	osutil.ExitIfError(app.Execute())
}

func httpHandler() http.Handler {
	uiTpl, err := template.ParseFS(templates, "item.html")
	if err != nil {
		panic(err)
	}

	happiness, err := static.Files.ReadDir("images")
	if err != nil {
		panic(err)
	}

	routes := http.NewServeMux()

	redirectToRandomItem := func(w http.ResponseWriter, r *http.Request) {
		idx := randBetween(0, len(happiness)-1)

		http.Redirect(w, r, "/happy/"+fileIdFromFilename(happiness[idx].Name()), http.StatusFound)
	}

	routes.Handle("/happy/static/", turbochargerapp.FileHandler("/happy/static", static.Files))

	routes.HandleFunc("/happy/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

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
			ImgSrc:      "/happy/static/images/" + id + ".jpg",
			Attribution: attribution,
		})
	})

	routes.HandleFunc("/happy", redirectToRandomItem)
	routes.HandleFunc("/happy/", redirectToRandomItem)

	return routes
}
