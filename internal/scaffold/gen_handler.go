// Package scaffold is a traffic controller for generating boilerplate code.
package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xRiot45/gocrafting/internal/core"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// GenerateHandler menentukan strategi template berdasarkan metadata project.
func GenerateHandler(meta *core.ProjectMetadata, name string) error {
	lowerName := strings.ToLower(name)
	titleName := cases.Title(language.English).String(lowerName)

	var templateFile string

	if meta.SelectedTemplate == "Simple API" {
		templateFile = "handlers/net_http.tmpl"
	} else {
		switch meta.SelectedFramework {
		case "Fiber":
			templateFile = "handlers/fiber.tmpl"
		case "Gin":
			templateFile = "handlers/gin.tmpl"
		default:
			return fmt.Errorf("framework %s not supported for handler generation", meta.SelectedFramework)
		}
	}

	targetDir := "internal/handlers"
	if meta.ProjectScale == "Small" {
		targetDir = "handlers"
	}

	fileName := fmt.Sprintf("%s_handler.go", lowerName)
	targetPath := filepath.Clean(filepath.Join(targetDir, fileName))

	data := TemplateData{
		PackageName: "handlers",
		StructName:  titleName,
		ModuleName:  meta.ModuleName,
	}

	if err := renderFile(templateFile, targetPath, data); err != nil {
		return err
	}

	fmt.Printf("   Created handler: %s\n", targetPath)
	return nil
}
