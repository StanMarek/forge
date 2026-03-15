package handlers

import (
	"net/http"

	"github.com/StanMarek/forge/core/tools"
	tooltempl "github.com/StanMarek/forge/ui/web/templates/tools"
)

func HandleHTMLEntityPage(w http.ResponseWriter, r *http.Request) {
	tooltempl.HTMLEntityPage(tools.Result{}, "encode", "").Render(r.Context(), w)
}

func HandleHTMLEntityProcess(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	input := r.FormValue("input")
	mode := r.FormValue("mode")

	var result tools.Result
	switch mode {
	case "decode":
		result = tools.HTMLEntityDecode(input)
	default: // encode
		result = tools.HTMLEntityEncode(input)
	}
	tooltempl.HTMLEntityOutput(result).Render(r.Context(), w)
}
