package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// StatusConfig defines the visual representation of an ADR status
type StatusConfig struct {
	Icon     string `yaml:"icon"`
	Color    string `yaml:"color"`
	CSSClass string `yaml:"css_class"`
}

// Config holds the complete ADR tool configuration
type Config struct {
	ADRDirectory      string                  `yaml:"adr_directory"`
	OutputDirectory   string                  `yaml:"output_directory"`
	BaseURL           string                  `yaml:"base_url"`
	DefaultCategory   string                  `yaml:"default_category"`
	AllowedCategories []string                `yaml:"allowed_categories"`
	AllowedStatuses   []string                `yaml:"allowed_statuses"`
	StatusConfig      map[string]StatusConfig `yaml:"status_config"`

	// Generator settings
	Minify  bool `yaml:"minify"`
	Verbose bool `yaml:"verbose"`
}

// DefaultConfig returns a configuration with sane defaults
func DefaultConfig() *Config {
	return &Config{
		ADRDirectory:    "adr",
		OutputDirectory: "docs",
		BaseURL:         "",
		DefaultCategory: "General",
		AllowedCategories: []string{
			"Core Architecture",
			"Data Management",
			"Frontend Development",
			"Security",
			"Infrastructure",
			"General",
		},
		AllowedStatuses: []string{
			"Proposed",
			"Accepted",
			"Deprecated",
			"Superseded",
		},
		StatusConfig: map[string]StatusConfig{
			"Accepted": {
				Icon:     "✓",
				Color:    "green",
				CSSClass: "bg-green-500",
			},
			"Proposed": {
				Icon:     "●",
				Color:    "yellow",
				CSSClass: "bg-yellow-500",
			},
			"Deprecated": {
				Icon:     "✗",
				Color:    "red",
				CSSClass: "bg-red-500",
			},
			"Superseded": {
				Icon:     "↑",
				Color:    "purple",
				CSSClass: "bg-purple-500",
			},
		},
		Minify:  false,
		Verbose: false,
	}
}

// LoadConfig loads configuration from a YAML file, falling back to defaults
func LoadConfig(configPath string) (*Config, error) {
	config := DefaultConfig()

	// If no config file specified, look for common names
	if configPath == "" {
		candidates := []string{
			"adr-config.yaml",
			"adr-config.yml",
			".adr-config.yaml",
			".adr-config.yml",
		}

		for _, candidate := range candidates {
			if _, err := os.Stat(candidate); err == nil {
				configPath = candidate
				break
			}
		}
	}

	// If still no config file found, return defaults
	if configPath == "" {
		return config, nil
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, nil
	}

	// Read and parse config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	// Parse YAML
	var fileConfig Config
	if err := yaml.Unmarshal(data, &fileConfig); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	// Merge with defaults (file config takes precedence)
	mergedConfig := mergeConfigs(config, &fileConfig)

	// Validate paths
	if err := validateConfig(mergedConfig); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return mergedConfig, nil
}

// mergeConfigs merges file config with defaults, with file config taking precedence
func mergeConfigs(defaults, fileConfig *Config) *Config {
	merged := *defaults // Copy defaults

	// Override with file config values if they're not zero values
	if fileConfig.ADRDirectory != "" {
		merged.ADRDirectory = fileConfig.ADRDirectory
	}
	if fileConfig.OutputDirectory != "" {
		merged.OutputDirectory = fileConfig.OutputDirectory
	}
	if fileConfig.BaseURL != "" {
		merged.BaseURL = fileConfig.BaseURL
	}
	if fileConfig.DefaultCategory != "" {
		merged.DefaultCategory = fileConfig.DefaultCategory
	}
	if len(fileConfig.AllowedCategories) > 0 {
		merged.AllowedCategories = fileConfig.AllowedCategories
	}
	if len(fileConfig.AllowedStatuses) > 0 {
		merged.AllowedStatuses = fileConfig.AllowedStatuses
	}

	// Merge status configs
	if len(fileConfig.StatusConfig) > 0 {
		if merged.StatusConfig == nil {
			merged.StatusConfig = make(map[string]StatusConfig)
		}
		for status, config := range fileConfig.StatusConfig {
			merged.StatusConfig[status] = config
		}
	}

	// Boolean flags
	merged.Minify = fileConfig.Minify
	merged.Verbose = fileConfig.Verbose

	return &merged
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	// Validate ADR directory exists
	if _, err := os.Stat(config.ADRDirectory); os.IsNotExist(err) {
		return fmt.Errorf("ADR directory does not exist: %s", config.ADRDirectory)
	}

	// Ensure output directory parent exists
	outputParent := filepath.Dir(config.OutputDirectory)
	if outputParent != "." {
		if _, err := os.Stat(outputParent); os.IsNotExist(err) {
			return fmt.Errorf("output directory parent does not exist: %s", outputParent)
		}
	}

	// Validate that all allowed statuses have status configs
	for _, status := range config.AllowedStatuses {
		if _, exists := config.StatusConfig[status]; !exists {
			// Add default config for missing status
			config.StatusConfig[status] = StatusConfig{
				Icon:     "?",
				Color:    "gray",
				CSSClass: "bg-gray-500",
			}
		}
	}

	return nil
}

// GetStatusIcon returns the icon for a given status
func (c *Config) GetStatusIcon(status string) string {
	statusConfig, exists := c.StatusConfig[status]
	if !exists {
		return "?"
	}
	return statusConfig.Icon
}

// GetStatusClass returns the CSS class for a given status
func (c *Config) GetStatusClass(status string) string {
	statusConfig, exists := c.StatusConfig[status]
	if !exists {
		return "bg-gray-500"
	}
	return statusConfig.CSSClass
}

// GetStatusColor returns the color for a given status
func (c *Config) GetStatusColor(status string) string {
	statusConfig, exists := c.StatusConfig[status]
	if !exists {
		return "gray"
	}
	return statusConfig.Color
}

// IsValidCategory checks if a category is in the allowed list
func (c *Config) IsValidCategory(category string) bool {
	for _, allowed := range c.AllowedCategories {
		if category == allowed {
			return true
		}
	}
	return false
}

// IsValidStatus checks if a status is in the allowed list
func (c *Config) IsValidStatus(status string) bool {
	for _, allowed := range c.AllowedStatuses {
		if status == allowed {
			return true
		}
	}
	return false
}
