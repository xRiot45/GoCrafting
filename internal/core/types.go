// Package core contains shared data structures and types used across the application.
package core

// ProjectConfig stores the data collected from the user via the TUI.
type ProjectConfig struct {
	ProjectName            string
	ModuleName             string
	ProjectScale           string
	SelectedTemplate       string
	SelectedFramework      string
	SelectedDatabaseDriver string
	SelectedAddons         []string
}

// ProjectMeta stores project metadata.
type ProjectMeta struct {
	CLIVersion     string   `json:"cli_version"`
	Name           string   `json:"project_name"`
	Module         string   `json:"module_name"`
	Scale          string   `json:"project_scale"`
	Template       string   `json:"selected_template"`
	Framework      string   `json:"selected_framework"`
	DatabaseDriver string   `json:"selected_database_driver"`
	Addons         []string `json:"selected_addons"`
	CreatedAt      string   `json:"created_at"`
}

// HasAddon checks if the given addonName is present in the SelectedAddons slice.
// It returns true if the addonName is found, and false otherwise.
func (c ProjectConfig) HasAddon(addonName string) bool {
	for _, a := range c.SelectedAddons {
		if a == addonName {
			return true
		}
	}
	return false
}
