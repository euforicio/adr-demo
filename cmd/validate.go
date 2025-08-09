package cmd

import (
	"fmt"
	"log"

	"github.com/euforicio/adr-demo/internal/validator"
	"github.com/spf13/cobra"
)

var (
	strict bool
	fix    bool
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate ADR structure and content",
	Long: `Validate all ADR files for:

Structure:
• Correct filename format (NNNN-kebab-case-title.md)
• Sequential numbering without gaps
• Required sections (Status, Context, Decision, Consequences)

Content:
• Valid front matter (if present)
• Proper heading hierarchy
• Valid Mermaid diagram syntax
• Working internal links
• Consistent status values

Use --strict for additional style checks and --fix to automatically
correct common issues.`,
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			fmt.Printf("🔍 Validating ADR files...\n")
			if strict {
				fmt.Printf("   Mode: strict validation\n")
			}
			if fix {
				fmt.Printf("   Auto-fix: enabled\n")
			}
		}

		config := &validator.Config{
			Strict:  strict,
			Fix:     fix,
			Verbose: verbose,
		}

		v := validator.New(config)
		result, err := v.ValidateAll()
		if err != nil {
			log.Fatalf("Validation failed: %v", err)
		}

		// Print results
		if result.HasErrors() {
			fmt.Printf("❌ Validation failed (%d errors, %d warnings)\n",
				result.ErrorCount, result.WarningCount)

			for _, issue := range result.Issues {
				icon := "⚠️"
				if issue.Level == "error" {
					icon = "❌"
				}
				fmt.Printf("%s %s:%d: %s\n", icon, issue.File, issue.Line, issue.Message)
			}

			if fix && result.FixCount > 0 {
				fmt.Printf("🔧 Auto-fixed %d issues\n", result.FixCount)
			}
		} else {
			fmt.Printf("✅ All ADRs are valid!\n")
			if verbose {
				fmt.Printf("📊 Validation stats:\n")
				fmt.Printf("   • %d ADRs checked\n", result.FileCount)
				fmt.Printf("   • %d diagrams validated\n", result.DiagramCount)
				if result.WarningCount > 0 {
					fmt.Printf("   • %d warnings (non-blocking)\n", result.WarningCount)
				}
			}
		}

		// Suggestions for improvement
		if !strict && !result.HasErrors() {
			fmt.Printf("💡 Tip: Run with --strict for additional style checks\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().BoolVar(&strict, "strict", false, "enable strict validation with additional style checks")
	validateCmd.Flags().BoolVar(&fix, "fix", false, "automatically fix common issues")
}
