package handlers

import (
	"net/http"

	"github.com/StanMarek/forge/core/tools"
	tooltempl "github.com/StanMarek/forge/ui/web/templates/tools"
)

func HandleXMLPage(w http.ResponseWriter, r *http.Request) {
	tooltempl.XMLPage(tools.Result{}, "format", "").Render(r.Context(), w)
}

func HandleXMLProcess(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	input := r.FormValue("input")
	mode := r.FormValue("mode")

	var result tools.Result
	switch mode {
	case "minify":
		result = tools.XMLMinify(input)
	default: // format
		result = tools.XMLFormat(input)
	}
	tooltempl.XMLOutput(result).Render(r.Context(), w)
}
