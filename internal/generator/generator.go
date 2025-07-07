package generator

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/euforicio/adr-demo/internal/config"
	"github.com/euforicio/adr-demo/internal/markdown"
)

// CacheEntry represents a cached rendered page
type CacheEntry struct {
	Content   []byte
	Timestamp time.Time
	FileHash  string
}

// Generator handles the static site generation
type Generator struct {
	config      *config.Config
	templates   *template.Template
	funcMap     template.FuncMap
	adrs        []*ADR
	stats       Stats
	renderCache map[string]*CacheEntry // Cache for rendered pages
	cacheMutex  sync.RWMutex           // Mutex for cache access
}

// Stats holds build statistics
type Stats struct {
	ADRCount     int
	PageCount    int
	AssetCount   int
	DiagramCount int
}

// ADR represents a parsed Architecture Decision Record
type ADR struct {
	Number      string
	Title       string
	Status      string
	Content     string
	HTMLContent template.HTML
	FilePath    string
	FileName    string
	DiagramType string
	Category    string // Category based on folder structure
	CreatedAt   time.Time
	ModifiedAt  time.Time
	FileHash    string // SHA256 hash of the source file content
}

// New creates a new generator instance
func New(cfg *config.Config) *Generator {
	return &Generator{
		config:      cfg,
		adrs:        make([]*ADR, 0),
		renderCache: make(map[string]*CacheEntry),
	}
}

// Build generates the complete static site
func (g *Generator) Build() error {
	if g.config.Verbose {
		fmt.Println("📚 Loading ADR files...")
	}

	// Load and parse all ADR files
	if err := g.loadADRs(); err != nil {
		return fmt.Errorf("failed to load ADRs: %w", err)
	}

	if g.config.Verbose {
		fmt.Println("🎨 Loading templates...")
	}

	// Load templates
	if err := g.loadTemplates(); err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	if g.config.Verbose {
		fmt.Println("🏗️  Generating pages...")
	}

	// Create output directory
	if err := os.MkdirAll(g.config.OutputDirectory, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate pages
	if err := g.generateIndexPage(); err != nil {
		return fmt.Errorf("failed to generate index page: %w", err)
	}

	if err := g.generateADRPages(); err != nil {
		return fmt.Errorf("failed to generate ADR pages: %w", err)
	}

	if err := g.generateSearchPage(); err != nil {
		return fmt.Errorf("failed to generate search page: %w", err)
	}

	if g.config.Verbose {
		fmt.Println("📦 Copying static assets...")
	}

	// Copy static assets
	if err := g.copyAssets(); err != nil {
		return fmt.Errorf("failed to copy assets: %w", err)
	}

	// Generate search index
	if err := g.generateSearchIndex(); err != nil {
		return fmt.Errorf("failed to generate search index: %w", err)
	}

	return nil
}

// loadADRs finds and parses all ADR markdown files from the flat structure
func (g *Generator) loadADRs() error {
	adrDir := g.config.ADRDirectory

	processor := markdown.NewSimple(&markdown.Config{
		EnableGFM:     true,
		EnableMermaid: true,
		Verbose:       g.config.Verbose,
		BaseURL:       g.config.BaseURL,
	})

	// Read all files in the ADR directory (flat structure)
	files, err := os.ReadDir(adrDir)
	if err != nil {
		return fmt.Errorf("failed to read ADR directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Only process .md files
		if !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		// Skip template file
		if file.Name() == "template.md" {
			continue
		}

		// Parse ADR number from filename
		if !isValidADRFilename(file.Name()) {
			if g.config.Verbose {
				fmt.Printf("⚠️  Skipping invalid filename: %s\n", file.Name())
			}
			continue
		}

		path := filepath.Join(adrDir, file.Name())
		adr, err := g.parseADR(path, processor)
		if err != nil {
			return fmt.Errorf("failed to parse ADR %s: %w", path, err)
		}

		// Extract category from ADR content
		adr.Category = g.extractCategoryFromContent(adr.Content)

		g.adrs = append(g.adrs, adr)
	}

	// Sort ADRs by number
	sort.Slice(g.adrs, func(i, j int) bool {
		return g.adrs[i].Number < g.adrs[j].Number
	})

	g.stats.ADRCount = len(g.adrs)

	if g.config.Verbose {
		fmt.Printf("📄 Loaded %d ADRs\n", len(g.adrs))
	}

	return nil
}

// parseADR parses a single ADR file
func (g *Generator) parseADR(filePath string, processor *markdown.SimpleProcessor) (*ADR, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Calculate file hash
	hash := sha256.Sum256(content)
	fileHash := hex.EncodeToString(hash[:])

	fileName := filepath.Base(filePath)
	number := extractADRNumber(fileName)
	title := extractTitleFromContent(string(content))
	status := extractStatusFromContent(string(content))

	// Process markdown to HTML
	htmlContent, err := processor.Process(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to process markdown: %w", err)
	}

	// Count diagrams
	diagramCount := strings.Count(string(content), "```mermaid")
	g.stats.DiagramCount += diagramCount

	// Determine diagram type
	diagramType := detectDiagramType(string(content))

	// Get file info
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	return &ADR{
		Number:      number,
		Title:       title,
		Status:      status,
		Content:     string(content),
		HTMLContent: template.HTML(htmlContent),
		FilePath:    filePath,
		FileName:    fileName,
		DiagramType: diagramType,
		CreatedAt:   info.ModTime(), // Approximation
		ModifiedAt:  info.ModTime(),
		FileHash:    fileHash,
	}, nil
}

// GetStats returns build statistics
func (g *Generator) GetStats() Stats {
	return g.stats
}

// GetADRs returns the current ADRs
func (g *Generator) GetADRs() []*ADR {
	return g.adrs
}

// Helper functions

func isValidADRFilename(filename string) bool {
	// Check if filename matches NNNN-*.md pattern
	if len(filename) < 8 {
		return false
	}

	// Check for 4 digits at the start
	for i := 0; i < 4; i++ {
		if filename[i] < '0' || filename[i] > '9' {
			return false
		}
	}

	// Check for dash after digits
	if filename[4] != '-' {
		return false
	}

	// Check for .md extension
	return strings.HasSuffix(filename, ".md")
}

func extractADRNumber(filename string) string {
	if len(filename) >= 4 {
		return filename[:4]
	}
	return "0000"
}

func extractTitleFromContent(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(line[2:])
		}
	}
	return "Untitled ADR"
}

