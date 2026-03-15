package handlers

import (
	"net/http"

	"github.com/StanMarek/forge/core/tools"
	tooltempl "github.com/StanMarek/forge/ui/web/templates/tools"
)

func HandleHashPage(w http.ResponseWriter, r *http.Request) {
	tooltempl.HashPage(tooltempl.HashResult{}, false, "").Render(r.Context(), w)
}

func HandleHashProcess(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	input := r.FormValue("input")
	uppercase := r.FormValue("uppercase") == "on"

	result := tooltempl.HashResult{
		MD5:    tools.Hash(input, "md5", uppercase).Output,
		SHA1:   tools.Hash(input, "sha1", uppercase).Output,
		SHA256: tools.Hash(input, "sha256", uppercase).Output,
		SHA512: tools.Hash(input, "sha512", uppercase).Output,
	}
	tooltempl.HashOutput(result).Render(r.Context(), w)
}
