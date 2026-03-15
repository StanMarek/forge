package cmd

import (
	"fmt"
	"os"

	"github.com/StanMarek/forge/internal/stdin"
)

// resolveInput gets input from args or stdin.
// If args[0] exists and is not "-", use it. Otherwise read stdin.
func resolveInput(args []string) (string, error) {
	if len(args) > 0 && args[0] != "-" {
		return args[0], nil
	}
	return stdin.Read()
}

// exitWithError prints an error to stderr and exits with code 1.
func exitWithError(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
