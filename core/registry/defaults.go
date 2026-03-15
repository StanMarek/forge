package registry

import "github.com/StanMarek/forge/core/tools"

// Default creates a registry pre-loaded with all tools.
func Default() *Registry {
	r := New()
	// Tier 1 — Encoders
	r.Register(tools.Base64Tool{})
	r.Register(tools.JWTTool{})
	r.Register(tools.URLTool{})
	r.Register(tools.HTMLEntityTool{})
	r.Register(tools.TextEscapeTool{})
	r.Register(tools.GZipTool{})
	// Tier 1 — Formatters
	r.Register(tools.JSONTool{})
	r.Register(tools.XMLTool{})
	// Tier 1 — Generators
	r.Register(tools.HashTool{})
	r.Register(tools.UUIDTool{})
	r.Register(tools.PasswordTool{})
	r.Register(tools.LoremTool{})
	// Tier 2 — Converters
	r.Register(tools.YAMLTool{})
	r.Register(tools.TimestampTool{})
	r.Register(tools.NumberBaseTool{})
	r.Register(tools.ColorTool{})
	r.Register(tools.CronTool{})
	r.Register(tools.CSVTool{})
	// Testers
	r.Register(tools.RegexTool{})
	// Text
	r.Register(tools.TextStatsTool{})
	r.Register(tools.DiffTool{})
	return r
}
