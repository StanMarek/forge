package tools

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------- Tool metadata ----------

func TestLoremTool_Metadata(t *testing.T) {
	tool := LoremTool{}

	assert.Equal(t, "Lorem Ipsum Generator", tool.Name())
	assert.Equal(t, "lorem", tool.ID())
	assert.Equal(t, "Generators", tool.Category())
	assert.NotEmpty(t, tool.Description())

	keywords := tool.Keywords()
	for _, kw := range []string{"lorem", "ipsum", "placeholder", "text"} {
		assert.Contains(t, keywords, kw)
	}
}

func TestLoremTool_DetectFromClipboard(t *testing.T) {
	tool := LoremTool{}
	assert.False(t, tool.DetectFromClipboard("anything"))
	assert.False(t, tool.DetectFromClipboard(""))
}

// ---------- Words mode ----------

func TestLoremGenerate_10Words(t *testing.T) {
	r := LoremGenerate(10, 0, 0)
	require.Empty(t, r.Error)

	words := strings.Fields(r.Output)
	assert.Len(t, words, 10)
}

func TestLoremGenerate_1Word(t *testing.T) {
	r := LoremGenerate(1, 0, 0)
	require.Empty(t, r.Error)

	words := strings.Fields(r.Output)
	assert.Len(t, words, 1)
}

// ---------- Sentences mode ----------

func TestLoremGenerate_3Sentences(t *testing.T) {
	r := LoremGenerate(0, 3, 0)
	require.Empty(t, r.Error)

	// Count sentences by splitting on ". " and trailing "."
	// Each sentence ends with a period.
	sentences := countSentences(r.Output)
	assert.Equal(t, 3, sentences)
}

func TestLoremGenerate_SentenceFormat(t *testing.T) {
	r := LoremGenerate(0, 1, 0)
	require.Empty(t, r.Error)

	// Should end with a period.
	assert.True(t, strings.HasSuffix(r.Output, "."))

	// First character should be uppercase.
	assert.True(t, r.Output[0] >= 'A' && r.Output[0] <= 'Z', "sentence should start with uppercase letter")

	// Word count should be 8-15.
	// Remove the trailing period for word counting.
	trimmed := strings.TrimSuffix(r.Output, ".")
	words := strings.Fields(trimmed)
	assert.GreaterOrEqual(t, len(words), 8)
	assert.LessOrEqual(t, len(words), 15)
}

// ---------- Paragraphs mode ----------

func TestLoremGenerate_2Paragraphs(t *testing.T) {
	r := LoremGenerate(0, 0, 2)
	require.Empty(t, r.Error)

	paragraphs := strings.Split(r.Output, "\n\n")
	assert.Len(t, paragraphs, 2)

	// Each paragraph should have 4-7 sentences.
	for _, p := range paragraphs {
		sc := countSentences(p)
		assert.GreaterOrEqual(t, sc, 4, "paragraph should have at least 4 sentences")
		assert.LessOrEqual(t, sc, 7, "paragraph should have at most 7 sentences")
	}
}

// ---------- Error cases ----------

func TestLoremGenerate_ErrorOnZero(t *testing.T) {
	r := LoremGenerate(0, 0, 0)
	assert.NotEmpty(t, r.Error)
	assert.Contains(t, r.Error, "one of words, sentences, or paragraphs must be greater than 0")
	assert.Empty(t, r.Output)
}

func TestLoremGenerate_ErrorOnMultipleModes(t *testing.T) {
	r := LoremGenerate(5, 3, 0)
	assert.NotEmpty(t, r.Error)
	assert.Contains(t, r.Error, "only one of words, sentences, or paragraphs may be greater than 0")
	assert.Empty(t, r.Output)
}

func TestLoremGenerate_ErrorOnAllModes(t *testing.T) {
	r := LoremGenerate(5, 3, 2)
	assert.NotEmpty(t, r.Error)
	assert.Empty(t, r.Output)
}

// ---------- Tool interface compliance ----------

func TestLoremTool_ImplementsToolInterface(t *testing.T) {
	var _ Tool = LoremTool{}
}

// ---------- helpers ----------

// countSentences counts sentences by counting periods that end sentences.
func countSentences(text string) int {
	count := 0
	for _, c := range text {
		if c == '.' {
			count++
		}
	}
	return count
}
