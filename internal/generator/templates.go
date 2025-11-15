package generator

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/euforicio/adr-demo/internal/config"
)

// loadTemplates loads and parses all HTML templates
func (g *Generator) loadTemplates() error {
	templateDir := "templates"

	// Check if templates directory exists
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return fmt.Errorf("templates directory not found: %s", templateDir)
	}

	// We'll store template functions for reuse
	g.funcMap = template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"statusClass": func(status string) string {
			return g.config.GetStatusClass(status)
		},
		"statusEmoji": func(status string) string {
			return g.config.GetStatusIcon(status)
		},
		"statusIcon": func(status string) string {
			return g.config.GetStatusIcon(status)
		},
		"groupByCategory": func(adrs []*ADR) map[string][]*ADR {
			groups := make(map[string][]*ADR)
			for _, adr := range adrs {
				category := adr.Category
				if category == "" {
					category = g.config.DefaultCategory
				}
				groups[category] = append(groups[category], adr)
			}
			return groups
		},
		"contains": func(s, substr string) bool {
			return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
		},
		"config": func() *config.Config {
			return g.config
		},
	}

	// We don't parse all templates together anymore to avoid block conflicts
	// Instead, we'll parse them individually in renderPage
	return nil
}

// generateIndexPage creates the main index page
func (g *Generator) generateIndexPage() error {
	data := struct {
		Title          string
		ADRs           []*ADR
		Stats          map[string]int
		BaseURL        string
		BreadcrumbType string
	}{
		Title:          "Architecture Decision Records",
		ADRs:           g.adrs,
		BaseURL:        g.config.BaseURL,
		BreadcrumbType: "index",
		Stats: map[string]int{
			"Total":      len(g.adrs),
			"Accepted":   g.countByStatus("Accepted"),
			"Proposed":   g.countByStatus("Proposed"),
			"Deprecated": g.countByStatus("Deprecated"),
			"Superseded": g.countByStatus("Superseded"),
			"Diagrams":   g.stats.DiagramCount,
		},
	}

	return g.renderPage("index.html", "index.html", data)
}

// generateADRPages creates individual pages for each ADR
func (g *Generator) generateADRPages() error {
	for i, adr := range g.adrs {
		data := struct {
			Title          string
			ADR            *ADR
			ADRs           []*ADR
			Previous       *ADR
			Next           *ADR
			BaseURL        string
			BreadcrumbType string
		}{
			Title:          fmt.Sprintf("ADR-%s: %s", adr.Number, adr.Title),
			ADR:            adr,
			ADRs:           g.adrs,
			BaseURL:        g.config.BaseURL,
			BreadcrumbType: "adr",
		}

		// Set previous/next navigation
		if i > 0 {
			data.Previous = g.adrs[i-1]
		}
		if i < len(g.adrs)-1 {
			data.Next = g.adrs[i+1]
		}

		filename := fmt.Sprintf("adr-%s.html", adr.Number)
		if err := g.renderPage("adr.html", filename, data); err != nil {
			return err
		}
		g.stats.PageCount++
	}

	return nil
}

// generateSearchPage creates the search page
func (g *Generator) generateSearchPage() error {
	data := struct {
		Title          string
		ADRs           []*ADR
		BaseURL        string
		BreadcrumbType string
	}{
		Title:          "Search ADRs",
		ADRs:           g.adrs,
		BaseURL:        g.config.BaseURL,
		BreadcrumbType: "search",
	}

	if err := g.renderPage("search.html", "search.html", data); err != nil {
		return err
	}
	g.stats.PageCount++
	return nil
}

// renderPage renders a template to a file
func (g *Generator) renderPage(templateName, filename string, data interface{}) error {
	outputPath := filepath.Join(g.config.OutputDirectory, filename)

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer file.Close()

	// Parse base template and specific template individually to avoid block conflicts
	templateDir := "templates"
	baseTemplatePath := filepath.Join(templateDir, "base.html")
	specificTemplatePath := filepath.Join(templateDir, templateName)

	// Create a fresh template for this specific page
	tmpl, err := template.New("base.html").
		Funcs(g.funcMap).
		ParseFiles(baseTemplatePath, specificTemplatePath)
	if err != nil {
		return fmt.Errorf("failed to parse templates for %s: %w", templateName, err)
	}

	if err := tmpl.ExecuteTemplate(file, templateName, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return nil
}

// countByStatus counts ADRs by status
func (g *Generator) countByStatus(status string) int {
	count := 0
	for _, adr := range g.adrs {
		if adr.Status == status {
			count++
		}
	}
	return count
}
