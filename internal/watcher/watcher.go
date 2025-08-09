package watcher

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Config holds the file watcher configuration
type Config struct {
	Paths   []string
	Verbose bool
}

// Event represents a file change event
type Event struct {
	Path string
	Type string // "create", "modify", "delete"
}

// FileWatcher watches for file changes
type FileWatcher struct {
	config   *Config
	watcher  *fsnotify.Watcher
	Events   chan Event
	debounce map[string]time.Time
}

// New creates a new file watcher
func New(config *Config) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}

	fw := &FileWatcher{
		config:   config,
		watcher:  watcher,
		Events:   make(chan Event, 100),
		debounce: make(map[string]time.Time),
	}

	// Add paths to watcher
	for _, path := range config.Paths {
		if err := fw.addPath(path); err != nil {
			return nil, fmt.Errorf("failed to watch path %s: %w", path, err)
		}
	}

	// Start watching
	go fw.watch()

	return fw, nil
}

// addPath adds a path to the watcher recursively
func (fw *FileWatcher) addPath(path string) error {
	return filepath.Walk(path, func(walkPath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		if !info.IsDir() {
			return nil
		}

		if fw.config.Verbose {
			fmt.Printf("👀 Watching directory: %s\n", walkPath)
		}
		return fw.watcher.Add(walkPath)
	})
}

// watch processes file system events
func (fw *FileWatcher) watch() {
	defer fw.watcher.Close()

	for {
		select {
		case event, ok := <-fw.watcher.Events:
			if !ok {
				return
			}

			// Filter relevant files
			if !fw.isRelevantFile(event.Name) {
				continue
			}

			// Debounce events (avoid multiple events for the same file)
			now := time.Now()
			if lastTime, exists := fw.debounce[event.Name]; exists {
				if now.Sub(lastTime) < 100*time.Millisecond {
					continue
				}
			}
			fw.debounce[event.Name] = now

			// Determine event type
			eventType := "modify"
			if event.Op&fsnotify.Create == fsnotify.Create {
				eventType = "create"
			} else if event.Op&fsnotify.Remove == fsnotify.Remove {
				eventType = "delete"
			}

			// Send event
			select {
			case fw.Events <- Event{Path: event.Name, Type: eventType}:
			default:
				// Channel full, skip event
			}

		case err, ok := <-fw.watcher.Errors:
			if !ok {
				return
			}
			if fw.config.Verbose {
				fmt.Printf("⚠️  Watcher error: %v\n", err)
			}
		}
	}
}

// isRelevantFile checks if a file should trigger a rebuild
func (fw *FileWatcher) isRelevantFile(path string) bool {
	// Skip hidden files and directories
	if strings.Contains(path, "/.") {
		return false
	}

	// Skip backup files
	if strings.HasSuffix(path, "~") || strings.HasSuffix(path, ".bak") {
		return false
	}

	// Check for relevant extensions
	ext := strings.ToLower(filepath.Ext(path))
	relevantExts := []string{
		".md", ".html", ".css", ".js", ".json", ".yaml", ".yml",
	}

	for _, relevantExt := range relevantExts {
		if ext == relevantExt {
			return true
		}
	}

	return false
}

// Close closes the file watcher
func (fw *FileWatcher) Close() error {
	close(fw.Events)
	return fw.watcher.Close()
}
