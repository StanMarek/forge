package tools

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const defaultSymbolSet = `!@#$%^&*()-_=+`

const (
	uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	digitChars     = "0123456789"
)

// PasswordTool provides metadata for the Password Generator tool.
type PasswordTool struct{}

func (p PasswordTool) Name() string        { return "Password Generator" }
func (p PasswordTool) ID() string          { return "password" }
func (p PasswordTool) Description() string { return "Generate cryptographically secure random passwords" }
func (p PasswordTool) Category() string    { return "Generators" }
func (p PasswordTool) Keywords() []string {
	return []string{"password", "generate", "random", "secure"}
}

// DetectFromClipboard always returns false for the password tool.
// Password generation is a generative tool with no meaningful clipboard detection.
func (p PasswordTool) DetectFromClipboard(_ string) bool {
	return false
}

// PasswordGenerate creates a random password of the specified length using the
// selected character types. It uses crypto/rand for secure randomness.
// If symbolSet is empty, the default symbol set is used.
func PasswordGenerate(length int, uppercase bool, lowercase bool, digits bool, symbols bool, symbolSet string) Result {
	if length < 1 || length > 256 {
		return Result{Error: fmt.Sprintf("length must be between 1 and 256, got %d", length)}
	}

	if !uppercase && !lowercase && !digits && !symbols {
		return Result{Error: "at least one character type must be selected"}
	}

	if symbolSet == "" {
		symbolSet = defaultSymbolSet
	}

	var charset strings.Builder
	if uppercase {
		charset.WriteString(uppercaseChars)
	}
	if lowercase {
		charset.WriteString(lowercaseChars)
	}
	if digits {
		charset.WriteString(digitChars)
	}
	if symbols {
		charset.WriteString(symbolSet)
	}

	chars := charset.String()
	charsetLen := big.NewInt(int64(len(chars)))

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		idx, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return Result{Error: fmt.Sprintf("crypto/rand failed: %s", err.Error())}
		}
		result[i] = chars[idx.Int64()]
	}

	return Result{Output: string(result)}
}
