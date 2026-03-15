package cmd

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/StanMarek/forge/internal/version"
	"github.com/StanMarek/forge/ui/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "forge",
	Short: "A developer's workbench for the terminal, browser, and desktop",
	Run: func(cmd *cobra.Command, args []string) {
		app := tui.New()
		p := tea.NewProgram(app)
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	},
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
