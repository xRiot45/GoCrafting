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

// GenerateHandler generates a handler file based on project scale and framework.
func GenerateHandler(meta *core.ProjectMetadata, name string) error {
	lowerName := strings.ToLower(name)
	titleName := cases.Title(language.English).String(lowerName)

	scaleFolder := strings.ToLower(meta.ProjectScale)

	if scaleFolder != "small" && scaleFolder != "medium" && scaleFolder != "enterprise" {
		return fmt.Errorf("project scale '%s' is not supported", meta.ProjectScale)
	}

	var templateFilename string
	if meta.SelectedTemplate == "Simple API" {
		templateFilename = "net_http.tmpl"
	} else {
		switch meta.SelectedFramework {
		case "Fiber":
			templateFilename = "fiber.tmpl"
		case "Gin":
			templateFilename = "gin.tmpl"
		default:
			return fmt.Errorf("framework '%s' not supported for handler generation", meta.SelectedFramework)
		}
	}

	templatePath := filepath.Join(scaleFolder, "handlers", templateFilename)

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

	if err := renderFile(templatePath, targetPath, data); err != nil {
		return err
	}

	fmt.Printf("   Created handler: %s\n", targetPath)
	return nil
}
