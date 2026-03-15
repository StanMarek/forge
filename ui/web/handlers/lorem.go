package handlers

import (
	"net/http"
	"strconv"

	"github.com/StanMarek/forge/core/tools"
	tooltempl "github.com/StanMarek/forge/ui/web/templates/tools"
)

func HandleLoremPage(w http.ResponseWriter, r *http.Request) {
	result := tools.LoremGenerate(0, 0, 3)
	tooltempl.LoremPage(result, "paragraphs", 3).Render(r.Context(), w)
}

func HandleLoremProcess(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	mode := r.FormValue("mode")
	count := 3
	if v, err := strconv.Atoi(r.FormValue("count")); err == nil && v > 0 {
		count = v
	}

	var result tools.Result
	switch mode {
	case "words":
		result = tools.LoremGenerate(count, 0, 0)
	case "sentences":
		result = tools.LoremGenerate(0, count, 0)
	default: // paragraphs
		result = tools.LoremGenerate(0, 0, count)
	}
	tooltempl.LoremOutput(result).Render(r.Context(), w)
}
