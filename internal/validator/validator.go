package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Config holds the validator configuration
type Config struct {
	Strict  bool
	Fix     bool
	Verbose bool
}

// ValidationResult holds the validation results
type ValidationResult struct {
	Issues       []Issue
	FileCount    int
	DiagramCount int
	ErrorCount   int
	WarningCount int
	FixCount     int
}

// Issue represents a validation issue
type Issue struct {
	File    string
	Line    int
	Level   string // "error" or "warning"
	Message string
}

// Validator validates ADR files
type Validator struct {
	config *Config
}

// New creates a new validator
func New(config *Config) *Validator {
	return &Validator{
		config: config,
	}
}

// ValidateAll validates all ADR files
func (v *Validator) ValidateAll() (*ValidationResult, error) {
	result := &ValidationResult{
		Issues: make([]Issue, 0),
	}

	// Find all ADR files
	adrDir := "adr"
	entries, err := os.ReadDir(adrDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read ADR directory: %w", err)
	}

	adrFiles := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		// Skip template file
		if entry.Name() == "template.md" {
			continue
		}

		// Check filename format
		if v.isValidADRFilename(entry.Name()) {
			adrFiles = append(adrFiles, entry.Name())
		} else {
			result.Issues = append(result.Issues, Issue{
				File:    entry.Name(),
				Line:    0,
				Level:   "error",
				Message: "Invalid ADR filename format. Expected: NNNN-kebab-case-title.md",
			})
			result.ErrorCount++
		}
	}

	result.FileCount = len(adrFiles)

	// Validate sequential numbering
	if err := v.validateSequentialNumbering(adrFiles, result); err != nil {
		return nil, err
	}

	// Validate each file
	for _, filename := range adrFiles {
		filePath := filepath.Join(adrDir, filename)
		if err := v.validateFile(filePath, result); err != nil {
			return nil, fmt.Errorf("failed to validate %s: %w", filename, err)
		}
	}

	return result, nil
}

// HasErrors returns true if there are validation errors
func (r *ValidationResult) HasErrors() bool {
	return r.ErrorCount > 0
}

// isValidADRFilename checks if filename follows NNNN-kebab-case.md format
func (v *Validator) isValidADRFilename(filename string) bool {
	pattern := `^[0-9]{4}-[a-z0-9-]+\.md$`
	matched, _ := regexp.MatchString(pattern, filename)
	return matched
}

// validateSequentialNumbering ensures ADR numbers are sequential
func (v *Validator) validateSequentialNumbering(filenames []string, result *ValidationResult) error {
	numbers := make([]int, 0)

	for _, filename := range filenames {
		if len(filename) >= 4 {
			numStr := filename[:4]
			num, err := strconv.Atoi(numStr)
			if err != nil {
				result.Issues = append(result.Issues, Issue{
					File:    filename,
					Line:    0,
					Level:   "error",
					Message: "Invalid ADR number format",
				})
				result.ErrorCount++
				continue
			}
			numbers = append(numbers, num)
		}
	}

	// Check for gaps in numbering
	for i := 0; i < len(numbers)-1; i++ {
		if numbers[i+1] != numbers[i]+1 {
			result.Issues = append(result.Issues, Issue{
				File:    fmt.Sprintf("ADR-%04d", numbers[i+1]),
				Line:    0,
				Level:   "error",
				Message: fmt.Sprintf("Gap in ADR numbering: %04d follows %04d", numbers[i+1], numbers[i]),
			})
			result.ErrorCount++
		}
	}

	return nil
}

// validateFile validates a single ADR file
func (v *Validator) validateFile(filePath string, result *ValidationResult) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	filename := filepath.Base(filePath)
	lines := strings.Split(string(content), "\n")

	// Check for required sections
	v.validateRequiredSections(filename, lines, result)

	// Check heading hierarchy
	v.validateHeadingHierarchy(filename, lines, result)

	// Check for Mermaid diagrams
	diagramCount := v.validateMermaidDiagrams(filename, lines, result)
	result.DiagramCount += diagramCount

	// Strict mode checks
	if v.config.Strict {
		v.validateStrictRules(filename, lines, result)
	}

	return nil
}

// validateRequiredSections checks for required ADR sections
func (v *Validator) validateRequiredSections(filename string, lines []string, result *ValidationResult) {
	requiredSections := []string{"Status", "Context", "Decision", "Consequences"}
	foundSections := make(map[string]bool)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			sectionName := strings.TrimSpace(trimmed[3:])
			for _, required := range requiredSections {
				if strings.EqualFold(sectionName, required) {
					foundSections[required] = true
				}
			}
		}
	}

	for _, section := range requiredSections {
		if !foundSections[section] {
			result.Issues = append(result.Issues, Issue{
				File:    filename,
				Line:    0,
				Level:   "error",
				Message: fmt.Sprintf("Missing required section: %s", section),
			})
			result.ErrorCount++
		}
	}
}

// validateHeadingHierarchy checks heading structure
func (v *Validator) validateHeadingHierarchy(filename string, lines []string, result *ValidationResult) {
	hasTitle := false

	for lineNum, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for main title (# Title)
		if strings.HasPrefix(trimmed, "# ") && !hasTitle {
			hasTitle = true
			title := strings.TrimSpace(trimmed[2:])
			if title == "" {
				result.Issues = append(result.Issues, Issue{
					File:    filename,
					Line:    lineNum + 1,
					Level:   "error",
					Message: "ADR title cannot be empty",
				})
				result.ErrorCount++
			}
		}

		// Check for multiple H1 headings
		if strings.HasPrefix(trimmed, "# ") && hasTitle && lineNum > 10 {
			result.Issues = append(result.Issues, Issue{
				File:    filename,
				Line:    lineNum + 1,
				Level:   "warning",
				Message: "Multiple H1 headings found. ADRs should have only one main title",
			})
			result.WarningCount++
		}
	}

	if !hasTitle {
		result.Issues = append(result.Issues, Issue{
			File:    filename,
			Line:    0,
			Level:   "error",
			Message: "ADR must have a main title (# Title)",
		})
		result.ErrorCount++
	}
}

// validateMermaidDiagrams validates Mermaid diagram syntax
func (v *Validator) validateMermaidDiagrams(filename string, lines []string, result *ValidationResult) int {
	diagramCount := 0
	inMermaid := false
	mermaidStart := 0

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "```mermaid" {
			inMermaid = true
			mermaidStart = i + 1
			diagramCount++
		} else if trimmed == "```" && inMermaid {
			inMermaid = false
		}
	}

	// Check for unclosed Mermaid blocks
	if inMermaid {
		result.Issues = append(result.Issues, Issue{
			File:    filename,
			Line:    mermaidStart,
			Level:   "error",
			Message: "Unclosed Mermaid diagram block",
		})
		result.ErrorCount++
	}

	return diagramCount
}

// validateStrictRules applies strict validation rules
func (v *Validator) validateStrictRules(filename string, lines []string, result *ValidationResult) {
	for lineNum, line := range lines {
		// Check line length
		if len(line) > 120 {
			result.Issues = append(result.Issues, Issue{
				File:    filename,
				Line:    lineNum + 1,
				Level:   "warning",
				Message: "Line exceeds 120 characters",
			})
			result.WarningCount++
		}

		// Check for trailing whitespace
		if len(line) > 0 && (line[len(line)-1] == ' ' || line[len(line)-1] == '\t') {
			result.Issues = append(result.Issues, Issue{
				File:    filename,
				Line:    lineNum + 1,
				Level:   "warning",
				Message: "Trailing whitespace",
			})
			result.WarningCount++
		}
	}
}
