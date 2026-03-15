package cmd

import (
	"fmt"

	"github.com/StanMarek/forge/core/tools"
	"github.com/spf13/cobra"
)

var colorCmd = &cobra.Command{
	Use:   "color [input]",
	Short: "Convert colors between hex, RGB, and HSL formats",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.ColorConvert(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}
