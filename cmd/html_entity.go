package cmd

import (
	"fmt"

	"github.com/StanMarek/forge/core/tools"
	"github.com/spf13/cobra"
)

var htmlEntityCmd = &cobra.Command{
	Use:   "html-entity",
	Short: "Encode or decode HTML entities",
}

var htmlEntityEncodeCmd = &cobra.Command{
	Use:   "encode [input]",
	Short: "Encode special HTML characters as entities",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.HTMLEntityEncode(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

var htmlEntityDecodeCmd = &cobra.Command{
	Use:   "decode [input]",
	Short: "Decode HTML entities back to characters",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.HTMLEntityDecode(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

func init() {
	htmlEntityCmd.AddCommand(htmlEntityEncodeCmd)
	htmlEntityCmd.AddCommand(htmlEntityDecodeCmd)
}
