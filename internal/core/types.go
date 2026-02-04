// Package core contains shared data structures and types used across the application.
package core

// ProjectConfig stores the data collected from the user via the TUI.
type ProjectConfig struct {
	ProjectName      string
	ModuleName       string
	ProjectScale     string
	SelectedTemplate string
	Persistence      string
}

// ProjectMeta stores project metadata.
type ProjectMeta struct {
	CLIVersion  string `json:"cli_version"`
	Name        string `json:"project_name"`
	Module      string `json:"module_name"`
	Scale       string `json:"project_scale"`
	Template    string `json:"selected_template"`
	Persistence string `json:"persistence"`
	CreatedAt   string `json:"created_at"`
}
