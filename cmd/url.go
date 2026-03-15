package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/StanMarek/forge/core/tools"
	"github.com/spf13/cobra"
)

var urlCmd = &cobra.Command{
	Use:   "url",
	Short: "Encode, decode, or parse URLs",
}

var urlEncodeCmd = &cobra.Command{
	Use:   "encode [input]",
	Short: "URL-encode a string",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		component, _ := cmd.Flags().GetBool("component")
		result := tools.URLEncode(input, component)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

var urlDecodeCmd = &cobra.Command{
	Use:   "decode [input]",
	Short: "Decode a URL-encoded string",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.URLDecode(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		fmt.Println(result.Output)
	},
}

var urlParseCmd = &cobra.Command{
	Use:   "parse [input]",
	Short: "Parse a URL into components",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := resolveInput(args)
		if err != nil {
			exitWithError(err.Error())
		}
		result := tools.URLParse(input)
		if result.Error != "" {
			exitWithError(result.Error)
		}
		asJSON, _ := cmd.Flags().GetBool("json")
		if asJSON {
			out := struct {
				Scheme   string              `json:"scheme"`
				Host     string              `json:"host"`
				Port     string              `json:"port,omitempty"`
				Path     string              `json:"path,omitempty"`
				Query    string              `json:"query,omitempty"`
				Fragment string              `json:"fragment,omitempty"`
				Params   map[string][]string `json:"params,omitempty"`
			}{
				Scheme:   result.Scheme,
				Host:     result.Host,
				Port:     result.Port,
				Path:     result.Path,
				Query:    result.Query,
				Fragment: result.Fragment,
				Params:   result.Params,
			}
			data, _ := json.MarshalIndent(out, "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Println(result.Output)
		}
	},
}

func init() {
	urlEncodeCmd.Flags().Bool("component", false, "Encode as URL component (encodes /, ?, &, =)")
	urlParseCmd.Flags().Bool("json", false, "Output as JSON")
	urlCmd.AddCommand(urlEncodeCmd)
	urlCmd.AddCommand(urlDecodeCmd)
	urlCmd.AddCommand(urlParseCmd)
}
