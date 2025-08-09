package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/euforicio/adr-demo/internal/config"
	"github.com/euforicio/adr-demo/internal/generator"
)

// Config holds the server configuration
type Config struct {
	Host    string
	Port    int
	Verbose bool
}

// Server represents the development server
type Server struct {
	config    *Config
	generator *generator.Generator
}

// New creates a new development server
func New(config *Config) *Server {
	return &Server{
		config: config,
	}
}

// Start starts the development server
func (s *Server) Start() error {
	// Create generator for initial build and dynamic serving
	genConfig := config.DefaultConfig()
	genConfig.Verbose = s.config.Verbose

	s.generator = generator.New(genConfig)

	// Initial load and cache of ADRs
	if s.config.Verbose {
		fmt.Println("🔨 Loading ADRs and building cache...")
	}

	if err := s.generator.LoadADRsOnly(); err != nil {
		return fmt.Errorf("initial ADR load failed: %w", err)
	}

	if s.config.Verbose {
		fmt.Println("🚀 Starting dynamic server with cached content...")
	}

	// Setup HTTP server
	s.setupRoutes()

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	if s.config.Verbose {
		fmt.Printf("🌐 Server listening on http://%s\n", addr)
	}

	return http.ListenAndServe(addr, nil)
}

// setupRoutes configures HTTP routes
func (s *Server) setupRoutes() {
	// Dynamic routes for ADR rendering
	http.HandleFunc("/", s.handleRequest)
	http.HandleFunc("/search", s.handleSearch)
	http.HandleFunc("/search.html", s.handleSearch)
	http.HandleFunc("/search-index.json", s.handleSearchIndex)
	http.HandleFunc("/docs", s.handleDocs)

	// Serve static assets if they exist
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}

// handleRequest routes requests based on URL pattern
func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	// Route based on URL pattern
	if r.URL.Path == "/" || r.URL.Path == "/index.html" {
		s.handleIndex(w, r)
		return
	}

	// Check if it's an ADR page
	if strings.HasPrefix(r.URL.Path, "/adr-") {
		s.handleADR(w, r)
		return
	}

	// Not found
	http.NotFound(w, r)
}

// handleIndex serves the main index page with all ADRs
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	// Render index page (will use cache if unchanged)
	if err := s.generator.RenderIndexPage(w); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to render index: %v", err),
			http.StatusInternalServerError,
		)
		return
	}
}

// handleADR serves individual ADR pages
func (s *Server) handleADR(w http.ResponseWriter, r *http.Request) {
	// Extract ADR number from URL path
	path := r.URL.Path

	// Remove /adr- prefix and .html suffix if present
	adrPart := strings.TrimPrefix(path, "/adr-")
	adrPart = strings.TrimSuffix(adrPart, ".html")

	// Render specific ADR page (will use cache if unchanged)
	if err := s.generator.RenderADRPage(w, adrPart); err != nil {
		http.Error(w, fmt.Sprintf("Failed to render ADR: %v", err), http.StatusNotFound)
		return
	}
}

// handleSearch serves the search page
func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	// Render search page (will use cache if unchanged)
	if err := s.generator.RenderSearchPage(w); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to render search: %v", err),
			http.StatusInternalServerError,
		)
		return
	}
}

// handleDocs serves the documentation page with README.md content
func (s *Server) handleDocs(w http.ResponseWriter, r *http.Request) {
	// Read README.md file
	content, err := os.ReadFile("README.md")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read documentation: %v", err), http.StatusInternalServerError)
		return
	}

	// Render documentation page (will use generator's markdown processing)
	if err := s.generator.RenderDocsPage(w, string(content)); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to render documentation: %v", err),
			http.StatusInternalServerError,
		)
		return
	}
}

// handleSearchIndex serves the search index JSON
func (s *Server) handleSearchIndex(w http.ResponseWriter, r *http.Request) {
	// Generate search index dynamically
	searchIndex := s.generateSearchIndex()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(searchIndex); err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to encode search index: %v", err),
			http.StatusInternalServerError,
		)
		return
	}
}

// generateSearchIndex creates search index from current ADRs
func (s *Server) generateSearchIndex() map[string]interface{} {
	type SearchItem struct {
		Number      string `json:"number"`
		Title       string `json:"title"`
		Status      string `json:"status"`
		Content     string `json:"content"`
		DiagramType string `json:"diagramType"`
		URL         string `json:"url"`
	}

	var searchItems []SearchItem
	adrs := s.generator.GetADRs()

	for _, adr := range adrs {
		// Clean content for search (remove markdown syntax)
		cleanContent := strings.ReplaceAll(adr.Content, "#", "")
		cleanContent = strings.ReplaceAll(cleanContent, "*", "")
		cleanContent = strings.ReplaceAll(cleanContent, "_", "")

		// Limit content length for search index
		if len(cleanContent) > 500 {
			cleanContent = cleanContent[:500] + "..."
		}

		item := SearchItem{
			Number:      adr.Number,
			Title:       adr.Title,
			Status:      adr.Status,
			Content:     cleanContent,
			DiagramType: adr.DiagramType,
			URL:         fmt.Sprintf("adr-%s.html", adr.Number),
		}

		searchItems = append(searchItems, item)
	}

	return map[string]interface{}{
		"generated": len(searchItems),
		"items":     searchItems,
	}
}
