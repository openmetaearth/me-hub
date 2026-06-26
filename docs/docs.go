package docs

import (
	"embed"
	httptemplate "html/template"
	"net/http"

	"github.com/gorilla/mux"
)

//go:embed static
var Docs embed.FS

//go:embed template
var template embed.FS

const (
	apiFile   = "/static/openapi.yml"
	indexFile = "template/index.tpl"
)

func RegisterOpenAPIService(appName string, rtr *mux.Router) {
	rtr.Handle(apiFile, http.FileServer(http.FS(Docs)))
	rtr.HandleFunc("/", handler(appName))
}

// handler returns an http handler that servers OpenAPI console for an OpenAPI spec at specURL.
func handler(title string) http.HandlerFunc {
	t, parseErr := httptemplate.ParseFS(template, indexFile)

	return func(w http.ResponseWriter, req *http.Request) {
		if parseErr != nil {
			http.Error(w, parseErr.Error(), http.StatusInternalServerError)
			return
		}

		if err := t.Execute(w, struct {
			Title string
			URL   string
		}{
			title,
			apiFile,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
