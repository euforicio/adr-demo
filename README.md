# Architecture Decision Records (ADR) Demo

This repository demonstrates best practices for documenting architectural decisions using Architecture Decision Records (ADRs) following the standards from [adr.github.io](https://adr.github.io/). It features a modern Go-based web interface with interactive diagrams, smart search, and comprehensive GitHub automation for managing ADRs in software development teams.

## 🚀 Quick Start

```bash
# Start the interactive ADR browser
go run main.go serve

# Open in your browser
open http://localhost:8080
```

## What are ADRs?

Architecture Decision Records (ADRs) are **lightweight documents that capture important architectural decisions** made during a project's development. Each ADR documents:

- **The decision made** and its rationale
- **The context** that led to the decision  
- **The consequences** (both positive and negative) of the decision
- **Alternative options** that were considered

ADRs help teams understand the reasoning behind architectural choices and maintain a historical record of decision-making that survives team changes and time.

## Why Use ADRs?

### Key Benefits

- **Knowledge Preservation**: Architectural decisions and their reasoning are documented for future reference
- **Onboarding**: New team members can understand the system's evolution and current state
- **Decision Transparency**: Everyone can see what decisions were made and why
- **Avoid Repeated Discussions**: Settled architectural matters don't need to be re-debated
- **Change Management**: Understanding past decisions helps evaluate future changes
- **Accountability**: Clear record of who decided what and when

### When to Create an ADR

Create an ADR when making decisions that:
- Affect the overall system architecture
- Have long-term consequences
- Involve significant trade-offs
- Impact multiple teams or components
- Introduce new technologies or patterns
- Change existing architectural patterns

## ADR Index

| ADR | Title | Status | C4 Diagram Type |
|-----|-------|--------|-----------------|
| [0001](adr/0001-record-architecture-decisions.md) | Record Architecture Decisions | Accepted | - |
| [0002](adr/0002-establish-architecture-review-board.md) | Establish Architecture Review Board | Accepted | Dynamic |
| [0003](adr/0003-adopt-microservices-architecture.md) | Adopt Microservices Architecture | Accepted | Context |
| [0004](adr/0004-choose-database-per-service.md) | Choose Database Per Service | Accepted | Container |
| [0005](adr/0005-implement-api-gateway-pattern.md) | Implement API Gateway Pattern | Accepted | Component |
| [0006](adr/0006-use-event-driven-communication.md) | Use Event-Driven Communication | Accepted | Dynamic |
| [0007](adr/0007-implement-graphql-api.md) | Implement GraphQL API | Proposed | Component |
| [0008](adr/0008-use-mongodb-for-session-storage.md) | Use MongoDB for Session Storage | Deprecated | Container |
| [0009](adr/0009-use-redis-for-session-storage.md) | Use Redis for Session Storage | Superseded | Container |
| [0010](adr/0010-adopt-hybrid-session-storage.md) | Adopt Hybrid Session Storage | Accepted | Component |

### ADR Status Definitions

- **Proposed**: Under review and discussion
- **Accepted**: Approved and being implemented
- **Deprecated**: No longer recommended but still in use
- **Superseded**: Replaced by a newer decision (link to replacement)

## ADR Process

See [ADR_PROCESS.md](ADR_PROCESS.md) for detailed process documentation including:

- Step-by-step ADR creation workflow
- Architecture Review Board procedures
- GitHub automation and validation
- Template and formatting guidelines
- Best practices and common pitfalls

## Example Project Context

These ADRs document architectural decisions for a fictional e-commerce platform called **"ShopFlow"** to demonstrate realistic scenarios and decision-making processes. The example shows how ADRs can track the evolution of a system from monolith to microservices, including decision lifecycle management.

### ShopFlow Architecture Evolution

1. **Foundation** (ADR-0001, 0002): Established ADR process and governance
2. **Architecture Shift** (ADR-0003): Moved from monolith to microservices  
3. **Data Strategy** (ADR-0004): Implemented database-per-service pattern
4. **API Management** (ADR-0005): Added API gateway for unified access
5. **Communication** (ADR-0006): Adopted event-driven architecture
6. **API Evolution** (ADR-0007): Proposed GraphQL API layer (under review)
7. **Session Storage Journey** (ADR-0008, 0009, 0010): Evolution from MongoDB → Redis → Hybrid approach

### Decision Lifecycle Examples

This repository demonstrates the full ADR lifecycle:

- **⏳ Proposed**: ADR-0007 shows a decision under active review
- **✅ Accepted**: ADRs 0001-0006, 0010 represent current active decisions  
- **🗑️ Deprecated**: ADR-0008 shows how decisions can become outdated
- **⬆️ Superseded**: ADR-0009 demonstrates replacement by better solutions

This progression demonstrates how architectural decisions build upon each other and how ADRs can document complex system evolution over time, including the natural lifecycle of architectural choices.

## Interactive ADR Browser

This repository features a modern Go-based web application that provides an interactive, real-time interface for browsing and exploring ADRs with advanced features.

### Key Features

- 🚀 **Dynamic Rendering**: Markdown files rendered on-the-fly with intelligent caching
- 📱 **Responsive Design**: Optimized for desktop, tablet, and mobile devices  
- 🌓 **Dark Mode**: Smart theme toggle with system preference detection
- 🔍 **Real-time Search**: Instant search across all ADRs with status filters
- 📊 **Status Tracking**: Visual indicators for all ADR lifecycle states
- 🗂️ **Smart Navigation**: Collapsible category groups and flat list views
- 📈 **Interactive Diagrams**: Advanced Mermaid diagram viewer with zoom controls
- ⌨️ **Keyboard Shortcuts**: Quick access with Ctrl/Cmd+K for search, Escape to clear

### Advanced Diagram Features

The ADR browser includes a sophisticated diagram viewer with:

- **🔍 Zoom Controls**: Zoom in/out with buttons, mouse wheel, or trackpad gestures
- **🖱️ Pan & Navigate**: Click and drag to explore large diagrams  
- **📱 Gesture Support**: Natural pinch-to-zoom with zoom-to-cursor positioning
- **🎯 Fullscreen Mode**: Click any diagram for detailed fullscreen view
- **🌙 Theme Aware**: Proper text colors in both light and dark modes
- **⚡ Smooth Performance**: Optimized interactions with gesture detection

### Technical Architecture

The ADR browser uses a modern, efficient architecture:

#### **Smart Caching System**
- **SHA-256 Content Hashing**: Tracks file changes for intelligent cache invalidation
- **In-Memory Rendering Cache**: Avoids re-processing unchanged markdown files
- **Dynamic Cache Management**: Automatically refreshes when source files change
- **Zero-Downtime Updates**: Content updates without server restart

#### **Performance Optimizations**
- **On-Demand Rendering**: Only processes requested ADRs, not entire site
- **Efficient Routing**: Fast URL pattern matching and request handling
- **Template Optimization**: Reusable template compilation and execution
- **Asset Serving**: Static assets served efficiently with proper headers

#### **Modern Web Features**
- **Tailwind CSS Integration**: Utility-first styling with custom prose plugin
- **GitHub-Style Syntax Highlighting**: Prism.js with proper dark mode support
- **Responsive Design**: Mobile-first approach with flexbox layouts
- **Progressive Enhancement**: Works without JavaScript, enhanced with JS

## GitHub Automation

This repository includes comprehensive GitHub Actions workflows that automate and validate the ADR process, ensuring consistency and quality across all architectural decisions.

### GitHub Actions Workflows

#### 🔍 ADR Validation (`adr-validation.yml`)
**Triggers**: Pull requests and pushes affecting ADR files  
**Purpose**: Comprehensive validation of ADR structure and consistency

**What it checks:**
- **Sequential Numbering**: Ensures ADRs are numbered sequentially (0001, 0002, etc.)
- **No Duplicates**: Prevents multiple ADRs with the same number
- **Required Sections**: Validates presence of Status, Context, Decision, Consequences
- **Valid Status**: Ensures status is one of: Proposed, Accepted, Deprecated, Superseded  
- **README Consistency**: Verifies all ADRs are referenced in README index
- **Broken Links**: Checks for broken internal markdown links
- **Markdown Linting**: Validates markdown formatting

#### 🔄 Auto-Update Index (`adr-index-update.yml`)
**Triggers**: Pushes to main branch with new ADR files  
**Purpose**: Automatically maintains the README ADR index

**What it does:**
- **Generates Index**: Creates up-to-date table with ADR numbers, titles, status
- **Detects C4 Diagrams**: Automatically identifies diagram types (Context, Container, Component, Dynamic)
- **Smart Updates**: Only commits changes when index actually changes
- **Auto-Commit**: Commits updates with descriptive message

#### 📋 Template & PR Check (`adr-template-check.yml`)
**Triggers**: Pull requests with ADR changes  
**Purpose**: Validates new ADRs follow template and provides PR feedback

**What it validates:**
- **Template Compliance**: Ensures new ADRs follow the standard template
- **No Placeholder Text**: Checks that template placeholders have been replaced
- **Meaningful Content**: Validates sections contain actual content, not just headers
- **Naming Convention**: Enforces `0000-kebab-case-title.md` format
- **Numbering Conflicts**: Prevents conflicts with existing ADR numbers
- **PR Comments**: Adds helpful validation results as PR comments

#### 🚀 Build & Deploy (`build-and-test.yml`)
**Triggers**: Pushes to main branch, pull requests  
**Purpose**: Build and test the Go ADR server

**What it does:**
- **Go Build**: Compiles the ADR server application
- **Run Tests**: Executes unit and integration tests
- **Static Analysis**: Runs Go vet and other static analysis tools
- **Generate Artifacts**: Creates deployable binaries
- **Validate Functionality**: Tests server startup and basic functionality

#### ✅ Setup Validation (`validate-setup.yml`)
**Triggers**: Pull requests, manual dispatch  
**Purpose**: Comprehensive validation of repository setup

**What it validates:**
- **File Structure**: Verifies all critical files are present
- **ADR Format**: Runs full ADR structure validation
- **Index Generation**: Tests the ADR index generation process
- **Web App**: Basic validation of JavaScript and CSS files
- **Statistics**: Provides repository health metrics

### GitHub Integration Features

#### Pull Request Workflow
1. **Create ADR**: Developer creates new ADR following template
2. **Open PR**: Submit ADR as pull request
3. **Automatic Validation**: GitHub Actions validate structure, numbering, format
4. **PR Comments**: Automated feedback posted to PR
5. **Review Process**: Architecture Review Board reviews via GitHub PR
6. **Merge & Update**: Upon merge, README index is automatically updated

#### Branch Protection (Recommended)
```yaml
# .github/branch-protection.yml (example configuration)
protection_rules:
  main:
    required_status_checks:
      - "ADR Validation"
      - "ADR Pull Request Validation"
    required_reviews: 2
    require_code_owner_reviews: true
    required_reviewers:
      - architecture-review-board
```

---

## Installation & Setup

### Prerequisites

Ensure you have Go installed (version 1.19 or later):

```bash
# Check Go version
go version

# Install Go if needed (macOS with Homebrew)
brew install go

# Or download from https://golang.org/dl/
```

### Quick Setup

```bash
# Clone the repository
git clone https://github.com/your-org/adr-demo.git
cd adr-demo

# Install dependencies
go mod tidy

# Start the ADR browser
go run main.go serve

# Open in your browser
open http://localhost:8080
```

### Command Options

```bash
# Build static files only
go run main.go build --output-dir ./dist

# Serve with custom configuration  
go run main.go serve --port 3000 --host 0.0.0.0

# Verbose output for debugging
go run main.go serve --verbose
```

### Environment Configuration

```bash
# Set custom port
export ADR_PORT=3000

# Enable verbose logging
export ADR_VERBOSE=true

# Custom ADR directory (default: ./adr)
export ADR_DIRECTORY=./docs/adrs
```

## Development Guide

### GitHub Actions Setup

1. **Enable Workflows**: GitHub Actions are automatically enabled when you add workflow files to `.github/workflows/`

2. **Configure Permissions**: Ensure your repository has these permissions:
   - **Contents**: Write (for committing README updates)
   - **Pull Requests**: Write (for adding PR comments)

3. **Set Up Branch Protection** (Optional):
   ```bash
   # Using GitHub CLI
   gh api repos/:owner/:repo/branches/main/protection \
     --method PUT \
     --field required_status_checks='{"strict":true,"contexts":["ADR Validation"]}' \
     --field required_pull_request_reviews='{"required_approving_review_count":1}'
   ```

4. **Configure Code Owners**: Create `.github/CODEOWNERS`:
   ```
   # Require Architecture Review Board approval for ADRs
   adr/ @your-org/architecture-review-board
   ```

### Development Tools

- **ADR Tools**: Command-line tools ([adr-tools](https://github.com/npryce/adr-tools))
- **Mermaid Live Editor**: [mermaid.live](https://mermaid.live) for diagram editing
- **PlantUML**: Alternative for complex diagrams
- **Markdown Linters**: Ensure consistent formatting

### Command Line Helpers

```bash
# Find next ADR number
ls adr/ | grep -E '^[0-9]{4}' | sort -n | tail -1 | sed 's/^\([0-9]*\).*/\1/' | xargs printf "%04d\n" | xargs -I {} expr {} + 1 | xargs printf "%04d\n"

# Create new ADR from template
cp adr/template.md adr/$(next_adr_number)-your-decision-title.md

# Validate local ADRs before committing
find adr -name '[0-9][0-9][0-9][0-9]-*.md' | while read file; do
  echo "Checking $file..."
  grep -q "^## Status" "$file" || echo "Missing Status section"
  grep -q "^## Context" "$file" || echo "Missing Context section" 
  grep -q "^## Decision" "$file" || echo "Missing Decision section"
  grep -q "^## Consequences" "$file" || echo "Missing Consequences section"
done
```

## Deployment Options

### Local Development
```bash
go run main.go serve --port 8080
```

### Production Build
```bash
go build -o adr-server
./adr-server serve --host 0.0.0.0 --port 80
```

### Docker (Optional)
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o adr-server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/adr-server .
COPY --from=builder /app/adr ./adr
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
CMD ["./adr-server", "serve", "--host", "0.0.0.0", "--port", "8080"]
```

## Troubleshooting

### Common Issues

#### ADR Validation Failures
- **Duplicate Numbers**: Use the next sequential number available
- **Missing Sections**: Ensure all required sections are present
- **Invalid Status**: Use only: Proposed, Accepted, Deprecated, Superseded
- **Template Text**: Replace all placeholder text from template

#### GitHub Actions Issues
- **Permission Errors**: Ensure repository has write permissions for contents
- **Workflow Failures**: Check Actions tab for detailed error messages
- **Manual Trigger**: Use workflow_dispatch to manually trigger workflows

#### Server Issues
- **Port Already in Use**: Change port with `--port` flag or `ADR_PORT` environment variable
- **File Not Found**: Ensure ADR files are in the correct directory structure
- **Template Errors**: Check Go template syntax in `templates/` directory

### Debugging Commands

```bash
# Check ADR file structure
find adr -name '[0-9][0-9][0-9][0-9]-*.md' | sort

# Test server locally
go run main.go serve --verbose --port 8081

# Validate ADR format manually
go run main.go build --verbose
```

### Getting Help

- **Issues**: Open an issue on the repository for bug reports
- **Discussions**: Use GitHub Discussions for questions and ideas
- **Documentation**: Refer to [adr.github.io](https://adr.github.io/) for ADR standards

## Example Project Context

These ADRs document architectural decisions for a fictional e-commerce platform called **"ShopFlow"** to demonstrate realistic scenarios and decision-making processes. The example shows how ADRs can track the evolution of a system from monolith to microservices, including decision lifecycle management.

### ShopFlow Architecture Evolution

1. **Foundation** (ADR-0001, 0002): Established ADR process and governance
2. **Architecture Shift** (ADR-0003): Moved from monolith to microservices  
3. **Data Strategy** (ADR-0004): Implemented database-per-service pattern
4. **API Management** (ADR-0005): Added API gateway for unified access
5. **Communication** (ADR-0006): Adopted event-driven architecture
6. **API Evolution** (ADR-0007): Proposed GraphQL API layer (under review)
7. **Session Storage Journey** (ADR-0008, 0009, 0010): Evolution from MongoDB → Redis → Hybrid approach

### Decision Lifecycle Examples

This repository demonstrates the full ADR lifecycle:

- **⏳ Proposed**: ADR-0007 shows a decision under active review
- **✅ Accepted**: ADRs 0001-0006, 0010 represent current active decisions  
- **🗑️ Deprecated**: ADR-0008 shows how decisions can become outdated
- **⬆️ Superseded**: ADR-0009 demonstrates replacement by better solutions

This progression demonstrates how architectural decisions build upon each other and how ADRs can document complex system evolution over time, including the natural lifecycle of architectural choices.

---

## Contributing

We welcome contributions to improve this ADR demo and make it more useful for the community!

### How to Contribute

1. **Fork** this repository
2. **Create** a feature branch for your changes
3. **Make** your improvements (code, documentation, examples)
4. **Test** your changes locally with `go run main.go serve`
5. **Submit** a pull request with a clear description
6. **Engage** with feedback and iterate as needed

### Types of Contributions

- 🐛 **Bug Fixes**: Report and fix issues with the web interface or automation
- 📚 **Documentation**: Improve README, add examples, clarify processes
- 🎨 **UI/UX**: Enhance the web interface, improve responsiveness, add features
- 🔧 **Tooling**: Improve GitHub Actions, add new automation features
- 📋 **ADR Examples**: Add realistic ADRs to demonstrate different scenarios
- 🏗️ **Architecture**: Improve the Go server, caching, or template system

### Development Setup

```bash
# Fork and clone your fork
git clone https://github.com/YOUR-USERNAME/adr-demo.git
cd adr-demo

# Install dependencies
go mod tidy

# Start development server
go run main.go serve --verbose

# Run tests (if any)
go test ./...
```

### Questions or Ideas?

- 🐛 **Bug Reports**: Open an issue with detailed reproduction steps
- 💡 **Feature Requests**: Open an issue to discuss new ideas
- 💬 **General Questions**: Use GitHub Discussions
- 📧 **Private Concerns**: Contact the maintainers directly

---

*This repository serves as both documentation and demonstration of effective ADR practices. Use it as a template for your own projects and adapt the process to fit your team's needs.*