func extractStatusFromContent(content string) string {
	lines := strings.Split(content, "\n")
	statusSection := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "## Status") {
			statusSection = true
			continue
		}

		if statusSection {
			if strings.HasPrefix(trimmed, "##") {
				break // Next section
			}
			if trimmed != "" {
				return trimmed
			}
		}
	}

	return "Unknown"
}

func (g *Generator) extractCategoryFromContent(content string) string {
	lines := strings.Split(content, "\n")

	// Look for Category field in the content
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for Category field (case insensitive)
		if strings.HasPrefix(strings.ToLower(trimmed), "category:") {
			category := strings.TrimSpace(trimmed[9:]) // Remove "Category:" prefix
			if category != "" && g.config.IsValidCategory(category) {
				return category
			}
		}

		// Check for Category in markdown heading format
		if strings.HasPrefix(trimmed, "## Category") {
			// Look for the next non-empty line
			for j := i + 1; j < len(lines) && j < i+11; j++ {
				nextTrimmed := strings.TrimSpace(lines[j])
				if strings.HasPrefix(nextTrimmed, "##") {
					break // Next section
				}
				if nextTrimmed != "" && g.config.IsValidCategory(nextTrimmed) {
					return nextTrimmed
				}
			}
		}
	}

	// Default to configured default category
	return g.config.DefaultCategory
}

func detectDiagramType(content string) string {
	if strings.Contains(content, "C4Context") {
		return "Context"
	}
	if strings.Contains(content, "C4Container") {
		return "Container"
	}
	if strings.Contains(content, "C4Component") {
		return "Component"
	}
	if strings.Contains(content, "sequenceDiagram") {
		return "Sequence"
	}
	if strings.Contains(content, "stateDiagram") {
		return "State"
	}
	if strings.Contains(content, "flowchart") {
		return "Flowchart"
	}
	if strings.Contains(content, "```mermaid") {
		return "Diagram"
	}
	return "-"
}

// LoadADRsOnly loads ADRs without generating any output files
func (g *Generator) LoadADRsOnly() error {
	if g.config.Verbose {
		fmt.Println("📚 Loading ADR files...")
	}

	// Reset stats and ADRs for fresh load
	g.adrs = make([]*ADR, 0)
	g.stats = Stats{}

	// Load and parse all ADR files
	if err := g.loadADRs(); err != nil {
		return fmt.Errorf("failed to load ADRs: %w", err)
	}

	if g.config.Verbose {
		fmt.Println("🎨 Loading templates...")
	}

	// Load templates
	if err := g.loadTemplates(); err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	return nil
}

// RenderIndexPage renders the index page to a writer
func (g *Generator) RenderIndexPage(w io.Writer) error {
	// Create cache key based on all ADR hashes
	cacheKey := g.generateIndexCacheKey()

	// Check cache first
	if cached := g.getCachedContent(cacheKey); cached != nil {
		_, err := w.Write(cached)
		return err
	}

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

	return g.renderPageToWriterWithCache("index.html", w, data, cacheKey)
}

