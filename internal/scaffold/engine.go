package scaffold

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed all:*
var templatesFS embed.FS

// TemplateData adalah struktur data universal yang dikirim ke semua template.
type TemplateData struct {
	PackageName string
	StructName  string
	ModuleName  string
}

// renderFile adalah fungsi generic untuk menulis file dari template embed.
func renderFile(templatePath string, targetPath string, data TemplateData) error {
	fullPath := "templates/" + templatePath

	tplContent, err := templatesFS.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("template not found in embed fs: %s", fullPath)
	}

	tmpl, err := template.New(templatePath).Parse(string(tplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	cleanPath := filepath.Clean(targetPath)

	dir := filepath.Dir(cleanPath)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	f, err := os.Create(cleanPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", cleanPath, err)
	}

	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close file %s: %v\n", cleanPath, closeErr)
		}
	}()

	// 7. Eksekusi Template ke File
	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
