// Package scaffold implements the scaffolding logic for generating boilerplate code.
package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/xRiot45/gocrafting/internal/core"
)

// Template String Sederhana (Bisa dipindah ke file .tmpl nanti)
const fiberHandlerTmpl = `package {{.PackageName}}

import "github.com/gofiber/fiber/v2"

type {{.StructName}}Handler struct {}

func New{{.StructName}}Handler() *{{.StructName}}Handler {
	return &{{.StructName}}Handler{}
}

func (h *{{.StructName}}Handler) Index(c *fiber.Ctx) error {
	return c.SendString("Hello form {{.StructName}}")
}
`

const ginHandlerTmpl = `package {{.PackageName}}

import "github.com/gin-gonic/gin"

type {{.StructName}}Handler struct {}

func New{{.StructName}}Handler() *{{.StructName}}Handler {
	return &{{.StructName}}Handler{}
}

func (h *{{.StructName}}Handler) Index(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Hello form {{.StructName}}"})
}
`

// GenerateResource membuat file handler baru berdasarkan framework yang dipilih.
func GenerateResource(meta *core.ProjectMetadata, name string) error {
	// 1. Normalisasi Nama (Alternative to strings.Title)
	lowerName := strings.ToLower(name)

	// Manual Title Case implementation to avoid Deprecated strings.Title and external deps
	runes := []rune(lowerName)
	if len(runes) > 0 {
		runes[0] = unicode.ToUpper(runes[0])
	}
	titleName := string(runes)

	// Tentukan Lokasi Folder
	targetDir := filepath.Join("internal", "handlers")
	if meta.ProjectScale == "Small" {
		targetDir = "handlers"
	}

	// FIX G301: Expect directory permissions to be 0750 or less
	if err := os.MkdirAll(targetDir, 0750); err != nil {
		return err
	}

	// Pilih Template
	var tmplString string
	if meta.SelectedFramework == "Fiber" {
		tmplString = fiberHandlerTmpl
	} else if meta.SelectedFramework == "Gin" {
		tmplString = ginHandlerTmpl
	} else {
		return fmt.Errorf("framework %s not supported yet", meta.SelectedFramework)
	}

	// Render Template
	data := map[string]string{
		"PackageName": "handlers",
		"StructName":  titleName,
	}

	// FIX G304: Clean path to prevent directory traversal
	filePath := filepath.Clean(filepath.Join(targetDir, lowerName+"_handler.go"))

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	// FIX errcheck: Check return value of f.Close
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close file %s: %v\n", filePath, closeErr)
		}
	}()

	tmpl, err := template.New("handler").Parse(tmplString)
	if err != nil {
		return err
	}

	return tmpl.Execute(f, data)
}
