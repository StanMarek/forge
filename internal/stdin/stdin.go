package stdin

import (
	"errors"
	"io"
	"os"
	"strings"
)

// Read reads all input from stdin if it is piped.
// Returns an error if stdin is a terminal (interactive).
func Read() (string, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return "", errors.New("no input provided")
	}
	if info.Mode()&os.ModeCharDevice != 0 {
		return "", errors.New("no input provided")
	}
	return readFrom(os.Stdin)
}

// readFrom reads all content from a reader, trimming the trailing newline.
func readFrom(r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(data), "\n"), nil
}
