package cmd

import (
	"os"

	"github.com/StanMarek/forge/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "forge",
	Short: "A developer's workbench for the terminal, browser, and desktop",
	// TODO: launch TUI when no subcommand is given
}

func init() {
	rootCmd.Version = version.Version
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(base64Cmd)
	rootCmd.AddCommand(jwtCmd)
	rootCmd.AddCommand(jsonCmd)
	rootCmd.AddCommand(hashCmd)
	rootCmd.AddCommand(urlCmd)
	rootCmd.AddCommand(uuidCmd)
	rootCmd.AddCommand(tuiCmd)
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(2)
	}
}
