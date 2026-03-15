package handlers

import (
	"net/http"

	"github.com/StanMarek/forge/core/tools"
	tooltempl "github.com/StanMarek/forge/ui/web/templates/tools"
)

func HandleCSVPage(w http.ResponseWriter, r *http.Request) {
	tooltempl.CSVPage(tools.Result{}, "json-to-csv", ",", "").Render(r.Context(), w)
}

func HandleCSVProcess(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	input := r.FormValue("input")
	mode := r.FormValue("mode")
	delimiter := r.FormValue("delimiter")
	if delimiter == "" {
		delimiter = ","
	}

	var result tools.Result
	switch mode {
	case "csv-to-json":
		result = tools.CSVToJSON(input, delimiter)
	default: // json-to-csv
		result = tools.JSONToCSV(input, delimiter)
	}
	tooltempl.CSVOutput(result).Render(r.Context(), w)
}
