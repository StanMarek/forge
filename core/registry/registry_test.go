package registry

import (
	"testing"

	"github.com/StanMarek/forge/core/tools"
	"github.com/stretchr/testify/assert"
)

func TestRegistryRegisterAndAll(t *testing.T) {
	r := New()
	r.Register(tools.Base64Tool{})
	r.Register(tools.JWTTool{})

	all := r.All()
	assert.Len(t, all, 2)
}

func TestRegistryDuplicateRegister(t *testing.T) {
	r := New()
	r.Register(tools.Base64Tool{})
	r.Register(tools.Base64Tool{})

	all := r.All()
	assert.Len(t, all, 1)
}

func TestRegistryByID(t *testing.T) {
	r := New()
	r.Register(tools.Base64Tool{})

	tool, ok := r.ByID("base64")
	assert.True(t, ok)
	assert.Equal(t, "base64", tool.ID())

	_, ok = r.ByID("nonexistent")
	assert.False(t, ok)
}

func TestRegistryByCategory(t *testing.T) {
	r := New()
	r.Register(tools.Base64Tool{})
	r.Register(tools.JWTTool{})
	r.Register(tools.HashTool{})

	encoders := r.ByCategory("Encoders")
	assert.Len(t, encoders, 2)

	generators := r.ByCategory("Generators")
	assert.Len(t, generators, 1)

	empty := r.ByCategory("Nonexistent")
	assert.Len(t, empty, 0)
}

func TestRegistrySearch(t *testing.T) {
	r := New()
	r.Register(tools.Base64Tool{})
	r.Register(tools.JWTTool{})
	r.Register(tools.JSONTool{})

	// "minify" only matches JSON tool's keywords
	results := r.Search("minify")
	assert.Len(t, results, 1)
	assert.Equal(t, "json", results[0].ID())

	// "encode" matches base64 and possibly others
	results = r.Search("encode")
	assert.GreaterOrEqual(t, len(results), 1)

	results = r.Search("zzzzz")
	assert.Len(t, results, 0)
}

func TestRegistrySearchCaseInsensitive(t *testing.T) {
	r := New()
	r.Register(tools.Base64Tool{})

	results := r.Search("BASE64")
	assert.Len(t, results, 1)
}

func TestRegistryDetect(t *testing.T) {
	r := Default()

	// JWT should be detected first (highest priority)
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"
	matches := r.Detect(jwt)
	assert.NotEmpty(t, matches)
	assert.Equal(t, "jwt", matches[0].ID())

	// UUID should be detected
	matches = r.Detect("550e8400-e29b-41d4-a716-446655440000")
	assert.NotEmpty(t, matches)
	assert.Equal(t, "uuid", matches[0].ID())

	// URL should be detected
	matches = r.Detect("https://example.com")
	assert.NotEmpty(t, matches)
	assert.Equal(t, "url", matches[0].ID())

	// JSON should be detected
	matches = r.Detect(`{"key":"value"}`)
	assert.NotEmpty(t, matches)
	assert.Equal(t, "json", matches[0].ID())

	// No match
	matches = r.Detect("just plain text")
	assert.Empty(t, matches)
}

func TestRegistryDetectPriority(t *testing.T) {
	r := Default()

	// Valid JSON that could also match base64 — JSON should rank higher
	matches := r.Detect(`{"a":1}`)
	if len(matches) > 1 {
		// If both json and base64 match, json should come first
		jsonIdx := -1
		base64Idx := -1
		for i, m := range matches {
			if m.ID() == "json" {
				jsonIdx = i
			}
			if m.ID() == "base64" {
				base64Idx = i
			}
		}
		if jsonIdx >= 0 && base64Idx >= 0 {
			assert.Less(t, jsonIdx, base64Idx)
		}
	}
}

func TestDefaultRegistry(t *testing.T) {
	r := Default()
	all := r.All()
	assert.Len(t, all, 21)

	ids := make(map[string]bool)
	for _, tool := range all {
		ids[tool.ID()] = true
	}
	// Tier 1
	assert.True(t, ids["base64"])
	assert.True(t, ids["jwt"])
	assert.True(t, ids["json"])
	assert.True(t, ids["hash"])
	assert.True(t, ids["url"])
	assert.True(t, ids["uuid"])
	// Tier 2
	assert.True(t, ids["yaml"])
	assert.True(t, ids["timestamp"])
	assert.True(t, ids["number-base"])
	assert.True(t, ids["regex"])
	assert.True(t, ids["html-entity"])
	assert.True(t, ids["password"])
	assert.True(t, ids["lorem"])
	// High-value additions
	assert.True(t, ids["color"])
	assert.True(t, ids["cron"])
	assert.True(t, ids["text-escape"])
	assert.True(t, ids["gzip"])
	assert.True(t, ids["text-stats"])
	assert.True(t, ids["diff"])
	assert.True(t, ids["xml"])
	assert.True(t, ids["csv"])
}
