package cmd

import (
	"fmt"

	"github.com/StanMarek/forge/core/tools"
	"github.com/spf13/cobra"
)

var passwordCmd = &cobra.Command{
	Use:   "password",
	Short: "Generate cryptographically secure random passwords",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		length, _ := cmd.Flags().GetInt("length")
		count, _ := cmd.Flags().GetInt("count")
		noUppercase, _ := cmd.Flags().GetBool("no-uppercase")
		noLowercase, _ := cmd.Flags().GetBool("no-lowercase")
		noDigits, _ := cmd.Flags().GetBool("no-digits")
		noSymbols, _ := cmd.Flags().GetBool("no-symbols")
		symbolSet, _ := cmd.Flags().GetString("symbols")

		for i := 0; i < count; i++ {
			result := tools.PasswordGenerate(
				length,
				!noUppercase,
				!noLowercase,
				!noDigits,
				!noSymbols,
				symbolSet,
			)
			if result.Error != "" {
				exitWithError(result.Error)
			}
			fmt.Println(result.Output)
		}
	},
}

func init() {
	passwordCmd.Flags().Int("length", 16, "Password length")
	passwordCmd.Flags().Int("count", 1, "Number of passwords to generate")
	passwordCmd.Flags().Bool("no-uppercase", false, "Exclude uppercase letters")
	passwordCmd.Flags().Bool("no-lowercase", false, "Exclude lowercase letters")
	passwordCmd.Flags().Bool("no-digits", false, "Exclude digits")
	passwordCmd.Flags().Bool("no-symbols", false, "Exclude symbols")
	passwordCmd.Flags().String("symbols", "", `Custom symbol set (default "!@#$%^&*()-_=+")`)
}
