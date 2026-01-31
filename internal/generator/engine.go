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
	if err := os.MkdirAll(config.ProjectName, 0750); err != nil {
		return fmt.Errorf("failed to create project folder: %w", err)
	}

	templatePath := filepath.Join("small", config.SelectedTemplate)
	return copyResources(templatePath, config)
}

// copyResources scans internal/templates and copies them to the destination.
func copyResources(sourceDir string, config ProjectConfig) error {
	return fs.WalkDir(templates.FS, sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

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

// forgeFile reads, processes, and writes a single file to disk.
func forgeFile(sourcePath, targetPath string, config ProjectConfig) error {
	content, err := fs.ReadFile(templates.FS, sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %w", sourcePath, err)
	}

	var processedContent []byte
	if strings.HasSuffix(sourcePath, ".tmpl") {
		processedContent, err = processTemplate(sourcePath, content, config)
		if err != nil {
			return err
		}
	} else {
		processedContent = content
	}

	return os.WriteFile(targetPath, processedContent, 0600)
}
