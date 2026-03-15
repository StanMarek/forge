package cmd

import (
	"fmt"

	"github.com/StanMarek/forge/core/tools"
	"github.com/spf13/cobra"
)

var loremCmd = &cobra.Command{
	Use:   "lorem",
	Short: "Generate lorem ipsum placeholder text",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		words, _ := cmd.Flags().GetInt("words")
		sentences, _ := cmd.Flags().GetInt("sentences")
		paragraphs, _ := cmd.Flags().GetInt("paragraphs")

		// Default: 1 paragraph if no mode specified.
		if words == 0 && sentences == 0 && paragraphs == 0 {
			paragraphs = 1
		}

		result := tools.LoremGenerate(words, sentences, paragraphs)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

func init() {
	loremCmd.Flags().Int("words", 0, "Generate N words")
	loremCmd.Flags().Int("sentences", 0, "Generate N sentences")
	loremCmd.Flags().Int("paragraphs", 0, "Generate N paragraphs (default: 1)")
}
