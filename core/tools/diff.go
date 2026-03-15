package tools

import (
	"fmt"
	"strings"
)

// DiffTool provides metadata for the Text Diff tool.
type DiffTool struct{}

func (DiffTool) Name() string        { return "Text Diff" }
func (DiffTool) ID() string          { return "diff" }
func (DiffTool) Description() string { return "Compare two texts and show differences" }
func (DiffTool) Category() string    { return "Text" }
func (DiffTool) Keywords() []string  { return []string{"diff", "compare", "text", "difference"} }

// DetectFromClipboard always returns false because diff requires two inputs.
func (DiffTool) DetectFromClipboard(_ string) bool {
	return false
}

// DiffText computes a unified diff between textA and textB using a longest
// common subsequence (LCS) algorithm. Lines present only in A are prefixed
// with "-", lines only in B with "+", and common lines with " ".
func DiffText(textA, textB string) Result {
	if textA == textB {
		return Result{Output: "Texts are identical"}
	}

	linesA := splitLines(textA)
	linesB := splitLines(textB)

	lcs := computeLCS(linesA, linesB)

	var buf strings.Builder
	buf.WriteString("--- Text A\n")
	buf.WriteString("+++ Text B\n")

	ia, ib, il := 0, 0, 0
	for il < len(lcs) {
		// Emit lines from A that are not in the LCS (removed).
		for ia < len(linesA) && linesA[ia] != lcs[il] {
			fmt.Fprintf(&buf, "-%s\n", linesA[ia])
			ia++
		}
		// Emit lines from B that are not in the LCS (added).
		for ib < len(linesB) && linesB[ib] != lcs[il] {
			fmt.Fprintf(&buf, "+%s\n", linesB[ib])
			ib++
		}
		// Emit the common line.
		fmt.Fprintf(&buf, " %s\n", lcs[il])
		ia++
		ib++
		il++
	}
	// Remaining lines in A after LCS exhausted (removed).
	for ia < len(linesA) {
		fmt.Fprintf(&buf, "-%s\n", linesA[ia])
		ia++
	}
	// Remaining lines in B after LCS exhausted (added).
	for ib < len(linesB) {
		fmt.Fprintf(&buf, "+%s\n", linesB[ib])
		ib++
	}

	return Result{Output: strings.TrimRight(buf.String(), "\n")}
}

// splitLines splits text into lines. An empty string produces a single
// empty-string element so that the diff can represent it.
func splitLines(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, "\n")
}

// computeLCS returns the longest common subsequence of two string slices
// using the classic dynamic-programming approach.
func computeLCS(a, b []string) []string {
	m, n := len(a), len(b)
	// Build the DP table.
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else if dp[i-1][j] >= dp[i][j-1] {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = dp[i][j-1]
			}
		}
	}
	// Backtrack to find the LCS.
	lcs := make([]string, 0, dp[m][n])
	i, j := m, n
	for i > 0 && j > 0 {
		if a[i-1] == b[j-1] {
			lcs = append(lcs, a[i-1])
			i--
			j--
		} else if dp[i-1][j] >= dp[i][j-1] {
			i--
		} else {
			j--
		}
	}
	// Reverse the LCS (built backwards).
	for left, right := 0, len(lcs)-1; left < right; left, right = left+1, right-1 {
		lcs[left], lcs[right] = lcs[right], lcs[left]
	}
	return lcs
}
