package handlers

import (
	"net/http"

	"github.com/StanMarek/forge/core/tools"
	tooltempl "github.com/StanMarek/forge/ui/web/templates/tools"
)

func HandleCronPage(w http.ResponseWriter, r *http.Request) {
	tooltempl.CronPage(tools.Result{}, "").Render(r.Context(), w)
}

func HandleCronProcess(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	input := r.FormValue("input")

	result := tools.CronParse(input)
	tooltempl.CronOutput(result).Render(r.Context(), w)
}
