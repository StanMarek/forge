package registry

import (
	"sort"
	"strings"

	"github.com/StanMarek/forge/core/tools"
)

// detectionPriority defines the order for smart clipboard detection.
// Lower index = higher priority.
var detectionPriority = map[string]int{
	"jwt":         0,
	"uuid":        1,
	"url":         2,
	"json":        3,
	"base64":      4,
	"timestamp":   5,
	"html-entity": 6,
	"number-base": 7,
}

// Registry holds registered tools and provides lookup, search, and detection.
type Registry struct {
	tools map[string]tools.Tool
	order []string
}

// New creates an empty Registry.
func New() *Registry {
	return &Registry{
		tools: make(map[string]tools.Tool),
	}
}

// Register adds a tool to the registry.
func (r *Registry) Register(tool tools.Tool) {
	id := tool.ID()
	if _, exists := r.tools[id]; !exists {
		r.order = append(r.order, id)
	}
	r.tools[id] = tool
}

// All returns all registered tools in registration order.
func (r *Registry) All() []tools.Tool {
	result := make([]tools.Tool, 0, len(r.order))
	for _, id := range r.order {
		result = append(result, r.tools[id])
	}
	return result
}

// ByID returns a tool by its ID.
func (r *Registry) ByID(id string) (tools.Tool, bool) {
	tool, ok := r.tools[id]
	return tool, ok
}

// ByCategory returns all tools in a given category.
func (r *Registry) ByCategory(category string) []tools.Tool {
	var result []tools.Tool
	for _, id := range r.order {
		if r.tools[id].Category() == category {
			result = append(result, r.tools[id])
		}
	}
	return result
}

// Search returns tools matching a query against name, ID, and keywords.
func (r *Registry) Search(query string) []tools.Tool {
	q := strings.ToLower(query)
	var result []tools.Tool
	for _, id := range r.order {
		tool := r.tools[id]
		if matchesTool(tool, q) {
			result = append(result, tool)
		}
	}
	return result
}

// Detect returns tools whose DetectFromClipboard returns true,
// sorted by detection priority (JWT > UUID > URL > JSON > Base64 > ...).
func (r *Registry) Detect(clipboard string) []tools.Tool {
	var matches []tools.Tool
	for _, id := range r.order {
		tool := r.tools[id]
		if tool.DetectFromClipboard(clipboard) {
			matches = append(matches, tool)
		}
	}
	sort.SliceStable(matches, func(i, j int) bool {
		return priorityOf(matches[i].ID()) < priorityOf(matches[j].ID())
	})
	return matches
}

func matchesTool(tool tools.Tool, query string) bool {
	if strings.Contains(strings.ToLower(tool.Name()), query) {
		return true
	}
	if strings.Contains(strings.ToLower(tool.ID()), query) {
		return true
	}
	for _, kw := range tool.Keywords() {
		if strings.Contains(strings.ToLower(kw), query) {
			return true
		}
	}
	return false
}

func priorityOf(id string) int {
	if p, ok := detectionPriority[id]; ok {
		return p
	}
	return 999
}
