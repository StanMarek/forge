package handlers

import (
	"net/http"

	"github.com/StanMarek/forge/core/tools"
	tooltempl "github.com/StanMarek/forge/ui/web/templates/tools"
)

func HandleDiffPage(w http.ResponseWriter, r *http.Request) {
	tooltempl.DiffPage(tools.Result{}, "", "").Render(r.Context(), w)
}

func HandleDiffProcess(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	textA := r.FormValue("text-a")
	textB := r.FormValue("text-b")

	result := tools.DiffText(textA, textB)
	tooltempl.DiffOutput(result).Render(r.Context(), w)
}
