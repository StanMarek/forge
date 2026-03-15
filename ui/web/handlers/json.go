package handlers

import (
	"net/http"

	"github.com/StanMarek/forge/core/tools"
	tooltempl "github.com/StanMarek/forge/ui/web/templates/tools"
)

func HandleJSONPage(w http.ResponseWriter, r *http.Request) {
	tooltempl.JSONPage(tools.Result{}, "format", false, "").Render(r.Context(), w)
}

func HandleJSONProcess(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	input := r.FormValue("input")
	mode := r.FormValue("mode")
	sortKeys := r.FormValue("sort-keys") == "on"

	var result tools.Result
	switch mode {
	case "minify":
		result = tools.JSONMinify(input)
	case "validate":
		result = tools.JSONValidate(input)
	default:
		result = tools.JSONFormat(input, 2, sortKeys, false)
	}
	tooltempl.JSONOutput(result).Render(r.Context(), w)
}
