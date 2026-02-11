package ui

import "github.com/xRiot45/gocrafting/internal/core"

// reconstructConfig menyusun ulang ProjectConfig dari data UI saat ini.
// Digunakan saat memulai instalasi dan saat install dependency.
func (uiModel MainModel) reconstructConfig() core.ProjectConfig {
	var selectedAddons []string

	// Konversi Map Indices ke Slice String
	for i, label := range AddonList {
		if uiModel.SelectedAddonsIndices[i] {
			selectedAddons = append(selectedAddons, label)
		}
	}

	return core.ProjectConfig{
		ProjectName:            uiModel.ProjectName,
		ModuleName:             uiModel.ModuleName,
		ProjectScale:           uiModel.ProjectScale,
		SelectedTemplate:       uiModel.SelectedTemplate,
		SelectedFramework:      uiModel.SelectedFramework,
		SelectedDatabaseDriver: uiModel.SelectedDatabaseDriver,
		SelectedAddons:         selectedAddons,
	}
}

// isDisabledTemplate mengecek apakah templateName sedang disabled.
func isDisabledTemplate(templateName string) bool {
	disabledTemplates := map[string]bool{
		// Small Project
		"CLI Tool":             true,
		"Telegram Bot Starter": true,
	}

	return disabledTemplates[templateName]
}

// isDisabledProjectScale mengecek apakah scaleName sedang disabled.
func isDisabledProjectScale(scale string) bool {
	disabledScales := map[string]bool{
		"Medium":     true,
		"Enterprise": true,
	}

	return disabledScales[scale]
}
