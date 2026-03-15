package handlers

import (
	"net/http"

	"github.com/StanMarek/forge/ui/web/templates"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	templates.IndexPage().Render(r.Context(), w)
}
