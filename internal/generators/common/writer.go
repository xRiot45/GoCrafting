package common

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/xRiot45/gocrafting/internal/core"
	"github.com/xRiot45/gocrafting/internal/templates"
)

// BaseGenerate is the foundation function used by all features (Small/Medium/Enterprise).
// Task: Create Project Folder -> Create JSON Meta File -> Copy Template.
func BaseGenerate(config core.ProjectConfig, templateSourcePath string) error {
	// 1. Create root project folder
	if err := os.MkdirAll(config.ProjectName, 0750); err != nil {
		return fmt.Errorf("failed to create project folder: %w", err)
	}

	// 2. Generating project metadata file (gocrafting-cli.json)
	if err := createMetaFile(config); err != nil {
		return fmt.Errorf("failed to create metadata file: %w", err)
	}

	// 3. Copy project template files
	// Note: templateSourcePath is sent from the caller (eg: "small/simple-api")
	if err := copyResources(templateSourcePath, config); err != nil {
		return err
	}

	return nil
}

// copyResources scans internal/templates and copies them to the destination.
func copyResources(sourceDir string, config core.ProjectConfig) error {
	fileSystem := templates.FS

	return fs.WalkDir(fileSystem, sourceDir, func(path string, d fs.DirEntry, err error) error {
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

		return forgeFile(fileSystem, path, targetPath, config)
	})
}

// forgeFile reads, processes, and writes a single file to disk.
func forgeFile(fileSystem fs.FS, sourcePath, targetPath string, config core.ProjectConfig) error {
	content, err := fs.ReadFile(fileSystem, sourcePath)
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

// createMetaFile creates the project metadata file using struct from core.
func createMetaFile(config core.ProjectConfig) error {

	// frameworkName := config.SelectedFramework
	// if frameworkName == "" {
	// 	frameworkName = "None"
	// }

	// 2. ISI STRUCT META
	meta := core.ProjectMetadata{
		CLIVersion:             "v1.0.0",
		ProjectName:            config.ProjectName,
		ModuleName:             config.ModuleName,
		ProjectScale:           config.ProjectScale,
		SelectedTemplate:       config.SelectedTemplate,
		SelectedFramework:      config.SelectedFramework,
		SelectedDatabaseDriver: config.SelectedDatabaseDriver,
		SelectedAddons:         config.SelectedAddons,
		CreatedAt:              time.Now(),
	}

	// 3. WRITE FILE
	if err := core.SaveMetadata(config.ProjectName, meta); err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	return nil
}

// processTemplate processes the raw template data using the Go text/template engine.
func processTemplate(name string, content []byte, config core.ProjectConfig) ([]byte, error) {
	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config); err != nil {
		return nil, fmt.Errorf("failed to execute template %s: %w", name, err)
	}

	return buf.Bytes(), nil
}
