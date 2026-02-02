package generator

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xRiot45/gocrafting/internal/runner"
	"github.com/xRiot45/gocrafting/internal/templates"
)

// Forge is the main entry point to initialize a new project.
func Forge(config ProjectConfig) error {
	// 1. Create root project
	if err := os.MkdirAll(config.ProjectName, 0750); err != nil {
		return fmt.Errorf("failed to create project folder: %w", err)
	}

	// 2. Generating project metadata file
	if err := createMetaFile(config); err != nil {
		return fmt.Errorf("failed to create metadata file: %w", err)
	}

	// 3. Copy project template
	templatePath := filepath.Join("small", config.SelectedTemplate)
	if err := copyResources(templatePath, config); err != nil {
		return err
	}

	if err := finalizeProject(config.ProjectName, config); err != nil {
		return fmt.Errorf("failed to finalize project: %w", err)
	}

	return nil
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

// createMetaFile creates the project metadata file.
func createMetaFile(config ProjectConfig) error {
	meta := ProjectMeta{
		CLIVersion:  "v1.0.0",
		Name:        config.ProjectName,
		Module:      config.ModuleName,
		Scale:       config.ProjectScale,
		Template:    config.SelectedTemplate,
		Persistence: config.Persistence,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}

	fileContent, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}

	targetPath := filepath.Join(config.ProjectName, "gocrafting-cli.json")
	return os.WriteFile(targetPath, fileContent, 0600)
}

// finalizeProject performs final tasks like running 'go mod tidy' and 'go fmt'
func finalizeProject(projectPath string, config ProjectConfig) error {
	// 1. Identify packages to install
	var packagesToInstall []string

	// Check Persistence
	if config.Persistence == "sqlite" {
		packagesToInstall = append(packagesToInstall, "modernc.org/sqlite")
	}

	// 2. Run 'go get' if there is a package that must be installed
	if len(packagesToInstall) > 0 {
		if err := runner.GoGet(projectPath, packagesToInstall...); err != nil {
			return err
		}
	}

	// 3. Run 'go mod tidy' (Final Cleaning)
	if err := runner.RunGoModTidy(projectPath); err != nil {
		return err
	}

	// 4. Run 'go fmt' (Formatting)
	if err := runner.RunGoFmt(projectPath); err != nil {
		fmt.Printf("Warning: %v\n", err)
	}

	return nil
}
