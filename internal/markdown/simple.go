package markdown

import (
	"bytes"
	"fmt"
	htmlpkg "html"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
)

// Config holds the markdown processor configuration
type Config struct {
	EnableGFM     bool
	EnableMermaid bool
	Verbose       bool
	BaseURL       string
}

// SimpleProcessor is a simplified markdown processor
type SimpleProcessor struct {
	config   *Config
	markdown goldmark.Markdown
}

// NewSimple creates a simplified markdown processor
func NewSimple(config *Config) *SimpleProcessor {
	// Configure goldmark with GitHub Flavored Markdown
	extensions := []goldmark.Extender{
		extension.GFM,            // GitHub Flavored Markdown
		extension.Footnote,       // Footnotes
		extension.DefinitionList, // Definition lists
	}

	// Configure HTML renderer
	htmlOptions := []renderer.Option{
		goldmarkhtml.WithHardWraps(),
		goldmarkhtml.WithXHTML(),
		goldmarkhtml.WithUnsafe(), // Allow raw HTML (needed for Mermaid)
	}

	// Configure parser
	parserOptions := []parser.Option{
		parser.WithAutoHeadingID(), // Automatic heading IDs
	}

	md := goldmark.New(
		goldmark.WithExtensions(extensions...),
		goldmark.WithParserOptions(parserOptions...),
		goldmark.WithRendererOptions(htmlOptions...),
	)

	return &SimpleProcessor{
		config:   config,
		markdown: md,
	}
}

// Process converts markdown content to HTML
func (p *SimpleProcessor) Process(content string) (string, error) {
	var buf bytes.Buffer

	// Pre-process content
	content = p.preprocessContent(content)

	// Convert markdown to HTML
	if err := p.markdown.Convert([]byte(content), &buf); err != nil {
		return "", fmt.Errorf("failed to convert markdown: %w", err)
	}

	html := buf.String()

	// Post-process HTML
	html = p.postprocessHTML(html)

	return html, nil
}

// preprocessContent handles content before markdown processing
func (p *SimpleProcessor) preprocessContent(content string) string {
	// Add line breaks for better rendering
	content = strings.ReplaceAll(content, "\r\n", "\n")

	// Ensure proper spacing around headings
	re := regexp.MustCompile(`\n(#{1,6}\s+.+)\n`)
	content = re.ReplaceAllString(content, "\n\n$1\n\n")

	return content
}

// postprocessHTML handles HTML after markdown processing
func (p *SimpleProcessor) postprocessHTML(html string) string {
	// Process ADR links
	html = p.processADRLinks(html)

	// Process Mermaid diagrams
	if p.config.EnableMermaid {
		html = p.processMermaidDiagrams(html)
	}

	return html
}


// processMermaidDiagrams handles Mermaid diagram blocks
func (p *SimpleProcessor) processMermaidDiagrams(html string) string {
	// Find Mermaid code blocks ((?s) allows . to match newlines)
	re := regexp.MustCompile(`(?s)<pre><code class="language-mermaid">(.*?)</code></pre>`)

	diagramCount := 0
	html = re.ReplaceAllStringFunc(html, func(block string) string {
		matches := re.FindStringSubmatch(block)
		if len(matches) != 2 {
			return block
		}

		diagramCount++
		code := matches[1]

		// Create Mermaid diagram container
		return fmt.Sprintf(`
<div class="mermaid-container" id="mermaid-%d">
	<div class="mermaid-toolbar">
		<button class="mermaid-fullscreen" onclick="openMermaidFullscreen('mermaid-%d')" title="View fullscreen">
			⛶
		</button>
		<button class="mermaid-copy" onclick="copyMermaidCode('mermaid-%d')" title="Copy diagram code">
			📋
		</button>
	</div>
	<div class="mermaid-diagram" data-diagram="%s" onclick="openMermaidFullscreen('mermaid-%d')" style="cursor: pointer;" title="Click to view fullscreen">
		<div class="mermaid">%s</div>
	</div>
</div>`, diagramCount, diagramCount, diagramCount, htmlpkg.EscapeString(code), diagramCount, code)
	})

	return html
}

// processADRLinks converts ADR markdown links to HTML URLs
func (p *SimpleProcessor) processADRLinks(html string) string {
	// Match links to ADR markdown files
	// Pattern: <a href="NNNN-some-title.md">Link Text</a>
	re := regexp.MustCompile(`<a href="([0-9]{4}-[a-z0-9-]+\.md)"([^>]*)>([^<]*)</a>`)
	
	html = re.ReplaceAllStringFunc(html, func(match string) string {
		matches := re.FindStringSubmatch(match)
		if len(matches) != 4 {
			return match
		}
		
		markdownFile := matches[1]   // e.g., "0001-record-architecture-decisions.md"
		attributes := matches[2]     // any additional attributes
		linkText := matches[3]       // the link text
		
		// Extract ADR number from filename
		adrNumber := extractADRNumberFromFilename(markdownFile)
		if adrNumber == "" {
			return match // If we can't extract number, leave unchanged
		}
		
		// Convert to HTML URL with base URL prefix
		htmlURL := fmt.Sprintf("%s/adr-%s.html", strings.TrimSuffix(p.config.BaseURL, "/"), adrNumber)
		if p.config.BaseURL == "" {
			htmlURL = fmt.Sprintf("/adr-%s.html", adrNumber)
		}
		
		return fmt.Sprintf(`<a href="%s"%s>%s</a>`, htmlURL, attributes, linkText)
	})
	
	return html
}

// extractADRNumberFromFilename extracts the 4-digit ADR number from a filename
func extractADRNumberFromFilename(filename string) string {
	if len(filename) < 4 {
		return ""
	}
	
	// Check if the first 4 characters are digits
	for i := 0; i < 4; i++ {
		if filename[i] < '0' || filename[i] > '9' {
			return ""
		}
	}
	
	return filename[:4]
}
