package tools

import (
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------- Tool metadata ----------

func TestPasswordTool_Metadata(t *testing.T) {
	tool := PasswordTool{}

	assert.Equal(t, "Password Generator", tool.Name())
	assert.Equal(t, "password", tool.ID())
	assert.Equal(t, "Generators", tool.Category())
	assert.NotEmpty(t, tool.Description())

	keywords := tool.Keywords()
	for _, kw := range []string{"password", "generate", "random", "secure"} {
		assert.Contains(t, keywords, kw)
	}
}

func TestPasswordTool_DetectFromClipboard(t *testing.T) {
	tool := PasswordTool{}
	assert.False(t, tool.DetectFromClipboard("anything"))
	assert.False(t, tool.DetectFromClipboard(""))
	assert.False(t, tool.DetectFromClipboard("P@ssw0rd!"))
}

// ---------- Default (all types, length 16) ----------

func TestPassword_Default(t *testing.T) {
	r := PasswordGenerate(16, true, true, true, true, "")
	require.Empty(t, r.Error)
	assert.Len(t, r.Output, 16)
}

// ---------- Verify output length ----------

func TestPassword_VerifyLength(t *testing.T) {
	for _, length := range []int{1, 8, 16, 32, 64, 128, 256} {
		r := PasswordGenerate(length, true, true, true, false, "")
		require.Empty(t, r.Error, "length %d should succeed", length)
		assert.Len(t, r.Output, length, "output should be exactly %d chars", length)
	}
}

// ---------- No symbols ----------

func TestPassword_NoSymbols(t *testing.T) {
	r := PasswordGenerate(100, true, true, true, false, "")
	require.Empty(t, r.Error)
	assert.Len(t, r.Output, 100)

	for _, ch := range r.Output {
		assert.True(t, unicode.IsLetter(ch) || unicode.IsDigit(ch),
			"char %q should be alphanumeric", string(ch))
	}
}

// ---------- Custom symbol set ----------

func TestPassword_CustomSymbolSet(t *testing.T) {
	r := PasswordGenerate(100, false, false, false, true, "!?")
	require.Empty(t, r.Error)
	assert.Len(t, r.Output, 100)

	for _, ch := range r.Output {
		assert.True(t, ch == '!' || ch == '?',
			"char %q should be from custom symbol set", string(ch))
	}
}

// ---------- Verify character types present ----------

func TestPassword_VerifyCharacterTypesPresent(t *testing.T) {
	// Generate a long password to make it statistically near-certain all types appear
	r := PasswordGenerate(200, true, true, true, true, "")
	require.Empty(t, r.Error)

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSymbol := false

	for _, ch := range r.Output {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		default:
			hasSymbol = true
		}
	}

	assert.True(t, hasUpper, "password should contain uppercase letters")
	assert.True(t, hasLower, "password should contain lowercase letters")
	assert.True(t, hasDigit, "password should contain digits")
	assert.True(t, hasSymbol, "password should contain symbols")
}

// ---------- Only uppercase ----------

func TestPassword_OnlyUppercase(t *testing.T) {
	r := PasswordGenerate(50, true, false, false, false, "")
	require.Empty(t, r.Error)

	for _, ch := range r.Output {
		assert.True(t, strings.ContainsRune(uppercaseChars, ch),
			"char %q should be uppercase letter", string(ch))
	}
}

// ---------- Only lowercase ----------

func TestPassword_OnlyLowercase(t *testing.T) {
	r := PasswordGenerate(50, false, true, false, false, "")
	require.Empty(t, r.Error)

	for _, ch := range r.Output {
		assert.True(t, strings.ContainsRune(lowercaseChars, ch),
			"char %q should be lowercase letter", string(ch))
	}
}

// ---------- Only digits ----------

func TestPassword_OnlyDigits(t *testing.T) {
	r := PasswordGenerate(50, false, false, true, false, "")
	require.Empty(t, r.Error)

	for _, ch := range r.Output {
		assert.True(t, unicode.IsDigit(ch),
			"char %q should be a digit", string(ch))
	}
}

// ---------- Length bounds ----------

func TestPassword_LengthTooSmall(t *testing.T) {
	r := PasswordGenerate(0, true, true, true, true, "")
	assert.NotEmpty(t, r.Error)
	assert.Contains(t, r.Error, "length must be between 1 and 256")
	assert.Empty(t, r.Output)
}

func TestPassword_LengthNegative(t *testing.T) {
	r := PasswordGenerate(-5, true, true, true, true, "")
	assert.NotEmpty(t, r.Error)
	assert.Contains(t, r.Error, "length must be between 1 and 256")
	assert.Empty(t, r.Output)
}

func TestPassword_LengthTooLarge(t *testing.T) {
	r := PasswordGenerate(257, true, true, true, true, "")
	assert.NotEmpty(t, r.Error)
	assert.Contains(t, r.Error, "length must be between 1 and 256")
	assert.Empty(t, r.Output)
}

// ---------- No character types selected ----------

func TestPassword_NoCharacterTypes(t *testing.T) {
	r := PasswordGenerate(16, false, false, false, false, "")
	assert.NotEmpty(t, r.Error)
	assert.Equal(t, "at least one character type must be selected", r.Error)
	assert.Empty(t, r.Output)
}

// ---------- Randomness (two consecutive calls differ) ----------

func TestPassword_Randomness(t *testing.T) {
	r1 := PasswordGenerate(32, true, true, true, true, "")
	r2 := PasswordGenerate(32, true, true, true, true, "")
	require.Empty(t, r1.Error)
	require.Empty(t, r2.Error)
	// With 32 chars from a large charset, collision probability is negligible
	assert.NotEqual(t, r1.Output, r2.Output, "two generated passwords should differ")
}

// ---------- Tool interface compliance ----------

func TestPasswordTool_ImplementsToolInterface(t *testing.T) {
	var _ Tool = PasswordTool{}
}
