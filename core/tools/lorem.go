package tools

import (
	"crypto/rand"
	"fmt"
	"math/big"
	mrand "math/rand"
	"strings"
	"unicode"
)

// LoremTool provides metadata for the Lorem Ipsum Generator tool.
type LoremTool struct{}

func (l LoremTool) Name() string        { return "Lorem Ipsum Generator" }
func (l LoremTool) ID() string          { return "lorem" }
func (l LoremTool) Description() string { return "Generate lorem ipsum placeholder text" }
func (l LoremTool) Category() string    { return "Generators" }
func (l LoremTool) Keywords() []string {
	return []string{"lorem", "ipsum", "placeholder", "text"}
}

// DetectFromClipboard always returns false for the lorem tool (generative tool).
func (l LoremTool) DetectFromClipboard(_ string) bool {
	return false
}

// loremWords is the classic lorem ipsum word list (100+ words).
var loremWords = []string{
	"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit",
	"sed", "do", "eiusmod", "tempor", "incididunt", "ut", "labore", "et", "dolore",
	"magna", "aliqua", "enim", "ad", "minim", "veniam", "quis", "nostrud",
	"exercitation", "ullamco", "laboris", "nisi", "aliquip", "ex", "ea", "commodo",
	"consequat", "duis", "aute", "irure", "in", "reprehenderit", "voluptate",
	"velit", "esse", "cillum", "fugiat", "nulla", "pariatur", "excepteur", "sint",
	"occaecat", "cupidatat", "non", "proident", "sunt", "culpa", "qui", "officia",
	"deserunt", "mollit", "anim", "id", "est", "laborum", "at", "vero", "eos",
	"accusamus", "iusto", "odio", "dignissimos", "ducimus", "blanditiis",
	"praesentium", "voluptatum", "deleniti", "atque", "corrupti", "quos", "dolores",
	"quas", "molestias", "excepturi", "obcaecati", "cupiditate", "provident",
	"similique", "mollitia", "animi", "sapiente", "delectus", "rerum", "hic",
	"tenetur", "a", "eligendi", "optio", "cumque", "nihil", "impedit", "quo",
	"minus", "maxime", "placeat", "facere", "possimus", "omnis", "voluptas",
	"assumenda", "repellendus", "temporibus", "autem", "quibusdam", "officiis",
	"debitis", "aut", "necessitatibus", "saepe", "eveniet", "voluptates",
	"repudiandae", "recusandae",
}

// newLoremRand creates a math/rand source seeded from crypto/rand.
func newLoremRand() *mrand.Rand {
	seed, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		// Fallback — should never happen in practice.
		return mrand.New(mrand.NewSource(42))
	}
	return mrand.New(mrand.NewSource(seed.Int64()))
}

// capitalize returns s with the first rune uppercased.
func capitalize(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// LoremGenerate generates lorem ipsum placeholder text.
// Exactly one of words, sentences, or paragraphs must be > 0.
func LoremGenerate(words int, sentences int, paragraphs int) Result {
	// Validate exactly one mode is active.
	modes := 0
	if words > 0 {
		modes++
	}
	if sentences > 0 {
		modes++
	}
	if paragraphs > 0 {
		modes++
	}
	if modes == 0 {
		return Result{Error: "one of words, sentences, or paragraphs must be greater than 0"}
	}
	if modes > 1 {
		return Result{Error: "only one of words, sentences, or paragraphs may be greater than 0"}
	}

	r := newLoremRand()

	if words > 0 {
		return Result{Output: generateWords(r, words)}
	}
	if sentences > 0 {
		return Result{Output: generateSentences(r, sentences)}
	}
	return Result{Output: generateParagraphs(r, paragraphs)}
}

// generateWords picks n words from the word list, cycling if needed.
func generateWords(r *mrand.Rand, n int) string {
	result := make([]string, n)
	for i := 0; i < n; i++ {
		result[i] = loremWords[r.Intn(len(loremWords))]
	}
	return strings.Join(result, " ")
}

// generateOneSentence generates a single sentence of 8-15 words.
func generateOneSentence(r *mrand.Rand) string {
	count := 8 + r.Intn(8) // 8 to 15 inclusive
	words := make([]string, count)
	for i := 0; i < count; i++ {
		words[i] = loremWords[r.Intn(len(loremWords))]
	}
	words[0] = capitalize(words[0])
	return strings.Join(words, " ") + "."
}

// generateSentences generates n sentences.
func generateSentences(r *mrand.Rand, n int) string {
	sents := make([]string, n)
	for i := 0; i < n; i++ {
		sents[i] = generateOneSentence(r)
	}
	return strings.Join(sents, " ")
}

// generateOneParagraph generates a single paragraph of 4-7 sentences.
func generateOneParagraph(r *mrand.Rand) string {
	count := 4 + r.Intn(4) // 4 to 7 inclusive
	sents := make([]string, count)
	for i := 0; i < count; i++ {
		sents[i] = generateOneSentence(r)
	}
	return strings.Join(sents, " ")
}

// generateParagraphs generates n paragraphs separated by double newlines.
func generateParagraphs(r *mrand.Rand, n int) string {
	paras := make([]string, n)
	for i := 0; i < n; i++ {
		paras[i] = generateOneParagraph(r)
	}
	return strings.Join(paras, fmt.Sprintf("\n\n"))
}
