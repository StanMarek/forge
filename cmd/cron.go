package cmd

import (
	"fmt"

	"github.com/StanMarek/forge/core/tools"
	"github.com/spf13/cobra"
)

var cronCmd = &cobra.Command{
	Use:   "cron [expression]",
	Short: "Parse cron expressions into human-readable descriptions",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.CronParse(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}
