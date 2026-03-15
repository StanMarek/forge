package handlers

import (
	"net/http"

	"github.com/StanMarek/forge/core/tools"
	tooltempl "github.com/StanMarek/forge/ui/web/templates/tools"
)

func HandleTextEscapePage(w http.ResponseWriter, r *http.Request) {
	tooltempl.TextEscapePage(tools.Result{}, "escape", "").Render(r.Context(), w)
}

func HandleTextEscapeProcess(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	input := r.FormValue("input")
	mode := r.FormValue("mode")

	var result tools.Result
	switch mode {
	case "unescape":
		result = tools.TextUnescape(input)
	default: // escape
		result = tools.TextEscape(input)
	}
	tooltempl.TextEscapeOutput(result).Render(r.Context(), w)
}
