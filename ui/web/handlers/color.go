package handlers

import (
	"net/http"

	"github.com/StanMarek/forge/core/tools"
	tooltempl "github.com/StanMarek/forge/ui/web/templates/tools"
)

func HandleColorPage(w http.ResponseWriter, r *http.Request) {
	tooltempl.ColorPage(tools.Result{}, "").Render(r.Context(), w)
}

func HandleColorProcess(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	input := r.FormValue("input")

	result := tools.ColorConvert(input)
	tooltempl.ColorOutput(result).Render(r.Context(), w)
}
