package cmd

import (
	"fmt"

	"github.com/StanMarek/forge/core/tools"
	"github.com/spf13/cobra"
)

var jwtCmd = &cobra.Command{
	Use:   "jwt",
	Short: "Decode and inspect JWT tokens",
}

var jwtDecodeCmd = &cobra.Command{
	Use:   "decode [token]",
	Short: "Decode a JWT token",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.JWTDecode(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}

		headerOnly, _ := cmd.Flags().GetBool("header-only")
		payloadOnly, _ := cmd.Flags().GetBool("payload-only")
		compact, _ := cmd.Flags().GetBool("compact")

		switch {
		case headerOnly:
			header := result.Header
			if compact {
				if m := tools.JSONMinify(header); m.Error == "" {
					header = m.Output
				}
			}
			fmt.Println(header)
		case payloadOnly:
			payload := result.Payload
			if compact {
				if m := tools.JSONMinify(payload); m.Error == "" {
					payload = m.Output
				}
			}
			fmt.Println(payload)
		case compact:
			h := result.Header
			if m := tools.JSONMinify(h); m.Error == "" {
				h = m.Output
			}
			p := result.Payload
			if m := tools.JSONMinify(p); m.Error == "" {
				p = m.Output
			}
			fmt.Printf("--- Header ---\n%s\n--- Payload ---\n%s\n--- Signature ---\n%s\n", h, p, result.Signature)
		default:
			fmt.Println(result.Output)
		}
	},
}

var jwtValidateCmd = &cobra.Command{
	Use:   "validate [token]",
	Short: "Check if a string is a structurally valid JWT",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.JWTValidate(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

func init() {
	jwtDecodeCmd.Flags().Bool("header-only", false, "Print only the header")
	jwtDecodeCmd.Flags().Bool("payload-only", false, "Print only the payload")
	jwtDecodeCmd.Flags().Bool("compact", false, "Print JSON on a single line")
	jwtCmd.AddCommand(jwtDecodeCmd)
	jwtCmd.AddCommand(jwtValidateCmd)
}
