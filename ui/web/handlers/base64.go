package handlers

import (
	"net/http"

	"github.com/StanMarek/forge/core/tools"
	tooltempl "github.com/StanMarek/forge/ui/web/templates/tools"
)

func HandleBase64Page(w http.ResponseWriter, r *http.Request) {
	tooltempl.Base64Page(tools.Result{}, "encode", false, false, "").Render(r.Context(), w)
}

func HandleBase64Process(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	input := r.FormValue("input")
	mode := r.FormValue("mode")
	urlSafe := r.FormValue("url-safe") == "on"
	noPadding := r.FormValue("no-padding") == "on"

	var result tools.Result
	if mode == "decode" {
		result = tools.Base64Decode(input, urlSafe)
	} else {
		result = tools.Base64Encode(input, urlSafe, noPadding)
	}
	tooltempl.Base64Output(result).Render(r.Context(), w)
}
