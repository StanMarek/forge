package cmd

import (
	"fmt"
	"runtime"

	"github.com/StanMarek/forge/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("forge %s\n", version.Version)
		fmt.Printf("commit:  %s\n", version.Commit)
		fmt.Printf("built:   %s\n", version.Date)
		fmt.Printf("go:      %s\n", runtime.Version())
	},
}
