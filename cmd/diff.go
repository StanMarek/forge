package cmd

import (
	"fmt"
	"os"

	"github.com/StanMarek/forge/core/tools"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff <file-a> <file-b>",
	Short: "Compare two files and show differences",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		dataA, err := os.ReadFile(args[0])
		if err != nil {
			exitWithError(fmt.Sprintf("cannot read file %s: %s", args[0], err.Error()))
		}
		dataB, err := os.ReadFile(args[1])
		if err != nil {
			exitWithError(fmt.Sprintf("cannot read file %s: %s", args[1], err.Error()))
		}
		result := tools.DiffText(string(dataA), string(dataB))
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}
