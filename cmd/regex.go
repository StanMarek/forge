package cmd

import (
	"fmt"

	"github.com/StanMarek/forge/core/tools"
	"github.com/spf13/cobra"
)

var regexCmd = &cobra.Command{
	Use:   "regex <pattern> [test-string]",
	Short: "Test a regular expression against input",
	Long:  "Test a regular expression. Pattern is the first argument, test string is the second argument or stdin.",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		pattern := args[0]
		var testString string
		if len(args) > 1 && args[1] != "-" {
			testString = args[1]
		} else {
			input, err := resolveInput(nil)
			if err != nil {
				exitWithError(err.Error())
			}
			testString = input
		}
		global, _ := cmd.Flags().GetBool("global")
		result := tools.RegexTest(pattern, testString, global)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

func init() {
	regexCmd.Flags().BoolP("global", "g", false, "Return all matches instead of just the first")
}
