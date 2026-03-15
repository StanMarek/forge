package cmd

import (
	"fmt"

	"github.com/StanMarek/forge/core/tools"
	"github.com/spf13/cobra"
)

var textStatsCmd = &cobra.Command{
	Use:   "text-stats [input]",
	Short: "Analyze text statistics or convert text case",
	Long:  "Analyze text statistics or convert case. Mode: stats (default), lower, upper, title, camel, snake, kebab.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		mode, _ := cmd.Flags().GetString("mode")
		if mode == "" {
			mode = "stats"
		}

		if mode == "stats" {
			result := tools.TextStats(input)
			if result.Error != "" {
				exitWithError(result.Error)
			}
			fmt.Println(result.Output)
		} else {
			result := tools.TextCaseConvert(input, mode)
			if result.Error != "" {
				exitWithError(result.Error)
			}
			fmt.Println(result.Output)
		}
	},
}

func init() {
	textStatsCmd.Flags().String("mode", "stats", "Operation mode: stats, lower, upper, title, camel, snake, kebab")
}
