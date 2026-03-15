package cmd

import (
	"fmt"

	"github.com/StanMarek/forge/core/tools"
	"github.com/spf13/cobra"
)

var yamlCmd = &cobra.Command{
	Use:   "yaml",
	Short: "Convert between JSON and YAML",
}

var yamlToJSONCmd = &cobra.Command{
	Use:   "to-json [input]",
	Short: "Convert YAML to JSON",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		compact, _ := cmd.Flags().GetBool("compact")
		result := tools.YAMLToJSON(input, compact)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

var yamlToYAMLCmd = &cobra.Command{
	Use:   "to-yaml [input]",
	Short: "Convert JSON to YAML",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.JSONToYAML(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

func init() {
	yamlToJSONCmd.Flags().Bool("compact", false, "Output compact JSON without indentation")
	yamlCmd.AddCommand(yamlToJSONCmd)
	yamlCmd.AddCommand(yamlToYAMLCmd)
}
