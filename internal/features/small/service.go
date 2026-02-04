package small

import (
	"path/filepath"

	"github.com/xRiot45/gocrafting/internal/core"
	"github.com/xRiot45/gocrafting/internal/features/common"
	"github.com/xRiot45/gocrafting/internal/runner"
)

// Generate adalah Entry Point khusus Small Project
func Generate(config core.ProjectConfig) error {
	// 1. Set Default
	if config.Persistence == "" {
		config.Persistence = "none"
	}

	// 2. Panggil Base Generator (Folder, JSON, Template)
	templatePath := filepath.Join("small", config.SelectedTemplate)
	if err := common.BaseGenerate(config, templatePath); err != nil {
		return err
	}

	// 3. Install Dependencies (Khusus Small)
	return installDependencies(config)
}

func installDependencies(config core.ProjectConfig) error {
	var packages []string

	// Logic Library: SQLite ModernC
	if config.Persistence == "SQLite" {
		packages = append(packages, "modernc.org/sqlite")
	}

	// Jika ada template FastHTTP (Fiber), tambah package disini
	if config.SelectedTemplate == "fast-http" {
		packages = append(packages, "github.com/gofiber/fiber/v2")
	}

	// Eksekusi Runner
	if len(packages) > 0 {
		if err := runner.GoGet(config.ProjectName, packages...); err != nil {
			return err
		}
	}

	if err := runner.RunGoModTidy(config.ProjectName); err != nil {
		return err
	}

	return runner.RunGoFmt(config.ProjectName)
}
