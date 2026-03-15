package handlers

import (
	"net/http"

	"github.com/StanMarek/forge/core/tools"
	tooltempl "github.com/StanMarek/forge/ui/web/templates/tools"
)

func HandleTimestampPage(w http.ResponseWriter, r *http.Request) {
	result := tools.TimestampNow("")
	tooltempl.TimestampPage(result, "now", false, "").Render(r.Context(), w)
}

func HandleTimestampProcess(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	input := r.FormValue("input")
	mode := r.FormValue("mode")
	millis := r.FormValue("millis") == "on"
	tz := r.FormValue("tz")

	var result tools.Result
	switch mode {
	case "from-unix":
		result = tools.TimestampFromUnix(input, tz)
	case "to-unix":
		result = tools.TimestampToUnix(input, millis)
	default: // now
		result = tools.TimestampNow(tz)
	}
	tooltempl.TimestampOutput(result).Render(r.Context(), w)
}
