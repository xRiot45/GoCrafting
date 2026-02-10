package core

// AddonOption represents a single add-on option with its metadata.
type AddonOption struct {
	ID          string
	Label       string
	Description string
}

// AvailableAddons is the list of all supported add-ons.
var AvailableAddons = []AddonOption{
	{
		ID:    "env",
		Label: "Environment File (.env)",
	},
	{
		ID:    "gitignore",
		Label: "Gitignore File",
	},
	{
		ID:    "readme",
		Label: "Readme File (Markdown)",
	},
	{
		ID:    "editorconfig",
		Label: "EditorConfig (.editorconfig)",
	},
	{
		ID:    "makefile",
		Label: "Makefile (Shortcut Commands)",
	},
	{
		ID:    "docker",
		Label: "Docker & Compose Support",
	},
	{
		ID:    "github_action",
		Label: "GitHub Actions (CI/CD Pipelines)",
	},
	{
		ID:    "lefthook",
		Label: "Lefthook (Git Hooks/Commit Linter)",
	},
}

// GetAddonLabelByID returns the label of an add-on given its ID.
func GetAddonLabelByID(id string) string {
	for _, addon := range AvailableAddons {
		if addon.ID == id {
			return addon.Label
		}
	}
	return "Unknown Addon"
}
