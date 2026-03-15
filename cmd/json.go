package cmd

import (
	"fmt"

	"github.com/StanMarek/forge/core/tools"
	"github.com/spf13/cobra"
)

var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "Format, minify, or validate JSON",
}

var jsonFormatCmd = &cobra.Command{
	Use:   "format [input]",
	Short: "Pretty-print JSON with indentation",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		indent, _ := cmd.Flags().GetInt("indent")
		tabs, _ := cmd.Flags().GetBool("tabs")
		sortKeys, _ := cmd.Flags().GetBool("sort-keys")
		result := tools.JSONFormat(input, indent, sortKeys, tabs)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

var jsonMinifyCmd = &cobra.Command{
	Use:   "minify [input]",
	Short: "Remove all whitespace from JSON",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.JSONMinify(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

var jsonValidateCmd = &cobra.Command{
	Use:   "validate [input]",
	Short: "Check if input is valid JSON",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.JSONValidate(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

func init() {
	jsonFormatCmd.Flags().Int("indent", 2, "Number of spaces for indentation")
	jsonFormatCmd.Flags().Bool("tabs", false, "Use tabs instead of spaces")
	jsonFormatCmd.Flags().Bool("sort-keys", false, "Sort object keys alphabetically")
	jsonCmd.AddCommand(jsonFormatCmd)
	jsonCmd.AddCommand(jsonMinifyCmd)
	jsonCmd.AddCommand(jsonValidateCmd)
}
