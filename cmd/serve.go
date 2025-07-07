package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"github.com/euforicio/adr-demo/internal/server"
	"github.com/spf13/cobra"
)

var (
	port int
	host string
	open bool
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start development server",
	Long: `Start a development server that renders ADR markdown files on-the-fly.

Features:
• Dynamic rendering of markdown files
• Intelligent caching with SHA-based invalidation
• Auto-open browser option
• Local development preview

Perfect for previewing ADRs during development - just edit and refresh!`,
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			fmt.Printf("🚀 Starting development server...\n")
			fmt.Printf("   Host: %s\n", host)
			fmt.Printf("   Port: %d\n", port)
		}

		config := &server.Config{
			Host:    host,
			Port:    port,
			Verbose: verbose,
		}

		srv := server.New(config)

		// Auto-open browser if requested
		if open {
			url := fmt.Sprintf("http://%s:%d", host, port)
			if verbose {
				fmt.Printf("🌐 Opening browser to %s\n", url)
			}
			openBrowser(url)
		}

		fmt.Printf("🌐 Development server running at http://%s:%d\n", host, port)
		fmt.Printf("📝 Serving dynamic ADR content (press Ctrl+C to stop)\n")

		// Start server (this blocks)
		if err := srv.Start(); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "port for the development server")
	serveCmd.Flags().StringVar(&host, "host", "localhost", "host for the development server")
	serveCmd.Flags().BoolVar(&open, "open", false, "automatically open the site in your default browser")
}

// openBrowser opens the default browser to the given URL
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
