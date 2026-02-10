package small

import (
	"fmt"

	"github.com/xRiot45/gocrafting/internal/core"
	common "github.com/xRiot45/gocrafting/internal/generators/common"
	"github.com/xRiot45/gocrafting/internal/shell"
)

// Generate generates a small-scale project based on the provided configuration.
// It will copy the template files to the project directory and install the required dependencies.
// If the SelectedDatabaseDriver field is empty, it will be set to "none".
//
// After generating the project files, it will call installDependencies to install the required dependencies.
//
// Returns an error if there is an issue during the generation process.
func Generate(config core.ProjectConfig) error {
	// 1. Base Structure (Folder, JSON, Template Files)
	// (Pastikan path templateSourcePath benar sesuai folder Anda)
	templatePath := "small/" + config.SelectedTemplate
	if err := common.BaseGenerate(config, templatePath); err != nil {
		return err
	}

	// 2. Install Dependencies
	if err := installDependencies(config); err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}

	// 3. [BARU] GENERATE ADDONS (.env, dll)
	// Panggil fungsi global yang baru kita buat
	if err := common.GenerateAddons(config); err != nil {
		return fmt.Errorf("failed to generate addons: %w", err)
	}

	return nil
}

// installDependencies installs the required dependencies for the generated project based on the provided configuration.
//
// It will get the required packages for the SelectedFramework and SelectedDatabaseDriver fields.
// If the SelectedDatabaseDriver field is empty or "none", it will not include the database driver packages.
//
// After getting the required packages, it will call GoGet to install them.
// Finally, it will call GoModTidy and GoFmt to clean up the project directory.
//
// Returns an error if there is an issue during the installation process.
func installDependencies(config core.ProjectConfig) error {
	var packages []string

	// ---------------------------------------------------------
	// 1. LOGIC BERDASARKAN FRAMEWORK / TEMPLATE
	// ---------------------------------------------------------

	if config.SelectedFramework != "" {
		packages = append(packages, core.GetPackages(config.SelectedFramework)...)
	} else {
		switch config.SelectedTemplate {
		case "cli-tool":
			packages = append(packages, core.GetPackages("Cobra")...)

		case "bot-starter":
			packages = append(packages, core.GetPackages("TelegramBot")...)

		case "simple-api":
			// No additional packages for simple-api without framework
		}
	}

	// ---------------------------------------------------------
	// 2. LOGIC BERDASARKAN PERSISTENCE (Database)
	// ---------------------------------------------------------

	if config.SelectedDatabaseDriver != "" && config.SelectedDatabaseDriver != "None" {
		packages = append(packages, core.GetPackages(config.SelectedDatabaseDriver)...)
	}

	if len(packages) > 0 {
		if err := shell.GoGet(config.ProjectName, packages...); err != nil {
			return err
		}
	}

	if err := shell.RunGoModTidy(config.ProjectName); err != nil {
		return err
	}

	return shell.RunGoFmt(config.ProjectName)
}
