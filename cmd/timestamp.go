package cmd

import (
	"fmt"

	"github.com/StanMarek/forge/core/tools"
	"github.com/spf13/cobra"
)

var timestampCmd = &cobra.Command{
	Use:   "timestamp",
	Short: "Convert between Unix timestamps and human-readable dates",
}

var timestampFromUnixCmd = &cobra.Command{
	Use:   "from-unix [input]",
	Short: "Convert a Unix timestamp to a human-readable date",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		tz, _ := cmd.Flags().GetString("tz")
		result := tools.TimestampFromUnix(input, tz)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

var timestampToUnixCmd = &cobra.Command{
	Use:   "to-unix [input]",
	Short: "Convert a date string to a Unix timestamp",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		millis, _ := cmd.Flags().GetBool("millis")
		result := tools.TimestampToUnix(input, millis)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

var timestampNowCmd = &cobra.Command{
	Use:   "now",
	Short: "Show the current time in multiple formats",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		tz, _ := cmd.Flags().GetString("tz")
		result := tools.TimestampNow(tz)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

func init() {
	timestampFromUnixCmd.Flags().String("tz", "UTC", "Timezone for output (e.g. America/New_York)")
	timestampToUnixCmd.Flags().Bool("millis", false, "Output timestamp in milliseconds")
	timestampNowCmd.Flags().String("tz", "UTC", "Timezone for output (e.g. America/New_York)")
	timestampCmd.AddCommand(timestampFromUnixCmd)
	timestampCmd.AddCommand(timestampToUnixCmd)
	timestampCmd.AddCommand(timestampNowCmd)
}
