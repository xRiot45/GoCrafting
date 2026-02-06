package common

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/xRiot45/gocrafting/internal/core"
	"github.com/xRiot45/gocrafting/internal/templates"
)

// GenerateAddons generates the selected add-ons based on the provided configuration.
func GenerateAddons(config core.ProjectConfig) error {
	fmt.Println("📦 Generating selected add-ons...")

	// 1. .env (Environment Variables)
	if config.HasAddon("Environment File (.env)") {
		// Create .env file
		if err := renderAndWrite(config, "shared/env.tmpl", ".env"); err != nil {
			return err
		}
		// Create .env.example
		if err := renderAndWrite(config, "shared/env.tmpl", ".env.example"); err != nil {
			return err
		}
	}

	return nil
}

// renderAndWrite renders the template based on the provided configuration and writes the result to the target file.
func renderAndWrite(config core.ProjectConfig, templateName string, targetFileName string) error {
	tmplContent, err := templates.FS.ReadFile(templateName)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templateName, err)
	}

	tmpl, err := template.New(targetFileName).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, config); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	targetPath := filepath.Join(config.ProjectName, targetFileName)
	if err := os.WriteFile(targetPath, buffer.Bytes(), 0600); err != nil {
		return fmt.Errorf("failed to write file %s: %w", targetPath, err)
	}

	return nil
}
