# ADR Demo Repository Makefile
# Provides local development commands for managing ADRs

.PHONY: help generate-index serve validate clean lint test

# Default target
help: ## Show this help message
	@echo "ADR Demo Repository Commands"
	@echo "=============================="
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Examples:"
	@echo "  make generate-index  # Generate adr-index.json"
	@echo "  make serve          # Start local development server"
	@echo "  make validate       # Validate all ADRs"

generate-index: ## Generate the adr-index.json file
	@echo "🔄 Generating ADR index..."
	@./scripts/generate-adr-index.sh
	@echo "✅ ADR index generated successfully"

serve: generate-index ## Start local development server with auto-generated index
	@echo "🚀 Starting local development server..."
	@echo "📂 Serving from docs/ directory"
	@echo "🌐 Open http://localhost:8080 in your browser"
	@echo "⏹️  Press Ctrl+C to stop"
	@cd docs && python3 -m http.server 8080

serve-simple: ## Start simple server without regenerating index
	@echo "🚀 Starting simple development server..."
	@cd docs && python3 -m http.server 8080

validate: ## Validate ADR structure and content
	@echo "🔍 Validating ADR structure..."
	@./scripts/validate-adrs.sh

validate-links: ## Check for broken links in ADRs
	@echo "🔗 Checking for broken links..."
	@./scripts/check-links.sh

lint: ## Lint markdown files
	@echo "📝 Linting markdown files..."
	@if command -v markdownlint >/dev/null 2>&1; then \
		markdownlint adr/*.md README.md; \
	else \
		echo "⚠️  markdownlint not installed. Install with: npm install -g markdownlint-cli"; \
	fi

new-adr: ## Create a new ADR from template
	@echo "📄 Creating new ADR..."
	@./scripts/new-adr.sh

test: validate lint ## Run all tests and validation
	@echo "🧪 Running all tests..."
	@make generate-index
	@echo "✅ All tests passed"

clean: ## Clean generated files
	@echo "🧹 Cleaning generated files..."
	@rm -f adr-index.json
	@echo "✅ Cleaned generated files"

install-deps: ## Install development dependencies
	@echo "📦 Installing development dependencies..."
	@if command -v npm >/dev/null 2>&1; then \
		npm install -g markdownlint-cli; \
		echo "✅ markdownlint installed"; \
	else \
		echo "⚠️  npm not found. Please install Node.js"; \
	fi

watch: ## Watch for changes and regenerate index
	@echo "👀 Watching for ADR changes..."
	@if command -v fswatch >/dev/null 2>&1; then \
		fswatch -o adr/ | xargs -n1 -I{} make generate-index; \
	elif command -v inotifywait >/dev/null 2>&1; then \
		while inotifywait -e modify,create,delete adr/; do make generate-index; done; \
	else \
		echo "⚠️  File watching requires fswatch (macOS) or inotify-tools (Linux)"; \
		echo "   Install with: brew install fswatch (macOS) or apt-get install inotify-tools (Linux)"; \
	fi

dev: ## Start development mode with file watching and server
	@echo "🛠️  Starting development mode..."
	@make generate-index
	@echo "🚀 Starting server in background..."
	@cd docs && python3 -m http.server 8080 &
	@SERVER_PID=$$!; \
	echo "🌐 Server running at http://localhost:8080"; \
	echo "👀 Watching for file changes..."; \
	echo "⏹️  Press Ctrl+C to stop both server and watcher"; \
	trap "kill $$SERVER_PID" EXIT; \
	make watch

# GitHub Actions simulation
ci-validate: ## Simulate GitHub Actions validation locally
	@echo "🔄 Simulating CI validation..."
	@make clean
	@make generate-index
	@make validate
	@make lint
	@echo "✅ CI validation simulation completed"

# Quick development commands
quick-serve: ## Quick serve with index regeneration (alias for serve)
	@make serve

qs: quick-serve ## Short alias for quick-serve

gi: generate-index ## Short alias for generate-index

# Platform-specific commands
serve-mac: generate-index ## Start server on macOS with open command
	@make serve &
	@sleep 2 && open http://localhost:8080

serve-linux: generate-index ## Start server on Linux with xdg-open
	@make serve &
	@sleep 2 && xdg-open http://localhost:8080

# Documentation
docs: ## Generate documentation
	@echo "📚 Documentation available at:"
	@echo "   README.md - Main project documentation"
	@echo "   adr/ - Architecture Decision Records"
	@echo "   docs/index.html - Web interface"

# GitHub Pages
pages-build: generate-index ## Prepare for GitHub Pages deployment
	@echo "🚀 Preparing for GitHub Pages deployment..."
	@echo "✅ ADR index generated"
	@echo "📄 Files ready in docs/ directory"
	@echo ""
	@echo "💡 To deploy:"
	@echo "   1. Push to main branch"
	@echo "   2. GitHub Actions will auto-deploy to Pages"
	@echo "   3. Site will be available at: https://USERNAME.github.io/REPO-NAME"

pages-test: generate-index serve-simple ## Test GitHub Pages setup locally
	@echo "🧪 Testing GitHub Pages setup locally..."

# Status and info
status: ## Show repository status
	@echo "📊 ADR Repository Status"
	@echo "======================="
	@echo "Total ADRs: $$(find adr -name '[0-9][0-9][0-9][0-9]-*.md' | wc -l | tr -d ' ')"
	@echo "Accepted:   $$(grep -l '^Accepted$$' adr/*.md 2>/dev/null | wc -l | tr -d ' ')"
	@echo "Proposed:   $$(grep -l '^Proposed$$' adr/*.md 2>/dev/null | wc -l | tr -d ' ')"
	@echo "Deprecated: $$(grep -l '^Deprecated$$' adr/*.md 2>/dev/null | wc -l | tr -d ' ')"
	@echo "Superseded: $$(grep -l '^Superseded$$' adr/*.md 2>/dev/null | wc -l | tr -d ' ')"
	@echo ""
	@if [ -f adr-index.json ]; then \
		echo "Index file: ✅ Present (generated $$(jq -r '.generated' adr-index.json 2>/dev/null || echo 'unknown'))"; \
	else \
		echo "Index file: ❌ Missing (run 'make generate-index')"; \
	fi