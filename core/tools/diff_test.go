package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------- Tool metadata ----------

func TestDiffTool_Metadata(t *testing.T) {
	tool := DiffTool{}

	assert.Equal(t, "Text Diff", tool.Name())
	assert.Equal(t, "diff", tool.ID())
	assert.Equal(t, "Text", tool.Category())
	assert.NotEmpty(t, tool.Description())

	keywords := tool.Keywords()
	for _, kw := range []string{"diff", "compare", "text", "difference"} {
		assert.Contains(t, keywords, kw)
	}
}

func TestDiffTool_DetectFromClipboard(t *testing.T) {
	tool := DiffTool{}
	assert.False(t, tool.DetectFromClipboard("anything"))
	assert.False(t, tool.DetectFromClipboard(""))
}

// ---------- Identical texts ----------

func TestDiffText_Identical(t *testing.T) {
	r := DiffText("hello\nworld", "hello\nworld")
	require.Empty(t, r.Error)
	assert.Equal(t, "Texts are identical", r.Output)
}

func TestDiffText_IdenticalSingleLine(t *testing.T) {
	r := DiffText("hello", "hello")
	require.Empty(t, r.Error)
	assert.Equal(t, "Texts are identical", r.Output)
}

func TestDiffText_BothEmpty(t *testing.T) {
	r := DiffText("", "")
	require.Empty(t, r.Error)
	assert.Equal(t, "Texts are identical", r.Output)
}

// ---------- Added lines ----------

func TestDiffText_AddedLines(t *testing.T) {
	r := DiffText("line1\nline3", "line1\nline2\nline3")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "+line2")
	assert.Contains(t, r.Output, " line1")
	assert.Contains(t, r.Output, " line3")
}

func TestDiffText_AddedAtEnd(t *testing.T) {
	r := DiffText("a\nb", "a\nb\nc")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "+c")
	assert.Contains(t, r.Output, " a")
	assert.Contains(t, r.Output, " b")
}

// ---------- Removed lines ----------

func TestDiffText_RemovedLines(t *testing.T) {
	r := DiffText("line1\nline2\nline3", "line1\nline3")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "-line2")
	assert.Contains(t, r.Output, " line1")
	assert.Contains(t, r.Output, " line3")
}

func TestDiffText_RemovedFromEnd(t *testing.T) {
	r := DiffText("a\nb\nc", "a\nb")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "-c")
}

// ---------- Changed lines ----------

func TestDiffText_ChangedLine(t *testing.T) {
	r := DiffText("hello\nold\nworld", "hello\nnew\nworld")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "-old")
	assert.Contains(t, r.Output, "+new")
	assert.Contains(t, r.Output, " hello")
	assert.Contains(t, r.Output, " world")
}

// ---------- Empty inputs ----------

func TestDiffText_EmptyA(t *testing.T) {
	r := DiffText("", "line1\nline2")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "+line1")
	assert.Contains(t, r.Output, "+line2")
}

func TestDiffText_EmptyB(t *testing.T) {
	r := DiffText("line1\nline2", "")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "-line1")
	assert.Contains(t, r.Output, "-line2")
}

// ---------- Completely different ----------

func TestDiffText_CompletelyDifferent(t *testing.T) {
	r := DiffText("aaa\nbbb", "xxx\nyyy")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "-aaa")
	assert.Contains(t, r.Output, "-bbb")
	assert.Contains(t, r.Output, "+xxx")
	assert.Contains(t, r.Output, "+yyy")
}

// ---------- Unified diff header ----------

func TestDiffText_HasHeader(t *testing.T) {
	r := DiffText("a", "b")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "--- Text A")
	assert.Contains(t, r.Output, "+++ Text B")
}

// ---------- Tool interface compliance ----------

func TestDiffTool_ImplementsToolInterface(t *testing.T) {
	var _ Tool = DiffTool{}
}
