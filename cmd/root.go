package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "adr-gen",
	Short: "A modern static site generator for Architecture Decision Records",
	Long: `adr-gen is a fast, modern static site generator specifically designed for
Architecture Decision Records (ADRs). It provides:

• Beautiful, responsive web interface
• GitHub Flavored Markdown support  
• Mermaid diagram rendering
• Development server with dynamic rendering
• CSS View Transitions for smooth navigation
• SEO-optimized static HTML output

Perfect for teams who want to maintain architectural decisions in a beautiful,
accessible format that works on any device.`,
	Version: "1.0.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.adr-gen.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Version flag
	rootCmd.Flags().BoolP("version", "", false, "version for adr-gen")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if verbose {
		fmt.Println("adr-gen starting...")
	}
}
