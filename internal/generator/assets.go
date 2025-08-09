package generator

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// copyAssets copies static assets to the output directory
func (g *Generator) copyAssets() error {
	staticDir := "static"
	outputStaticDir := filepath.Join(g.config.OutputDirectory, "static")

	// Check if static directory exists
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		// Create minimal CSS if no static directory exists
		return g.createMinimalAssets()
	}

	// Create output static directory
	if err := os.MkdirAll(outputStaticDir, 0755); err != nil {
		return fmt.Errorf("failed to create static output directory: %w", err)
	}

	// Copy all static files
	err := filepath.Walk(staticDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path
		relPath, err := filepath.Rel(staticDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(outputStaticDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		// Copy file
		return g.copyFile(path, destPath)
	})

	if err != nil {
		return fmt.Errorf("failed to copy static assets: %w", err)
	}

	return nil
}

// copyFile copies a single file
func (g *Generator) copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create destination directory if it doesn't exist
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	g.stats.AssetCount++
	return nil
}

// createMinimalAssets creates comprehensive CSS/JS if no static directory exists
func (g *Generator) createMinimalAssets() error {
	staticDir := filepath.Join(g.config.OutputDirectory, "static")
	cssDir := filepath.Join(staticDir, "css")
	jsDir := filepath.Join(staticDir, "js")

	// Create directories
	if err := os.MkdirAll(cssDir, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(jsDir, 0755); err != nil {
		return err
	}

	// Create comprehensive CSS with layout
	css := `/* ADR Generator - Complete Styles */
:root {
    --primary-color: #5e6ad2;
    --primary-hover: #4f5bd1;
    --bg-primary: #ffffff;
    --bg-secondary: #f9fafb;
    --bg-tertiary: #f3f4f6;
    --text-primary: #111827;
    --text-secondary: #6b7280;
    --text-muted: #9ca3af;
    --border-color: #e5e7eb;
    --border-light: #f3f4f6;
    --shadow-sm: 0 1px 2px 0 rgb(0 0 0 / 0.05);
    --shadow-md: 0 4px 6px -1px rgb(0 0 0 / 0.1);
    --shadow-lg: 0 10px 15px -3px rgb(0 0 0 / 0.1);
    --radius-sm: 0.375rem;
    --radius-md: 0.5rem;
    --radius-lg: 0.75rem;
}

@media (prefers-color-scheme: dark) {
    :root {
        --bg-primary: #0f0f0f;
        --bg-secondary: #1a1a1a;
        --bg-tertiary: #262626;
        --text-primary: #f9fafb;
        --text-secondary: #d1d5db;
        --text-muted: #9ca3af;
        --border-color: #374151;
        --border-light: #4b5563;
        --shadow-sm: 0 1px 2px 0 rgb(0 0 0 / 0.3);
        --shadow-md: 0 4px 6px -1px rgb(0 0 0 / 0.3);
        --shadow-lg: 0 10px 15px -3px rgb(0 0 0 / 0.3);
    }
}

@view-transition {
    navigation: auto;
}

* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
    background: var(--bg-primary);
    color: var(--text-primary);
    line-height: 1.6;
    overflow-x: hidden;
}

/* ===== LAYOUT ===== */
.app-container {
    display: flex;
    min-height: 100vh;
}

.sidebar {
    width: 320px;
    background: var(--bg-secondary);
    border-right: 1px solid var(--border-color);
    display: flex;
    flex-direction: column;
    position: fixed;
    height: 100vh;
    overflow-y: auto;
    view-transition-name: sidebar;
    z-index: 1000;
}

.main-content {
    flex: 1;
    margin-left: 320px;
    display: flex;
    flex-direction: column;
    min-height: 100vh;
}

.main-header {
    background: var(--bg-primary);
    border-bottom: 1px solid var(--border-color);
    padding: 1rem 2rem;
    position: sticky;
    top: 0;
    z-index: 100;
    view-transition-name: main-header;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.adr-content {
    flex: 1;
    padding: 2rem;
    max-width: none;
    view-transition-name: adr-content;
}

/* ===== SIDEBAR ===== */
.sidebar-header {
    padding: 1.5rem 1rem;
    border-bottom: 1px solid var(--border-color);
    background: var(--bg-primary);
}

.site-title {
    margin: 0;
    font-size: 1.25rem;
    font-weight: 700;
}

.home-link {
    color: var(--text-primary);
    text-decoration: none;
    display: flex;
    align-items: center;
    gap: 0.5rem;
}

.home-link:hover {
    color: var(--primary-color);
}

.icon {
    font-size: 1.5rem;
}

.search-container {
    padding: 1rem;
    border-bottom: 1px solid var(--border-color);
    position: relative;
}

.search-input {
    width: 100%;
    padding: 0.75rem 1rem;
    border: 1px solid var(--border-color);
    border-radius: var(--radius-md);
    background: var(--bg-primary);
    color: var(--text-primary);
    font-size: 0.875rem;
}

.search-input:focus {
    outline: none;
    border-color: var(--primary-color);
    box-shadow: 0 0 0 3px rgb(94 106 210 / 0.1);
}

.search-clear {
    position: absolute;
    right: 1.5rem;
    top: 50%;
    transform: translateY(-50%);
    background: none;
    border: none;
    color: var(--text-muted);
    cursor: pointer;
    font-size: 1.25rem;
    padding: 0.25rem;
}

.filter-container {
    padding: 0 1rem 1rem;
    border-bottom: 1px solid var(--border-color);
}

.status-filter {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid var(--border-color);
    border-radius: var(--radius-sm);
    background: var(--bg-primary);
    color: var(--text-primary);
    font-size: 0.875rem;
}

.adr-nav {
    flex: 1;
    overflow-y: auto;
    padding: 1rem 0;
}

.adr-category {
    margin-bottom: 1.5rem;
}

.adr-category-title {
    padding: 0.5rem 1rem;
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--text-secondary);
    text-transform: uppercase;
    letter-spacing: 0.05em;
    margin-bottom: 0.5rem;
}

.adr-category-list {
    list-style: none;
}

.adr-item {
    margin-bottom: 0.25rem;
}

.adr-link {
    display: block;
    padding: 0.75rem 1rem;
    margin: 0 0.5rem;
    border-radius: var(--radius-md);
    color: var(--text-primary);
    text-decoration: none;
    transition: all 0.2s ease;
}

.adr-link:hover {
    background: var(--bg-tertiary);
    text-decoration: none;
}

.adr-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 0.5rem;
}

.adr-number {
    font-family: 'JetBrains Mono', 'Fira Code', monospace;
    font-size: 0.75rem;
    font-weight: 600;
    color: var(--text-muted);
}

.adr-status {
    font-size: 0.75rem;
    padding: 0.125rem 0.375rem;
    border-radius: var(--radius-sm);
    font-weight: 500;
}

.adr-title {
    font-size: 0.875rem;
    font-weight: 500;
    line-height: 1.4;
    margin-bottom: 0.25rem;
}

.adr-diagram-indicator {
    font-size: 0.75rem;
    color: var(--text-muted);
    font-style: italic;
}

.sidebar-footer {
    padding: 1rem;
    border-top: 1px solid var(--border-color);
    background: var(--bg-primary);
    text-align: center;
}

.project-info {
    font-size: 0.875rem;
    color: var(--text-secondary);
    margin-bottom: 0.5rem;
}

.footer-link {
    color: var(--primary-color);
    text-decoration: none;
    font-size: 0.875rem;
}

/* ===== MAIN CONTENT ===== */
.breadcrumb {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.875rem;
}

.breadcrumb-home {
    color: var(--text-secondary);
    text-decoration: none;
}

.breadcrumb-home:hover {
    color: var(--primary-color);
}

.breadcrumb-separator {
    color: var(--text-muted);
}

.breadcrumb-category {
    color: var(--text-secondary);
}

.breadcrumb-current {
    color: var(--text-primary);
    font-weight: 500;
}

.content-actions {
    display: flex;
    gap: 0.5rem;
}

.action-btn {
    padding: 0.5rem;
    border: 1px solid var(--border-color);
    background: var(--bg-primary);
    color: var(--text-secondary);
    border-radius: var(--radius-sm);
    cursor: pointer;
    font-size: 1rem;
    transition: all 0.2s ease;
}

.action-btn:hover {
    background: var(--bg-secondary);
    color: var(--text-primary);
}

/* ===== TYPOGRAPHY ===== */
h1, h2, h3, h4, h5, h6 {
    margin-bottom: 1rem;
    color: var(--text-primary);
    font-weight: 600;
    line-height: 1.3;
}

h1 { font-size: 2.25rem; }
h2 { font-size: 1.875rem; }
h3 { font-size: 1.5rem; }
h4 { font-size: 1.25rem; }

p {
    margin-bottom: 1rem;
    color: var(--text-primary);
}

a {
    color: var(--primary-color);
    text-decoration: none;
}

a:hover {
    text-decoration: underline;
}

/* ===== STATUS BADGES ===== */
.status-badge {
    display: inline-flex;
    align-items: center;
    gap: 0.25rem;
    padding: 0.25rem 0.5rem;
    border-radius: var(--radius-sm);
    font-size: 0.75rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
}

.status-badge.accepted { 
    background: #dcfce7; 
    color: #166534; 
    border: 1px solid #bbf7d0;
}
.status-badge.proposed { 
    background: #fef3c7; 
    color: #92400e; 
    border: 1px solid #fde68a;
}
.status-badge.deprecated { 
    background: #fee2e2; 
    color: #991b1b; 
    border: 1px solid #fecaca;
}
.status-badge.superseded { 
    background: #ede9fe; 
    color: #5b21b6; 
    border: 1px solid #ddd6fe;
}

@media (prefers-color-scheme: dark) {
    .status-badge.accepted { 
        background: #064e3b; 
        color: #bbf7d0; 
        border: 1px solid #166534;
    }
    .status-badge.proposed { 
        background: #78350f; 
        color: #fde68a; 
        border: 1px solid #92400e;
    }
    .status-badge.deprecated { 
        background: #7f1d1d; 
        color: #fecaca; 
        border: 1px solid #991b1b;
    }
    .status-badge.superseded { 
        background: #3730a3; 
        color: #ddd6fe; 
        border: 1px solid #5b21b6;
    }
}

/* ===== LISTS ===== */
.adr-list {
    list-style: none;
    margin-bottom: 1rem;
}

.adr-ul {
    list-style: disc;
    margin-left: 1.5rem;
    margin-bottom: 1rem;
}

.adr-ol {
    list-style: decimal;
    margin-left: 1.5rem;
    margin-bottom: 1rem;
}

.adr-list li, .adr-ul li, .adr-ol li {
    margin-bottom: 0.5rem;
    line-height: 1.6;
}

/* ===== RESPONSIVE ===== */
@media (max-width: 768px) {
    .sidebar {
        position: fixed;
        transform: translateX(-100%);
        transition: transform 0.3s ease;
        z-index: 1000;
    }
    
    .sidebar.open {
        transform: translateX(0);
    }
    
    .main-content {
        margin-left: 0;
    }
    
    .adr-content {
        padding: 1rem;
    }
    
    .main-header {
        padding: 1rem;
    }
}

/* ===== WELCOME PAGE ===== */
.welcome-screen {
    max-width: 800px;
    margin: 0 auto;
    text-align: center;
}

.welcome-title {
    font-size: 3rem;
    margin-bottom: 1rem;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 1rem;
}

.welcome-subtitle {
    font-size: 1.25rem;
    color: var(--text-secondary);
    margin-bottom: 2rem;
}

.welcome-description {
    margin-bottom: 3rem;
}

.stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
    gap: 1rem;
    margin: 2rem 0;
}

.stat-card {
    background: var(--bg-secondary);
    padding: 1.5rem;
    border-radius: var(--radius-lg);
    border: 1px solid var(--border-color);
}

.stat-number {
    display: block;
    font-size: 2rem;
    font-weight: 700;
    color: var(--primary-color);
    margin-bottom: 0.5rem;
}

.stat-label {
    color: var(--text-secondary);
    font-size: 0.875rem;
}

.getting-started {
    text-align: left;
    margin-bottom: 3rem;
}

.getting-started h2 {
    text-align: center;
    margin-bottom: 1.5rem;
}

.getting-started ul {
    list-style: none;
    max-width: 500px;
    margin: 0 auto 2rem;
}

.getting-started li {
    padding: 0.5rem 0;
    padding-left: 1.5rem;
    position: relative;
}

.getting-started li::before {
    content: "→";
    position: absolute;
    left: 0;
    color: var(--primary-color);
    font-weight: bold;
}

.welcome-actions {
    display: flex;
    gap: 1rem;
    justify-content: center;
    flex-wrap: wrap;
}

.action-button {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.75rem 1.5rem;
    border-radius: var(--radius-md);
    text-decoration: none;
    font-weight: 500;
    transition: all 0.2s ease;
}

.action-button.primary {
    background: var(--primary-color);
    color: white;
}

.action-button.primary:hover {
    background: var(--primary-hover);
    text-decoration: none;
}

.action-button.secondary {
    background: var(--bg-secondary);
    color: var(--text-primary);
    border: 1px solid var(--border-color);
}

.action-button.secondary:hover {
    background: var(--bg-tertiary);
    text-decoration: none;
}

.recent-adrs {
    text-align: left;
}

.recent-adrs h2 {
    text-align: center;
    margin-bottom: 2rem;
}

.adr-cards {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 1rem;
}

.adr-card {
    background: var(--bg-secondary);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-lg);
    padding: 1.5rem;
    transition: all 0.2s ease;
}

.adr-card:hover {
    box-shadow: var(--shadow-md);
    border-color: var(--primary-color);
}

.adr-card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
}

.adr-card-number {
    font-family: 'JetBrains Mono', 'Fira Code', monospace;
    font-size: 0.875rem;
    color: var(--text-muted);
    font-weight: 600;
}

.adr-card-status {
    font-size: 0.75rem;
}

.adr-card-title {
    margin-bottom: 1rem;
}

.adr-card-title a {
    color: var(--text-primary);
    text-decoration: none;
    font-weight: 600;
}

.adr-card-title a:hover {
    color: var(--primary-color);
}

.adr-card-diagram {
    font-size: 0.875rem;
    color: var(--text-muted);
    font-style: italic;
}

/* ===== ADR DOCUMENT STYLES ===== */
.adr-document {
    max-width: 4xl;
    margin: 0 auto;
}

.adr-header {
    text-align: center;
    margin-bottom: 3rem;
    padding-bottom: 2rem;
    border-bottom: 2px solid var(--border-color);
}

.adr-meta {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 2rem;
    margin-bottom: 2rem;
    flex-wrap: wrap;
}

.adr-number-large {
    font-family: 'JetBrains Mono', 'Fira Code', monospace;
    font-size: 2rem;
    font-weight: 700;
    color: var(--text-muted);
}

.adr-status-container {
    display: flex;
    align-items: center;
    gap: 1rem;
    flex-wrap: wrap;
}

.diagram-badge {
    padding: 0.5rem 1rem;
    background: var(--bg-secondary);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-md);
    font-size: 0.875rem;
    color: var(--text-secondary);
}

.adr-title {
    font-size: 2.5rem;
    margin: 0;
    color: var(--text-primary);
}

.adr-content-wrapper {
    display: grid;
    grid-template-columns: 250px 1fr;
    gap: 3rem;
    margin-bottom: 3rem;
}

.table-of-contents {
    position: sticky;
    top: 6rem;
    height: fit-content;
    background: var(--bg-secondary);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-lg);
    padding: 1.5rem;
}

.table-of-contents h3 {
    font-size: 1rem;
    margin-bottom: 1rem;
    color: var(--text-secondary);
}

#toc-list {
    list-style: none;
}

.toc-level-2 {
    margin-bottom: 0.5rem;
}

.toc-level-3 {
    margin-left: 1rem;
    margin-bottom: 0.25rem;
}

.toc-level-4 {
    margin-left: 2rem;
    margin-bottom: 0.25rem;
}

.toc-link {
    color: var(--text-secondary);
    text-decoration: none;
    font-size: 0.875rem;
    display: block;
    padding: 0.25rem 0;
    border-radius: var(--radius-sm);
    transition: all 0.2s ease;
}

.toc-link:hover,
.toc-link.active {
    color: var(--primary-color);
    background: var(--bg-tertiary);
    padding-left: 0.5rem;
}

.adr-main-content {
    max-width: none;
}

.adr-paragraph {
    margin-bottom: 1.5rem;
    line-height: 1.7;
}

.adr-heading {
    margin-top: 2rem;
    margin-bottom: 1rem;
    scroll-margin-top: 6rem;
}

.adr-h1 { 
    font-size: 2.25rem;
    border-bottom: 2px solid var(--border-color);
    padding-bottom: 0.5rem;
}

.adr-h2 { 
    font-size: 1.875rem;
    color: var(--primary-color);
}

.adr-h3 { 
    font-size: 1.5rem;
    color: var(--text-primary);
}

.adr-table {
    width: 100%;
    border-collapse: collapse;
    margin: 2rem 0;
    background: var(--bg-primary);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-lg);
    overflow: hidden;
}

.adr-table th,
.adr-table td {
    padding: 1rem;
    text-align: left;
    border-bottom: 1px solid var(--border-color);
}

.adr-table th {
    background: var(--bg-secondary);
    font-weight: 600;
    color: var(--text-primary);
}

.adr-table tr:last-child td {
    border-bottom: none;
}

.adr-table code {
    background: var(--bg-tertiary);
    padding: 0.25rem 0.5rem;
    border-radius: var(--radius-sm);
    font-family: 'JetBrains Mono', 'Fira Code', monospace;
    font-size: 0.875rem;
}

.adr-navigation {
    display: flex;
    justify-content: space-between;
    gap: 2rem;
    margin: 3rem 0;
    padding: 2rem 0;
    border-top: 2px solid var(--border-color);
}

.nav-link {
    flex: 1;
    max-width: 300px;
    padding: 1.5rem;
    background: var(--bg-secondary);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-lg);
    text-decoration: none;
    transition: all 0.2s ease;
}

.nav-link:hover {
    border-color: var(--primary-color);
    box-shadow: var(--shadow-md);
    text-decoration: none;
}

.nav-previous {
    text-align: left;
}

.nav-next {
    text-align: right;
    margin-left: auto;
}

.nav-direction {
    display: block;
    font-size: 0.875rem;
    color: var(--text-muted);
    margin-bottom: 0.5rem;
}

.nav-title {
    display: block;
    font-weight: 600;
    color: var(--text-primary);
    line-height: 1.4;
}

.adr-footer {
    text-align: center;
    padding: 2rem 0;
    border-top: 1px solid var(--border-color);
    color: var(--text-secondary);
}

.adr-meta-info {
    margin-bottom: 1rem;
    font-size: 0.875rem;
}

.adr-source a {
    color: var(--primary-color);
    text-decoration: none;
    font-size: 0.875rem;
}

.adr-source a:hover {
    text-decoration: underline;
}

/* ===== RESPONSIVE FOR ADR CONTENT ===== */
@media (max-width: 1024px) {
    .adr-content-wrapper {
        grid-template-columns: 1fr;
        gap: 2rem;
    }
    
    .table-of-contents {
        position: static;
        order: -1;
    }
}

@media (max-width: 768px) {
    .adr-meta {
        flex-direction: column;
        gap: 1rem;
    }
    
    .adr-number-large {
        font-size: 1.5rem;
    }
    
    .adr-title {
        font-size: 2rem;
    }
    
    .adr-navigation {
        flex-direction: column;
    }
    
    .nav-link {
        max-width: none;
    }
    
    .nav-next {
        margin-left: 0;
        text-align: left;
    }
}
`

	if err := os.WriteFile(filepath.Join(cssDir, "main.css"), []byte(css), 0644); err != nil {
		return err
	}

	// Create enhanced JavaScript
	js := `// ADR Generator - Enhanced JavaScript
console.log('ADR site loaded with enhanced features');

// Enhanced search functionality
function initializeSearch() {
    const searchInput = document.getElementById('search-input');
    const searchClear = document.getElementById('search-clear');
    const statusFilter = document.getElementById('status-filter');
    const adrItems = document.querySelectorAll('.adr-item');
    
    if (!searchInput) return;
    
    function performSearch() {
        const query = searchInput.value.toLowerCase();
        const selectedStatus = statusFilter ? statusFilter.value.toLowerCase() : '';
        
        adrItems.forEach(item => {
            const text = item.textContent.toLowerCase();
            const statusElement = item.querySelector('.adr-status');
            const itemStatus = statusElement ? statusElement.textContent.toLowerCase() : '';
            
            const matchesQuery = !query || text.includes(query);
            const matchesStatus = !selectedStatus || itemStatus.includes(selectedStatus);
            
            item.style.display = (matchesQuery && matchesStatus) ? 'block' : 'none';
        });
    }
    
    searchInput.addEventListener('input', performSearch);
    
    if (statusFilter) {
        statusFilter.addEventListener('change', performSearch);
    }
    
    if (searchClear) {
        searchClear.addEventListener('click', function() {
            searchInput.value = '';
            if (statusFilter) statusFilter.value = '';
            performSearch();
            searchInput.focus();
        });
    }
}

// Mobile sidebar toggle
function initializeMobileSidebar() {
    const sidebar = document.querySelector('.sidebar');
    const mainContent = document.querySelector('.main-content');
    
    if (!sidebar || window.innerWidth > 768) return;
    
    // Create mobile menu button
    const menuButton = document.createElement('button');
    menuButton.className = 'mobile-menu-btn';
    menuButton.innerHTML = '☰';
    menuButton.style.cssText = ` + "`" + `
        position: fixed;
        top: 1rem;
        left: 1rem;
        z-index: 1001;
        background: #5e6ad2;
        color: white;
        border: none;
        padding: 0.5rem;
        border-radius: 0.375rem;
        font-size: 1.25rem;
        cursor: pointer;
        display: block;
    ` + "`" + `;
    
    document.body.appendChild(menuButton);
    
    menuButton.addEventListener('click', function() {
        sidebar.classList.toggle('open');
    });
    
    // Close sidebar when clicking outside
    document.addEventListener('click', function(e) {
        if (!sidebar.contains(e.target) && !menuButton.contains(e.target)) {
            sidebar.classList.remove('open');
        }
    });
}

// Enhanced keyboard navigation
function initializeKeyboardNavigation() {
    document.addEventListener('keydown', function(e) {
        // Ctrl/Cmd + K to focus search
        if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
            e.preventDefault();
            const searchInput = document.getElementById('search-input');
            if (searchInput) {
                searchInput.focus();
                searchInput.select();
            }
        }
        
        // Escape to clear search
        if (e.key === 'Escape') {
            const searchInput = document.getElementById('search-input');
            const searchClear = document.getElementById('search-clear');
            if (searchInput && searchInput.value) {
                if (searchClear) {
                    searchClear.click();
                } else {
                    searchInput.value = '';
                    searchInput.dispatchEvent(new Event('input'));
                }
            }
        }
    });
}

// Smooth scrolling for anchor links
function initializeSmoothScrolling() {
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function(e) {
            e.preventDefault();
            const target = document.querySelector(this.getAttribute('href'));
            if (target) {
                target.scrollIntoView({
                    behavior: 'smooth',
                    block: 'start'
                });
            }
        });
    });
}

// Initialize all features when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    initializeSearch();
    initializeMobileSidebar();
    initializeKeyboardNavigation();
    initializeSmoothScrolling();
    
    console.log('All ADR site features initialized');
});

// Handle window resize for mobile sidebar
window.addEventListener('resize', function() {
    if (window.innerWidth > 768) {
        const sidebar = document.querySelector('.sidebar');
        const menuButton = document.querySelector('.mobile-menu-btn');
        if (sidebar) sidebar.classList.remove('open');
        if (menuButton) menuButton.style.display = 'none';
    } else {
        const menuButton = document.querySelector('.mobile-menu-btn');
        if (menuButton) menuButton.style.display = 'block';
    }
});
`

	if err := os.WriteFile(filepath.Join(jsDir, "main.js"), []byte(js), 0644); err != nil {
		return err
	}

	g.stats.AssetCount += 2
	return nil
}

// generateSearchIndex creates a search index JSON file
func (g *Generator) generateSearchIndex() error {
	type SearchItem struct {
		Number      string `json:"number"`
		Title       string `json:"title"`
		Status      string `json:"status"`
		Content     string `json:"content"`
		DiagramType string `json:"diagramType"`
		URL         string `json:"url"`
	}

	var searchItems []SearchItem

	for _, adr := range g.adrs {
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

	searchIndex := map[string]interface{}{
		"generated": fmt.Sprintf("%d", g.stats.ADRCount),
		"items":     searchItems,
	}

	data, err := json.MarshalIndent(searchIndex, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal search index: %w", err)
	}

	indexPath := filepath.Join(g.config.OutputDirectory, "search-index.json")
	if err := os.WriteFile(indexPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write search index: %w", err)
	}

	return nil
}
