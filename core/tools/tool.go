package tools

// Result is the standard return type for tool operations.
// Success is indicated by Error == "".
type Result struct {
	Output string
	Error  string
}

// Tool defines the metadata interface for tool discovery and routing.
// Tool logic lives in standalone functions, NOT on this interface.
type Tool interface {
	Name() string
	ID() string
	Description() string
	Category() string
	Keywords() []string
	DetectFromClipboard(s string) bool
}
