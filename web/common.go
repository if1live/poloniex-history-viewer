package web

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"path"

	"github.com/if1live/poloniex-history-viewer/yui"
)

func renderErrorJSON(w http.ResponseWriter, err error, errcode int) {
	type Response struct {
		Error string `json:"error"`
	}
	resp := Response{
		Error: err.Error(),
	}
	data, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errcode)
	w.Write(data)
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		renderErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func renderPoloniexStatic(w http.ResponseWriter, r *http.Request, target string) {
	cleaned := path.Clean(target)
	basePath := yui.GetExecutablePath()
	fp := path.Join(basePath, "web", "html", cleaned)
	cleanedFp := path.Clean(fp)
	http.ServeFile(w, r, cleanedFp)
}

func renderStatic(w http.ResponseWriter, r *http.Request, target string) {
	cleaned := path.Clean(target)
	basePath := yui.GetExecutablePath()
	fp := path.Join(basePath, "web", "static", cleaned)
	cleanedFp := path.Clean(fp)
	http.ServeFile(w, r, cleanedFp)
}

func renderTemplate(w io.Writer, tplfile string, v interface{}) error {
	basePath := yui.GetExecutablePath()
	fp := path.Join(basePath, "web", "templates", tplfile)
	tpl, err := template.New(tplfile).ParseFiles(fp)
	if err != nil {
		return err
	}
	return tpl.Execute(w, v)

}
