package registry

import "github.com/StanMarek/forge/core/tools"

// Default creates a registry pre-loaded with all Tier-1 tools.
func Default() *Registry {
	r := New()
	r.Register(tools.Base64Tool{})
	r.Register(tools.JWTTool{})
	r.Register(tools.JSONTool{})
	r.Register(tools.HashTool{})
	r.Register(tools.URLTool{})
	r.Register(tools.UUIDTool{})
	return r
}
