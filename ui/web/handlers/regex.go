package handlers

import (
	"net/http"

	"github.com/StanMarek/forge/core/tools"
	tooltempl "github.com/StanMarek/forge/ui/web/templates/tools"
)

func HandleRegexPage(w http.ResponseWriter, r *http.Request) {
	tooltempl.RegexPage(tools.Result{}, true, "", "").Render(r.Context(), w)
}

func HandleRegexProcess(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	pattern := r.FormValue("pattern")
	input := r.FormValue("input")
	global := r.FormValue("global") == "on"

	result := tools.RegexTest(pattern, input, global)
	tooltempl.RegexOutput(result).Render(r.Context(), w)
}
