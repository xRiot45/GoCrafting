// Package generator provides the engine and configuration for project creation.
package generator

// ProjectConfig stores the data collected from the user via the TUI.
type ProjectConfig struct {
	ProjectName      string
	ModuleName       string
	ProjectScale     string
	SelectedTemplate string
	Persistence      string
}
