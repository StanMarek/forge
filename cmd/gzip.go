package cmd

import (
	"fmt"

	"github.com/StanMarek/forge/core/tools"
	"github.com/spf13/cobra"
)

var gzipCmd = &cobra.Command{
	Use:   "gzip",
	Short: "Compress or decompress data using gzip",
}

var gzipCompressCmd = &cobra.Command{
	Use:   "compress [input]",
	Short: "Compress input with gzip (output is base64-encoded)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.GZipCompress(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

var gzipDecompressCmd = &cobra.Command{
	Use:   "decompress [input]",
	Short: "Decompress base64-encoded gzip data",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.GZipDecompress(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

func init() {
	gzipCmd.AddCommand(gzipCompressCmd)
	gzipCmd.AddCommand(gzipDecompressCmd)
}
