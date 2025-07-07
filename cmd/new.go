package cmd

import (
	"fmt"
	"log"

	"github.com/euforicio/adr-demo/internal/generator"
	"github.com/spf13/cobra"
)

var (
	status   string
	category string
	force    bool
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [title]",
	Short: "Create a new ADR from template",
	Long: `Create a new Architecture Decision Record from the template.

The new ADR will:
• Use the next available number (e.g., 0011-title.md)
• Be pre-filled with the standard ADR template
• Have the title automatically formatted in kebab-case
• Be created in the adr/ directory

Examples:
  adr-gen new "Use Redis for Caching"
  adr-gen new "Implement API Gateway" --status accepted
  adr-gen new "Frontend Framework Choice" --category frontend
  adr-gen new "Security Policy" --category security --status accepted`,
	Args: cobra.RangeArgs(1, 10), // Allow multiple words for title
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Please provide a title for the new ADR")
		}

		// Join all arguments to form the title
		title := ""
		for i, arg := range args {
			if i > 0 {
				title += " "
			}
			title += arg
		}

		if verbose {
			fmt.Printf("📄 Creating new ADR...\n")
			fmt.Printf("   Title: %s\n", title)
			fmt.Printf("   Status: %s\n", status)
			if category != "" {
				fmt.Printf("   Category: %s\n", category)
			}
		}

		creator := generator.NewADRCreator(&generator.ADRConfig{
			Title:    title,
			Status:   status,
			Category: category,
			Force:    force,
			Verbose:  verbose,
		})

		filename, err := creator.Create()
		if err != nil {
			log.Fatalf("Failed to create ADR: %v", err)
		}

		fmt.Printf("✅ Created new ADR: %s\n", filename)
		fmt.Printf("💡 Edit the file to add your decision details\n")

		if verbose {
			fmt.Printf("📝 Next steps:\n")
			fmt.Printf("   1. Edit %s\n", filename)
			fmt.Printf("   2. Add context and decision details\n")
			fmt.Printf("   3. Consider adding a Mermaid diagram\n")
			fmt.Printf("   4. Run 'adr-gen serve' to preview\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringVarP(&status, "status", "s", "Proposed", "initial status for the ADR (Proposed, Accepted, Deprecated, Superseded)")
	newCmd.Flags().StringVarP(&category, "category", "c", "", "category/folder for the ADR (e.g., frontend, infrastructure, security)")
	newCmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite existing ADR if it exists")
}
