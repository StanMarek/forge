package cmd

import (
	"fmt"

	"github.com/StanMarek/forge/core/tools"
	"github.com/spf13/cobra"
)

var csvCmd = &cobra.Command{
	Use:   "csv",
	Short: "Convert between JSON and CSV",
}

var csvToCSVCmd = &cobra.Command{
	Use:   "to-csv [input]",
	Short: "Convert a JSON array of objects to CSV",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		delimiter, _ := cmd.Flags().GetString("delimiter")
		result := tools.JSONToCSV(input, delimiter)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

var csvToJSONCmd = &cobra.Command{
	Use:   "to-json [input]",
	Short: "Convert CSV to a JSON array of objects",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		delimiter, _ := cmd.Flags().GetString("delimiter")
		result := tools.CSVToJSON(input, delimiter)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

func init() {
	csvToCSVCmd.Flags().String("delimiter", ",", "Field delimiter character")
	csvToJSONCmd.Flags().String("delimiter", ",", "Field delimiter character")
	csvCmd.AddCommand(csvToCSVCmd)
	csvCmd.AddCommand(csvToJSONCmd)
}
