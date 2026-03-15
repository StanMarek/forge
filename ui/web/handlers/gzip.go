package handlers

import (
	"net/http"

	"github.com/StanMarek/forge/core/tools"
	tooltempl "github.com/StanMarek/forge/ui/web/templates/tools"
)

func HandleGZipPage(w http.ResponseWriter, r *http.Request) {
	tooltempl.GZipPage(tools.Result{}, "compress", "").Render(r.Context(), w)
}

func HandleGZipProcess(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	input := r.FormValue("input")
	mode := r.FormValue("mode")

	var result tools.Result
	switch mode {
	case "decompress":
		result = tools.GZipDecompress(input)
	default: // compress
		result = tools.GZipCompress(input)
	}
	tooltempl.GZipOutput(result).Render(r.Context(), w)
}
