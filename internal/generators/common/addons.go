// Package common contains functions to generate common add-ons for Go projects.
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

	if config.HasAddon("Environment File (.env)") {
		envFiles := map[string]string{
			"common/env/env_development.tmpl": ".env.development", // Dev config
			"common/env/env_example.tmpl":     ".env.example",     // Master Documentation
			"common/env/env_production.tmpl":  ".env.production",  // Prod config
			"common/env/env_staging.tmpl":     ".env.staging",     // Staging config
			"common/env/env_test.tmpl":        ".env.test",        // CI/CD config
		}

		for tpl, output := range envFiles {
			if err := renderAndWrite(config, tpl, output); err != nil {
				return fmt.Errorf("failed to create %s: %w", output, err)
			}
		}

	}

	if config.HasAddon("Gitignore File") {
		if err := renderAndWrite(config, "common/gitignore.tmpl", ".gitignore"); err != nil {
			return fmt.Errorf("failed to create .gitignore: %w", err)
		}
	}

	if config.HasAddon("Readme File") {
		if err := renderAndWrite(config, "common/readme.tmpl", "README.md"); err != nil {
			return fmt.Errorf("failed to create README.md: %w", err)
		}
	}

	if config.HasAddon("Dockerfile") {
		dockerFiles := map[string]string{
			"common/docker/Dockerfile.tmpl":     "Dockerfile",
			"common/docker/.dockerignore.tmpl":  ".dockerignore",
			"common/docker/docker-compose.tmpl": "docker-compose.yaml",
		}

		for tpl, output := range dockerFiles {
			if err := renderAndWrite(config, tpl, output); err != nil {
				return fmt.Errorf("failed to create %s: %w", output, err)
			}
		}
	}

	if config.HasAddon("GitHub Actions (CI/CD)") {
		workflowsDir := filepath.Join(config.ProjectName, ".github", "workflows")
		if err := os.MkdirAll(workflowsDir, 0750); err != nil {
			return fmt.Errorf("failed to create workflows dir: %w", err)
		}

		// 2. Map Template -> Output
		ciFiles := map[string]string{
			"common/github/ci.tmpl":         ".github/workflows/ci.yaml",
			"common/github/release.tmpl":    ".github/workflows/release.yaml",
			"common/github/dependabot.tmpl": ".github/dependabot.yaml",
			"common/github/goreleaser.tmpl": ".goreleaser.yaml",
		}

		for tpl, output := range ciFiles {
			if err := renderAndWrite(config, tpl, output); err != nil {
				return fmt.Errorf("failed to generate %s: %w", output, err)
			}
		}
	}

	if config.HasAddon("Editor Config File") {
		if err := renderAndWrite(config, "common/editorconfig.tmpl", ".editorconfig"); err != nil {
			return fmt.Errorf("failed to create .editorconfig: %w", err)
		}
	}

	if config.HasAddon("Makefile (Shortcut Commands)") {
		if err := renderAndWrite(config, "common/makefile.tmpl", "Makefile"); err != nil {
			return fmt.Errorf("failed to create Makefile: %w", err)
		}
	}

	if config.HasAddon("Lefthook (Commit Linter)") {
		if err := renderAndWrite(config, "common/lefthook.tmpl", "lefthook.yaml"); err != nil {
			return fmt.Errorf("failed to create lefthook.yaml: %w", err)
		}

		fmt.Println("   ⚠️  IMPORTANT: Run 'lefthook install' inside the project to activate hooks.")
	}

	return nil
}

// renderAndWrite renders the template based on the provided configuration and writes the result to the target file.
func renderAndWrite(config core.ProjectConfig, templatePath string, outputPath string) error {
	// 1. Read Template from Embed FS
	tplContent, err := templates.FS.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("template not found: %s", templatePath)
	}

	// 2. Parse Template
	tmpl, err := template.New(outputPath).Parse(string(tplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	// 3. Execute Template (Inject Data)
	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, config); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	// 4. Prepare Output Path
	// Handle jika outputPath mengandung folder (misal: .github/workflows/ci.yml)
	fullPath := filepath.Join(config.ProjectName, outputPath)
	dir := filepath.Dir(fullPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0750); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// 5. Write File
	if err := os.WriteFile(fullPath, buffer.Bytes(), 0600); err != nil {
		return fmt.Errorf("failed to write file %s: %w", fullPath, err)
	}

	return nil
}
