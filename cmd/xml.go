package cmd

import (
	"fmt"

	"github.com/StanMarek/forge/core/tools"
	"github.com/spf13/cobra"
)

var xmlCmd = &cobra.Command{
	Use:   "xml",
	Short: "Format or minify XML",
}

var xmlFormatCmd = &cobra.Command{
	Use:   "format [input]",
	Short: "Pretty-print XML with indentation",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.XMLFormat(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

var xmlMinifyCmd = &cobra.Command{
	Use:   "minify [input]",
	Short: "Minify XML by removing unnecessary whitespace",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.XMLMinify(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

func init() {
	xmlCmd.AddCommand(xmlFormatCmd)
	xmlCmd.AddCommand(xmlMinifyCmd)
}
