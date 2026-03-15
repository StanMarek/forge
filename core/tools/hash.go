package tools

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"strings"
)

// HashTool provides metadata for the Hash Generator tool.
type HashTool struct{}

func (h HashTool) Name() string        { return "Hash Generator" }
func (h HashTool) ID() string          { return "hash" }
func (h HashTool) Description() string { return "Generate MD5, SHA1, SHA256, and SHA512 hashes" }
func (h HashTool) Category() string    { return "Generators" }
func (h HashTool) Keywords() []string {
	return []string{"hash", "md5", "sha1", "sha256", "sha512", "digest", "checksum"}
}

// DetectFromClipboard always returns false for the hash tool.
func (h HashTool) DetectFromClipboard(_ string) bool {
	return false
}

// Hash computes the hash digest of input using the specified algorithm.
// Supported algorithms: md5, sha1, sha256, sha512 (case-insensitive).
// If uppercase is true, the hex-encoded output is returned in uppercase.
func Hash(input string, algorithm string, uppercase bool) Result {
	algo := strings.ToLower(algorithm)

	var hasher hash.Hash
	switch algo {
	case "md5":
		hasher = md5.New()
	case "sha1":
		hasher = sha1.New()
	case "sha256":
		hasher = sha256.New()
	case "sha512":
		hasher = sha512.New()
	default:
		return Result{
			Error: fmt.Sprintf("unsupported algorithm: %s (supported: md5, sha1, sha256, sha512)", algorithm),
		}
	}

	hasher.Write([]byte(input))
	digest := fmt.Sprintf("%x", hasher.Sum(nil))

	if uppercase {
		digest = strings.ToUpper(digest)
	}

	return Result{Output: digest}
}
