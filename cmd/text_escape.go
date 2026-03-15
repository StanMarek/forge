package cmd

import (
	"fmt"

	"github.com/StanMarek/forge/core/tools"
	"github.com/spf13/cobra"
)

var textEscapeCmd = &cobra.Command{
	Use:   "text-escape",
	Short: "Escape or unescape special characters in text",
}

var textEscapeEscapeCmd = &cobra.Command{
	Use:   "escape [input]",
	Short: "Escape special characters in text",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.TextEscape(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

var textEscapeUnescapeCmd = &cobra.Command{
	Use:   "unescape [input]",
	Short: "Unescape escape sequences in text",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.TextUnescape(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

func init() {
	textEscapeCmd.AddCommand(textEscapeEscapeCmd)
	textEscapeCmd.AddCommand(textEscapeUnescapeCmd)
}
