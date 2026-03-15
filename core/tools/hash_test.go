package tools

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------- Tool metadata ----------

func TestHashTool_Metadata(t *testing.T) {
	tool := HashTool{}

	assert.Equal(t, "Hash Generator", tool.Name())
	assert.Equal(t, "hash", tool.ID())
	assert.Equal(t, "Generators", tool.Category())
	assert.NotEmpty(t, tool.Description())

	keywords := tool.Keywords()
	for _, kw := range []string{"hash", "md5", "sha1", "sha256", "sha512", "digest", "checksum"} {
		assert.Contains(t, keywords, kw)
	}
}

func TestHashTool_DetectFromClipboard(t *testing.T) {
	tool := HashTool{}
	assert.False(t, tool.DetectFromClipboard("anything"))
	assert.False(t, tool.DetectFromClipboard(""))
}

// ---------- Known test vectors ----------

func TestHash_MD5(t *testing.T) {
	r := Hash("hello world", "md5", false)
	require.Empty(t, r.Error)
	assert.Equal(t, "5eb63bbbe01eeed093cb22bb8f5acdc3", r.Output)
}

func TestHash_SHA1(t *testing.T) {
	r := Hash("hello world", "sha1", false)
	require.Empty(t, r.Error)
	assert.Equal(t, "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed", r.Output)
}

func TestHash_SHA256(t *testing.T) {
	r := Hash("hello world", "sha256", false)
	require.Empty(t, r.Error)
	assert.Equal(t, "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9", r.Output)
}

func TestHash_SHA512(t *testing.T) {
	r := Hash("hello world", "sha512", false)
	require.Empty(t, r.Error)
	assert.Equal(t, "309ecc489c12d6eb4cc40f50c902f2b4d0ed77ee511a7c7a9bcd3ca86d4cd86f989dd35bc5ff499670da34255b45b0cfd830e81f605dcf7dc5542e93ae9cd76f", r.Output)
}

func TestHash_SHA256_Empty(t *testing.T) {
	r := Hash("", "sha256", false)
	require.Empty(t, r.Error)
	assert.Equal(t, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", r.Output)
}

// ---------- Uppercase ----------

func TestHash_Uppercase(t *testing.T) {
	r := Hash("hello world", "md5", true)
	require.Empty(t, r.Error)
	assert.Equal(t, "5EB63BBBE01EEED093CB22BB8F5ACDC3", r.Output)
}

func TestHash_SHA256_Uppercase(t *testing.T) {
	r := Hash("hello world", "sha256", true)
	require.Empty(t, r.Error)
	assert.Equal(t, strings.ToUpper("b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"), r.Output)
}

// ---------- Case-insensitive algorithm name ----------

func TestHash_AlgorithmCaseInsensitive(t *testing.T) {
	for _, algo := range []string{"MD5", "Md5", "mD5", "md5"} {
		r := Hash("hello world", algo, false)
		require.Empty(t, r.Error, "algorithm %q should be accepted", algo)
		assert.Equal(t, "5eb63bbbe01eeed093cb22bb8f5acdc3", r.Output)
	}

	for _, algo := range []string{"SHA256", "Sha256", "sha256"} {
		r := Hash("hello world", algo, false)
		require.Empty(t, r.Error, "algorithm %q should be accepted", algo)
		assert.Equal(t, "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9", r.Output)
	}
}

// ---------- Unsupported algorithm ----------

func TestHash_UnsupportedAlgorithm(t *testing.T) {
	r := Hash("hello", "sha384", false)
	assert.NotEmpty(t, r.Error)
	assert.Contains(t, r.Error, "unsupported algorithm: sha384")
	assert.Contains(t, r.Error, "supported: md5, sha1, sha256, sha512")
	assert.Empty(t, r.Output)
}

func TestHash_UnsupportedAlgorithmPreservesOriginalCase(t *testing.T) {
	r := Hash("hello", "Blake2b", false)
	assert.Contains(t, r.Error, "unsupported algorithm: Blake2b")
}

// ---------- Tool interface compliance ----------

func TestHashTool_ImplementsToolInterface(t *testing.T) {
	var _ Tool = HashTool{}
}
