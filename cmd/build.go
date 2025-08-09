package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/euforicio/adr-demo/internal/config"
	"github.com/euforicio/adr-demo/internal/generator"
	"github.com/spf13/cobra"
)

var (
	configFile string
	outputDir  string
	baseURL    string
	minify     bool
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the static site",
	Long: `Build generates a complete static website from your ADR markdown files.
	
The generated site includes:
• Individual pages for each ADR
• Navigation sidebar with search
• Responsive design that works on all devices
• Optimized assets and SEO meta tags
• Mermaid diagrams rendered as SVG

Output is generated in the docs/ directory by default, ready for GitHub Pages deployment.`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()

		// Load configuration
		cfg, err := config.LoadConfig(configFile)
		if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}

		// Override config with command line flags if provided
		if outputDir != "" {
			cfg.OutputDirectory = outputDir
		}
		if baseURL != "" {
			cfg.BaseURL = baseURL
		}
		if minify {
			cfg.Minify = minify
		}
		if verbose {
			cfg.Verbose = verbose
		}

		if cfg.Verbose {
			fmt.Printf("🔧 Building static site...\n")
			fmt.Printf("   ADR Directory: %s\n", cfg.ADRDirectory)
			fmt.Printf("   Output: %s\n", cfg.OutputDirectory)
			if cfg.BaseURL != "" {
				fmt.Printf("   Base URL: %s\n", cfg.BaseURL)
			}
			if cfg.Minify {
				fmt.Printf("   Minification: enabled\n")
			}
			fmt.Printf("   Default Category: %s\n", cfg.DefaultCategory)
			fmt.Printf("   Allowed Categories: %v\n", cfg.AllowedCategories)
		}

		gen := generator.New(cfg)
		if err := gen.Build(); err != nil {
			log.Fatalf("Build failed: %v", err)
		}

		duration := time.Since(start)
		fmt.Printf("✅ Site built successfully in %s (%.2fs)\n", cfg.OutputDirectory, duration.Seconds())

		if cfg.Verbose {
			stats := gen.GetStats()
			fmt.Printf("📊 Build stats:\n")
			fmt.Printf("   • %d ADRs processed\n", stats.ADRCount)
			fmt.Printf("   • %d pages generated\n", stats.PageCount)
			fmt.Printf("   • %d assets copied\n", stats.AssetCount)
			if stats.DiagramCount > 0 {
				fmt.Printf("   • %d Mermaid diagrams rendered\n", stats.DiagramCount)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringVarP(&configFile, "config", "c", "", "config file path (default: look for adr-config.yaml)")
	buildCmd.Flags().StringVarP(&outputDir, "output", "o", "", "output directory for generated site (overrides config)")
	buildCmd.Flags().StringVar(&baseURL, "base-url", "", "base URL for the site (overrides config)")
	buildCmd.Flags().BoolVar(&minify, "minify", false, "minify HTML, CSS, and JavaScript (overrides config)")
}
