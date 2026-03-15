package handlers

import (
	"net/http"
	"strconv"

	"github.com/StanMarek/forge/core/tools"
	tooltempl "github.com/StanMarek/forge/ui/web/templates/tools"
)

func HandlePasswordPage(w http.ResponseWriter, r *http.Request) {
	result := tools.PasswordGenerate(16, true, true, true, true, "")
	tooltempl.PasswordPage(result, 16, true, true, true, true).Render(r.Context(), w)
}

func HandlePasswordProcess(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	length := 16
	if v, err := strconv.Atoi(r.FormValue("length")); err == nil && v > 0 {
		length = v
	}
	uppercase := r.FormValue("uppercase") == "on"
	lowercase := r.FormValue("lowercase") == "on"
	digits := r.FormValue("digits") == "on"
	symbols := r.FormValue("symbols") == "on"

	result := tools.PasswordGenerate(length, uppercase, lowercase, digits, symbols, "")
	tooltempl.PasswordOutput(result).Render(r.Context(), w)
}
