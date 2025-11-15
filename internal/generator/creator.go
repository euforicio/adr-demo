package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ADRConfig holds the ADR creation configuration
type ADRConfig struct {
	Title    string
	Status   string
	Category string
	Force    bool
	Verbose  bool
}

// ADRCreator handles creating new ADRs
type ADRCreator struct {
	config *ADRConfig
}

// NewADRCreator creates a new ADR creator
func NewADRCreator(config *ADRConfig) *ADRCreator {
	return &ADRCreator{
		config: config,
	}
}

// Create creates a new ADR file
func (c *ADRCreator) Create() (string, error) {
	// Get next ADR number
	nextNumber, err := c.getNextADRNumber()
	if err != nil {
		return "", fmt.Errorf("failed to determine next ADR number: %w", err)
	}

	// Create filename and path
	filename := c.createFilename(nextNumber, c.config.Title)

	// Use flat structure following ADR spec
	adrDir := "adr"

	// Create directory if it doesn't exist
	if err := os.MkdirAll(adrDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create ADR directory: %w", err)
	}

	filePath := filepath.Join(adrDir, filename)

	// Check if file already exists
	if !c.config.Force {
		if _, err := os.Stat(filePath); err == nil {
			return "", fmt.Errorf("ADR file already exists: %s (use --force to overwrite)", filename)
		}
	}

	// Generate ADR content
	content := c.generateADRContent(nextNumber, c.config.Title, c.config.Status)

	// Write file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write ADR file: %w", err)
	}

	// Return the relative path from the project root
	return filepath.Join(strings.TrimPrefix(adrDir, "./"), filename), nil
}

// getNextADRNumber determines the next available ADR number from flat structure
func (c *ADRCreator) getNextADRNumber() (int, error) {
	adrDir := "adr"

	// If directory doesn't exist, start with 0001
	if _, err := os.Stat(adrDir); os.IsNotExist(err) {
		return 1, nil
	}

	maxNumber := 0

	// Read all files in the ADR directory (flat structure)
	files, err := os.ReadDir(adrDir)
	if err != nil {
		return 0, fmt.Errorf("failed to read ADR directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Only process .md files
		if !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		// Skip template
		if file.Name() == "template.md" {
			continue
		}

		// Extract number from filename
		if len(file.Name()) >= 4 {
			numStr := file.Name()[:4]
			if num, err := strconv.Atoi(numStr); err == nil {
				if num > maxNumber {
					maxNumber = num
				}
			}
		}
	}

	return maxNumber + 1, nil
}

// createFilename creates a filename from title
func (c *ADRCreator) createFilename(number int, title string) string {
	// Convert title to kebab-case
	kebabTitle := c.toKebabCase(title)
	return fmt.Sprintf("%04d-%s.md", number, kebabTitle)
}

// toKebabCase converts a string to kebab-case
func (c *ADRCreator) toKebabCase(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace spaces and special characters with hyphens
	re := regexp.MustCompile(`[^a-z0-9]+`)
	s = re.ReplaceAllString(s, "-")

	// Remove leading/trailing hyphens
	s = strings.Trim(s, "-")

	// Collapse multiple hyphens
	re = regexp.MustCompile(`-+`)
	s = re.ReplaceAllString(s, "-")

	return s
}

// generateADRContent generates the content for a new ADR
func (c *ADRCreator) generateADRContent(number int, title, status string) string {
	now := time.Now()

	template := `# %s

## Status

%s

## Context

*Describe the context and problem statement that led to this decision.*

The issue motivating this decision, and any context that influences or constrains the decision.

## Decision

*Describe the decision that was made.*

We will...

### Rationale

*Explain why this decision was made.*

### Alternatives Considered

*List other options that were considered and why they were not chosen.*

## Consequences

### Positive

- *List positive consequences of this decision*

### Negative

- *List negative consequences of this decision*

### Neutral

- *List neutral consequences that should be noted*

## Implementation

### Next Steps

- [ ] Task 1
- [ ] Task 2
- [ ] Task 3

### Timeline

*Describe the implementation timeline and milestones.*

## Related Decisions

*Link to related ADRs or decisions.*

---

*This ADR was created on %s*`

	return fmt.Sprintf(template,
		title,
		status,
		now.Format("January 2, 2006"))
}
