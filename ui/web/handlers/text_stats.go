package handlers

import (
	"net/http"

	"github.com/StanMarek/forge/core/tools"
	tooltempl "github.com/StanMarek/forge/ui/web/templates/tools"
)

func HandleTextStatsPage(w http.ResponseWriter, r *http.Request) {
	tooltempl.TextStatsPage(tools.Result{}, "stats", "").Render(r.Context(), w)
}

func HandleTextStatsProcess(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	input := r.FormValue("input")
	mode := r.FormValue("mode")

	var result tools.Result
	switch mode {
	case "lower", "upper", "title", "camel", "snake", "kebab":
		result = tools.TextCaseConvert(input, mode)
	default: // stats
		result = tools.TextStats(input)
	}
	tooltempl.TextStatsOutput(result).Render(r.Context(), w)
}
