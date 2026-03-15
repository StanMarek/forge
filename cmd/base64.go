package cmd

import (
	"fmt"

	"github.com/StanMarek/forge/core/tools"
	"github.com/spf13/cobra"
)

var base64Cmd = &cobra.Command{
	Use:   "base64",
	Short: "Encode or decode Base64 strings",
}

var base64EncodeCmd = &cobra.Command{
	Use:   "encode [input]",
	Short: "Encode input as Base64",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		urlSafe, _ := cmd.Flags().GetBool("url-safe")
		noPadding, _ := cmd.Flags().GetBool("no-padding")
		result := tools.Base64Encode(input, urlSafe, noPadding)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

var base64DecodeCmd = &cobra.Command{
	Use:   "decode [input]",
	Short: "Decode a Base64-encoded string",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		urlSafe, _ := cmd.Flags().GetBool("url-safe")
		result := tools.Base64Decode(input, urlSafe)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

func init() {
	base64EncodeCmd.Flags().Bool("url-safe", false, "Use URL-safe Base64 alphabet (RFC 4648 §5)")
	base64EncodeCmd.Flags().Bool("no-padding", false, "Omit trailing = padding characters")
	base64DecodeCmd.Flags().Bool("url-safe", false, "Expect URL-safe Base64 alphabet")
	base64Cmd.AddCommand(base64EncodeCmd)
	base64Cmd.AddCommand(base64DecodeCmd)
}
