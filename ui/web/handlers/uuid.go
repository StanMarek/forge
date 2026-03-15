package handlers

import (
	"net/http"
	"strconv"

	"github.com/StanMarek/forge/core/tools"
	tooltempl "github.com/StanMarek/forge/ui/web/templates/tools"
)

func HandleUUIDPage(w http.ResponseWriter, r *http.Request) {
	// Generate one UUID on page load
	result := tools.UUIDGenerate(4, false, false)
	tooltempl.UUIDPage(result.Output, "", "generate", 4, false, false, "").Render(r.Context(), w)
}

func HandleUUIDProcess(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	mode := r.FormValue("mode")
	input := r.FormValue("input")
	uppercase := r.FormValue("uppercase") == "on"
	noHyphens := r.FormValue("no-hyphens") == "on"
	version := 4
	if v, err := strconv.Atoi(r.FormValue("version")); err == nil && (v == 4 || v == 7) {
		version = v
	}

	var output, errMsg string

	switch mode {
	case "validate":
		result := tools.UUIDValidate(input)
		output, errMsg = result.Output, result.Error
	case "parse":
		result := tools.UUIDParse(input)
		output, errMsg = result.Output, result.Error
	default: // generate
		result := tools.UUIDGenerate(version, uppercase, noHyphens)
		output, errMsg = result.Output, result.Error
	}
	tooltempl.UUIDOutput(output, errMsg).Render(r.Context(), w)
}
