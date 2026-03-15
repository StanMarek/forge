package cmd

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/StanMarek/forge/ui/web"
	"github.com/spf13/cobra"
)

var desktopCmd = &cobra.Command{
	Use:   "desktop",
	Short: "Launch Forge in a desktop window",
	Run: func(cmd *cobra.Command, args []string) {
		// Find a random available port
		listener, err := net.Listen("tcp", "localhost:0")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error finding available port: %v\n", err)
			os.Exit(1)
		}
		port := listener.Addr().(*net.TCPAddr).Port
		listener.Close()

		url := fmt.Sprintf("http://localhost:%d", port)

		// Start web server in background
		go func() {
			if err := web.Start("localhost", port); err != nil {
				fmt.Fprintf(os.Stderr, "server error: %v\n", err)
				os.Exit(1)
			}
		}()

		fmt.Printf("Forge desktop running at %s\n", url)

		// Try to open in app mode (Chrome/Edge) for native-looking window
		if !tryAppMode(url) {
			// Fall back to system browser
			openBrowser(url)
		}

		// Wait for interrupt
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		fmt.Println("\nShutting down...")
	},
}

// tryAppMode attempts to open the URL in Chrome/Edge app mode (borderless window).
func tryAppMode(url string) bool {
	browsers := chromePaths()
	for _, browser := range browsers {
		path, err := exec.LookPath(browser)
		if err != nil {
			continue
		}
		cmd := exec.Command(path, "--app="+url, "--new-window")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err == nil {
			return true
		}
	}
	return false
}

// chromePaths returns possible Chrome/Edge binary names for the current OS.
func chromePaths() []string {
	switch runtime.GOOS {
	case "darwin":
		return []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
			"/Applications/Brave Browser.app/Contents/MacOS/Brave Browser",
		}
	case "windows":
		return []string{
			"chrome",
			"msedge",
			"chromium",
			"brave",
		}
	default: // linux
		return []string{
			"google-chrome",
			"google-chrome-stable",
			"chromium",
			"chromium-browser",
			"microsoft-edge",
			"brave-browser",
		}
	}
}

// openBrowser opens the URL in the system default browser.
func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	cmd.Start()
}