// RenderADRPage renders a specific ADR page to a writer
func (g *Generator) RenderADRPage(w io.Writer, adrNumber string) error {
	// Find the ADR by number
	var targetADR *ADR
	var adrIndex int
	for i, adr := range g.adrs {
		if adr.Number == adrNumber {
			targetADR = adr
			adrIndex = i
			break
		}
	}

	if targetADR == nil {
		return fmt.Errorf("ADR %s not found", adrNumber)
	}

	// Create cache key based on the specific ADR hash
	cacheKey := fmt.Sprintf("adr-%s-%s", adrNumber, targetADR.FileHash)

	// Check cache first
	if cached := g.getCachedContent(cacheKey); cached != nil {
		_, err := w.Write(cached)
		return err
	}

	data := struct {
		Title          string
		ADR            *ADR
		ADRs           []*ADR
		Previous       *ADR
		Next           *ADR
		BaseURL        string
		BreadcrumbType string
	}{
		Title:          fmt.Sprintf("ADR-%s: %s", targetADR.Number, targetADR.Title),
		ADR:            targetADR,
		ADRs:           g.adrs,
		BaseURL:        g.config.BaseURL,
		BreadcrumbType: "adr",
	}

	// Set previous/next navigation
	if adrIndex > 0 {
		data.Previous = g.adrs[adrIndex-1]
	}
	if adrIndex < len(g.adrs)-1 {
		data.Next = g.adrs[adrIndex+1]
	}

	return g.renderPageToWriterWithCache("adr.html", w, data, cacheKey)
}

// RenderSearchPage renders the search page to a writer
func (g *Generator) RenderSearchPage(w io.Writer) error {
	// Create cache key based on all ADR hashes (search page shows all ADRs)
	cacheKey := g.generateSearchCacheKey()

	// Check cache first
	if cached := g.getCachedContent(cacheKey); cached != nil {
		_, err := w.Write(cached)
		return err
	}

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

	return g.renderPageToWriterWithCache("search.html", w, data, cacheKey)
}

// RenderDocsPage renders the documentation page to a writer
func (g *Generator) RenderDocsPage(w io.Writer, readmeContent string) error {
	// Create cache key based on README content hash
	hash := sha256.Sum256([]byte(readmeContent))
	fileHash := hex.EncodeToString(hash[:])
	cacheKey := fmt.Sprintf("docs-%s", fileHash)

	// Check cache first
	if cached := g.getCachedContent(cacheKey); cached != nil {
		_, err := w.Write(cached)
		return err
	}

	// Process markdown to HTML
	processor := markdown.NewSimple(&markdown.Config{
		EnableGFM:     true,
		EnableMermaid: true,
		Verbose:       g.config.Verbose,
		BaseURL:       g.config.BaseURL,
	})

	htmlContent, err := processor.Process(readmeContent)
	if err != nil {
		return fmt.Errorf("failed to process markdown: %w", err)
	}

	data := struct {
		Title          string
		Content        template.HTML
		ADRs           []*ADR
		BaseURL        string
		BreadcrumbType string
	}{
		Title:          "Documentation",
		Content:        template.HTML(htmlContent),
		ADRs:           g.adrs,
		BaseURL:        g.config.BaseURL,
		BreadcrumbType: "docs",
	}

	return g.renderPageToWriterWithCache("docs.html", w, data, cacheKey)
}

// renderPageToWriter renders a template to a writer (for dynamic serving)
func (g *Generator) renderPageToWriter(templateName string, w io.Writer, data interface{}) error {
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

	if err := tmpl.ExecuteTemplate(w, templateName, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return nil
}

// getCachedContent retrieves content from cache if it exists and is valid
func (g *Generator) getCachedContent(cacheKey string) []byte {
	g.cacheMutex.RLock()
	defer g.cacheMutex.RUnlock()

	entry, exists := g.renderCache[cacheKey]
	if !exists {
		return nil
	}

	// Cache is valid, return content
	return entry.Content
}

// setCachedContent stores content in cache
func (g *Generator) setCachedContent(cacheKey string, content []byte, fileHash string) {
	g.cacheMutex.Lock()
	defer g.cacheMutex.Unlock()

	g.renderCache[cacheKey] = &CacheEntry{
		Content:   content,
		Timestamp: time.Now(),
		FileHash:  fileHash,
	}
}

// generateIndexCacheKey creates a cache key for the index page based on all ADR hashes
func (g *Generator) generateIndexCacheKey() string {
	hash := sha256.New()
	for _, adr := range g.adrs {
		hash.Write([]byte(adr.FileHash))
	}
	return fmt.Sprintf("index-%x", hash.Sum(nil))
}

// generateSearchCacheKey creates a cache key for the search page based on all ADR hashes
func (g *Generator) generateSearchCacheKey() string {
	hash := sha256.New()
	for _, adr := range g.adrs {
		hash.Write([]byte(adr.FileHash))
	}
	return fmt.Sprintf("search-%x", hash.Sum(nil))
}

// renderPageToWriterWithCache renders a template to a writer with caching
func (g *Generator) renderPageToWriterWithCache(templateName string, w io.Writer, data interface{}, cacheKey string) error {
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

	// Render to a buffer first so we can cache it
	var buf strings.Builder
	if err := tmpl.ExecuteTemplate(&buf, templateName, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	// Cache the rendered content
	content := []byte(buf.String())
	g.setCachedContent(cacheKey, content, "")

	// Write to the actual writer
	_, err = w.Write(content)
	return err
}
