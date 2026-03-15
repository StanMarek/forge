package stdin

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFromPipe(t *testing.T) {
	r, w, err := os.Pipe()
	assert.NoError(t, err)

	_, err = w.WriteString("hello world\n")
	assert.NoError(t, err)
	w.Close()

	result, err := readFrom(r)
	assert.NoError(t, err)
	assert.Equal(t, "hello world", result)
}

func TestReadFromPipeMultiline(t *testing.T) {
	r, w, err := os.Pipe()
	assert.NoError(t, err)

	_, err = w.WriteString("line1\nline2\nline3\n")
	assert.NoError(t, err)
	w.Close()

	result, err := readFrom(r)
	assert.NoError(t, err)
	assert.Equal(t, "line1\nline2\nline3", result)
}

func TestReadFromEmpty(t *testing.T) {
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	w.Close()

	result, err := readFrom(r)
	assert.NoError(t, err)
	assert.Equal(t, "", result)
}

func TestReadFromPreservesInternalWhitespace(t *testing.T) {
	r, w, err := os.Pipe()
	assert.NoError(t, err)

	_, err = w.WriteString("  hello  world  \n")
	assert.NoError(t, err)
	w.Close()

	result, err := readFrom(r)
	assert.NoError(t, err)
	assert.Equal(t, "  hello  world  ", result)
}
