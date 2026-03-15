package handlers

import (
	"net/http"

	"github.com/StanMarek/forge/core/tools"
	tooltempl "github.com/StanMarek/forge/ui/web/templates/tools"
)

func HandleYAMLPage(w http.ResponseWriter, r *http.Request) {
	tooltempl.YAMLPage(tools.Result{}, "yaml-to-json", false, "").Render(r.Context(), w)
}

func HandleYAMLProcess(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	input := r.FormValue("input")
	mode := r.FormValue("mode")
	compact := r.FormValue("compact") == "on"

	var result tools.Result
	switch mode {
	case "json-to-yaml":
		result = tools.JSONToYAML(input)
	default: // yaml-to-json
		result = tools.YAMLToJSON(input, compact)
	}
	tooltempl.YAMLOutput(result).Render(r.Context(), w)
}
