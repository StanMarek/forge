package handlers

import (
	"net/http"

	"github.com/StanMarek/forge/core/tools"
	tooltempl "github.com/StanMarek/forge/ui/web/templates/tools"
)

func HandleNumberBasePage(w http.ResponseWriter, r *http.Request) {
	tooltempl.NumberBasePage(tools.Result{}, "").Render(r.Context(), w)
}

func HandleNumberBaseProcess(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	input := r.FormValue("input")

	result := tools.NumberBaseConvert(input)
	tooltempl.NumberBaseOutput(result).Render(r.Context(), w)
}
