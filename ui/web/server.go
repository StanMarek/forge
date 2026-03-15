package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/StanMarek/forge/ui/web/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed static/*
var staticFS embed.FS

// NewRouter creates the Chi router with all routes.
func NewRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Static files
	sub, _ := fs.Sub(staticFS, "static")
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServerFS(sub)))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Index
	r.Get("/", handlers.HandleIndex)

	// Tools
	r.Get("/tools/base64", handlers.HandleBase64Page)
	r.Post("/tools/base64", handlers.HandleBase64Process)

	r.Get("/tools/jwt", handlers.HandleJWTPage)
	r.Post("/tools/jwt", handlers.HandleJWTProcess)

	r.Get("/tools/json", handlers.HandleJSONPage)
	r.Post("/tools/json", handlers.HandleJSONProcess)

	r.Get("/tools/hash", handlers.HandleHashPage)
	r.Post("/tools/hash", handlers.HandleHashProcess)

	r.Get("/tools/url", handlers.HandleURLPage)
	r.Post("/tools/url", handlers.HandleURLProcess)

	r.Get("/tools/uuid", handlers.HandleUUIDPage)
	r.Post("/tools/uuid", handlers.HandleUUIDProcess)

	r.Get("/tools/yaml", handlers.HandleYAMLPage)
	r.Post("/tools/yaml", handlers.HandleYAMLProcess)

	r.Get("/tools/timestamp", handlers.HandleTimestampPage)
	r.Post("/tools/timestamp", handlers.HandleTimestampProcess)

	r.Get("/tools/number-base", handlers.HandleNumberBasePage)
	r.Post("/tools/number-base", handlers.HandleNumberBaseProcess)

	r.Get("/tools/regex", handlers.HandleRegexPage)
	r.Post("/tools/regex", handlers.HandleRegexProcess)

	r.Get("/tools/html-entity", handlers.HandleHTMLEntityPage)
	r.Post("/tools/html-entity", handlers.HandleHTMLEntityProcess)

	r.Get("/tools/password", handlers.HandlePasswordPage)
	r.Post("/tools/password", handlers.HandlePasswordProcess)

	r.Get("/tools/lorem", handlers.HandleLoremPage)
	r.Post("/tools/lorem", handlers.HandleLoremProcess)

	r.Get("/tools/color", handlers.HandleColorPage)
	r.Post("/tools/color", handlers.HandleColorProcess)

	r.Get("/tools/cron", handlers.HandleCronPage)
	r.Post("/tools/cron", handlers.HandleCronProcess)

	r.Get("/tools/text-escape", handlers.HandleTextEscapePage)
	r.Post("/tools/text-escape", handlers.HandleTextEscapeProcess)

	r.Get("/tools/gzip", handlers.HandleGZipPage)
	r.Post("/tools/gzip", handlers.HandleGZipProcess)

	r.Get("/tools/text-stats", handlers.HandleTextStatsPage)
	r.Post("/tools/text-stats", handlers.HandleTextStatsProcess)

	r.Get("/tools/diff", handlers.HandleDiffPage)
	r.Post("/tools/diff", handlers.HandleDiffProcess)

	r.Get("/tools/xml", handlers.HandleXMLPage)
	r.Post("/tools/xml", handlers.HandleXMLProcess)

	r.Get("/tools/csv", handlers.HandleCSVPage)
	r.Post("/tools/csv", handlers.HandleCSVProcess)

	return r
}

// Start launches the web server.
func Start(host string, port int) error {
	r := NewRouter()
	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Forge web server running at http://%s\n", addr)
	return http.ListenAndServe(addr, r)
}
