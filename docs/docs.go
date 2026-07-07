package docs

import (
	"embed"
	"github.com/gorilla/mux"
	httptemplate "html/template"
	"net/http"
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
	t, _ := httptemplate.ParseFS(template, indexFile)

	return func(w http.ResponseWriter, req *http.Request) {
		t.Execute(w, struct {
			Title string
			URL   string
		}{
			title,
			apiFile,
		})
	}
}
