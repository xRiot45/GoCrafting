package generator

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/xRiot45/gocrafting/internal/templates"
)

// Forge is the main entry point to initialize a new project.
func Forge(config ProjectConfig) error {
	// 1. Create root project folder with secure permissions (0750)
	if err := os.MkdirAll(config.ProjectName, 0750); err != nil {
		return fmt.Errorf("failed to create project folder: %w", err)
	}

	// 2. Execute resource copying and template processing
	return copyResources(config.ProjectScale, config)
}

// copyResources scans internal/templates and copies them to the destination.
func copyResources(sourceDir string, config ProjectConfig) error {
	return fs.WalkDir(templates.FS, sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(sourceDir, path)
		if relPath == "." {
			return nil
		}

		targetPath := filepath.Join(config.ProjectName, strings.TrimSuffix(relPath, ".tmpl"))

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0750)
		}

		return forgeFile(path, targetPath, config)
	})
}

// forgeFile reads, processes, and writes a single file to disk with secure permissions (0600).
func forgeFile(sourcePath, targetPath string, config ProjectConfig) error {
	content, err := fs.ReadFile(templates.FS, sourcePath)
	if err != nil {
		return err
	}

	processedContent := content
	if strings.HasSuffix(sourcePath, ".tmpl") {
		processedContent, err = processTemplate(sourcePath, content, config)
		if err != nil {
			return err
		}
	}

	return os.WriteFile(targetPath, processedContent, 0600)
}
