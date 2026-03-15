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

	return r
}

// Start launches the web server.
func Start(host string, port int) error {
	r := NewRouter()
	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Forge web server running at http://%s\n", addr)
	return http.ListenAndServe(addr, r)
}
