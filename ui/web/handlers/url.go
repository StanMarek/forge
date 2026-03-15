package handlers

import (
	"net/http"

	"github.com/StanMarek/forge/core/tools"
	tooltempl "github.com/StanMarek/forge/ui/web/templates/tools"
)

func HandleURLPage(w http.ResponseWriter, r *http.Request) {
	tooltempl.URLPage("", "", "parse", false, "").Render(r.Context(), w)
}

func HandleURLProcess(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	input := r.FormValue("input")
	mode := r.FormValue("mode")
	component := r.FormValue("component") == "on"

	var output, errMsg string

	switch mode {
	case "encode":
		result := tools.URLEncode(input, component)
		output, errMsg = result.Output, result.Error
	case "decode":
		result := tools.URLDecode(input)
		output, errMsg = result.Output, result.Error
	default: // parse
		result := tools.URLParse(input)
		output, errMsg = result.Output, result.Error
	}
	tooltempl.URLOutput(output, errMsg).Render(r.Context(), w)
}
