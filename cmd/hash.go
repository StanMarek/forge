package cmd

import (
	"fmt"
	"os"

	"github.com/StanMarek/forge/core/tools"
	"github.com/StanMarek/forge/internal/stdin"
	"github.com/spf13/cobra"
)

var hashCmd = &cobra.Command{
	Use:   "hash <algorithm> [input]",
	Short: "Generate hash digests",
	Long:  "Generate hash digests. Supported algorithms: md5, sha1, sha256, sha512.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		algorithm := args[0]
		uppercase, _ := cmd.Flags().GetBool("uppercase")
		filePath, _ := cmd.Flags().GetString("file")

		var input string
		var err error

		if filePath != "" {
			data, readErr := os.ReadFile(filePath)
			if readErr != nil {
				exitWithError(fmt.Sprintf("cannot read file: %s", readErr.Error()))
			}
			input = string(data)
		} else if len(args) > 1 && args[1] != "-" {
			input = args[1]
		} else {
			input, err = stdin.Read()
			if err != nil {
				exitWithError(err.Error())
			}
		}

		result := tools.Hash(input, algorithm, uppercase)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

func init() {
	hashCmd.Flags().Bool("uppercase", false, "Output hash in uppercase hex")
	hashCmd.Flags().String("file", "", "Hash the contents of a file instead of a string")
}
