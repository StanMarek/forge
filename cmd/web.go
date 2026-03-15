package cmd

import (
	"fmt"
	"os"

	"github.com/StanMarek/forge/ui/web"
	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Launch the web server",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		host, _ := cmd.Flags().GetString("host")
		if err := web.Start(host, port); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	webCmd.Flags().Int("port", 8080, "Port to listen on")
	webCmd.Flags().String("host", "localhost", "Host to bind to")
}
