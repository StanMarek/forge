package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------- Tool metadata ----------

func TestColorTool_Metadata(t *testing.T) {
	tool := ColorTool{}

	assert.Equal(t, "Color Converter", tool.Name())
	assert.Equal(t, "color", tool.ID())
	assert.Equal(t, "Converters", tool.Category())
	assert.NotEmpty(t, tool.Description())

	keywords := tool.Keywords()
	for _, kw := range []string{"color", "hex", "rgb", "hsl", "convert"} {
		assert.Contains(t, keywords, kw)
	}
}

func TestColorTool_DetectFromClipboard(t *testing.T) {
	tool := ColorTool{}
	assert.True(t, tool.DetectFromClipboard("#ff9800"))
	assert.True(t, tool.DetectFromClipboard("#f90"))
	assert.True(t, tool.DetectFromClipboard("rgb(255, 152, 0)"))
	assert.True(t, tool.DetectFromClipboard("hsl(36, 100%, 50%)"))
	assert.False(t, tool.DetectFromClipboard("hello"))
	assert.False(t, tool.DetectFromClipboard(""))
	assert.False(t, tool.DetectFromClipboard("#gg0000"))
}

// ---------- Hex 6-digit ----------

func TestColorConvert_Hex6(t *testing.T) {
	r := ColorConvert("#ff9800")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "Hex:  #ff9800")
	assert.Contains(t, r.Output, "RGB:  rgb(255, 152, 0)")
	assert.Contains(t, r.Output, "HSL:  hsl(36, 100%, 50%)")
}

// ---------- Hex 3-digit ----------

func TestColorConvert_Hex3(t *testing.T) {
	r := ColorConvert("#f90")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "Hex:  #ff9900")
	assert.Contains(t, r.Output, "RGB:  rgb(255, 153, 0)")
}

// ---------- RGB ----------

func TestColorConvert_RGB(t *testing.T) {
	r := ColorConvert("rgb(255,152,0)")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "Hex:  #ff9800")
	assert.Contains(t, r.Output, "RGB:  rgb(255, 152, 0)")
}

// ---------- HSL ----------

func TestColorConvert_HSL(t *testing.T) {
	r := ColorConvert("hsl(36, 100%, 50%)")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "RGB:  rgb(255, 153, 0)")
	assert.Contains(t, r.Output, "Hex:  #ff9900")
}

// ---------- Black ----------

func TestColorConvert_Black(t *testing.T) {
	r := ColorConvert("#000000")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "Hex:  #000000")
	assert.Contains(t, r.Output, "RGB:  rgb(0, 0, 0)")
	assert.Contains(t, r.Output, "HSL:  hsl(0, 0%, 0%)")
}

// ---------- White ----------

func TestColorConvert_White(t *testing.T) {
	r := ColorConvert("#ffffff")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "Hex:  #ffffff")
	assert.Contains(t, r.Output, "RGB:  rgb(255, 255, 255)")
	assert.Contains(t, r.Output, "HSL:  hsl(0, 0%, 100%)")
}

// ---------- Invalid input ----------

func TestColorConvert_Invalid(t *testing.T) {
	r := ColorConvert("not a color")
	assert.NotEmpty(t, r.Error)
	assert.Contains(t, r.Error, "unrecognized color format")
	assert.Empty(t, r.Output)
}

func TestColorConvert_EmptyInput(t *testing.T) {
	r := ColorConvert("")
	assert.NotEmpty(t, r.Error)
	assert.Empty(t, r.Output)
}

// ---------- RGB with spaces ----------

func TestColorConvert_RGBWithSpaces(t *testing.T) {
	r := ColorConvert("rgb( 255 , 152 , 0 )")
	require.Empty(t, r.Error)
	assert.Contains(t, r.Output, "Hex:  #ff9800")
}

// ---------- Tool interface compliance ----------

func TestColorTool_ImplementsToolInterface(t *testing.T) {
	var _ Tool = ColorTool{}
}
