package handlers

import (
	"net/http"

	"github.com/StanMarek/forge/core/tools"
	tooltempl "github.com/StanMarek/forge/ui/web/templates/tools"
)

func HandleJWTPage(w http.ResponseWriter, r *http.Request) {
	tooltempl.JWTPage(tools.JWTDecodeResult{}, "").Render(r.Context(), w)
}

func HandleJWTProcess(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	input := r.FormValue("input")

	result := tools.JWTDecode(input)
	tooltempl.JWTOutput(result).Render(r.Context(), w)
}
