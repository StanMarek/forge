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
	rootCmd.AddCommand(webCmd)
	rootCmd.AddCommand(yamlCmd)
	rootCmd.AddCommand(timestampCmd)
	rootCmd.AddCommand(numberBaseCmd)
	rootCmd.AddCommand(regexCmd)
	rootCmd.AddCommand(htmlEntityCmd)
	rootCmd.AddCommand(passwordCmd)
	rootCmd.AddCommand(loremCmd)
	rootCmd.AddCommand(colorCmd)
	rootCmd.AddCommand(cronCmd)
	rootCmd.AddCommand(textEscapeCmd)
	rootCmd.AddCommand(gzipCmd)
	rootCmd.AddCommand(textStatsCmd)
	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(xmlCmd)
	rootCmd.AddCommand(csvCmd)
	rootCmd.AddCommand(desktopCmd)
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(2)
	}
}
