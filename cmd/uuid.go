package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/StanMarek/forge/core/tools"
	"github.com/spf13/cobra"
)

var uuidCmd = &cobra.Command{
	Use:   "uuid",
	Short: "Generate, validate, or parse UUIDs",
}

var uuidGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new UUID",
	Run: func(cmd *cobra.Command, args []string) {
		ver, _ := cmd.Flags().GetInt("version")
		count, _ := cmd.Flags().GetInt("count")
		uppercase, _ := cmd.Flags().GetBool("uppercase")
		noHyphens, _ := cmd.Flags().GetBool("no-hyphens")

		for i := 0; i < count; i++ {
			result := tools.UUIDGenerate(ver, uppercase, noHyphens)
			if result.Error != "" {
				exitWithError(result.Error)
			}
			fmt.Println(result.Output)
		}
	},
}

var uuidValidateCmd = &cobra.Command{
	Use:   "validate [input]",
	Short: "Check if a string is a valid UUID",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.UUIDValidate(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

var uuidParseCmd = &cobra.Command{
	Use:   "parse [input]",
	Short: "Parse a UUID and show its components",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.UUIDParse(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		asJSON, _ := cmd.Flags().GetBool("json")
		if asJSON {
			out := struct {
				UUID      string `json:"uuid"`
				Version   int    `json:"version"`
				Variant   string `json:"variant"`
				Timestamp string `json:"timestamp,omitempty"`
			}{
				UUID:      result.UUID,
				Version:   result.Version,
				Variant:   result.Variant,
				Timestamp: result.Timestamp,
			}
			data, _ := json.MarshalIndent(out, "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Println(result.Output)
		}
	},
}

func init() {
	uuidGenerateCmd.Flags().Int("version", 4, "UUID version: 4 (random) or 7 (time-ordered)")
	uuidGenerateCmd.Flags().Int("count", 1, "Number of UUIDs to generate")
	uuidGenerateCmd.Flags().Bool("uppercase", false, "Output in uppercase")
	uuidGenerateCmd.Flags().Bool("no-hyphens", false, "Output without hyphens")
	uuidParseCmd.Flags().Bool("json", false, "Output as JSON")
	uuidCmd.AddCommand(uuidGenerateCmd)
	uuidCmd.AddCommand(uuidValidateCmd)
	uuidCmd.AddCommand(uuidParseCmd)
}
