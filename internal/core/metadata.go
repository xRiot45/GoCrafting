package core

import (
	"encoding/json"
	"os"
	"time"
)

// ProjectMetadata menyimpan state konfigurasi project.
type ProjectMetadata struct {
	CLIVersion             string    `json:"cli_version"`
	ProjectName            string    `json:"project_name"`
	ModuleName             string    `json:"module_name"`
	ProjectScale           string    `json:"project_scale"`
	SelectedTemplate       string    `json:"selected_template"`
	SelectedFramework      string    `json:"selected_framework"`
	SelectedDatabaseDriver string    `json:"selected_database_driver"`
	SelectedAddons         []string  `json:"selected_addons"`
	CreatedAt              time.Time `json:"created_at"`
}

// SaveMetadata menulis file gocrafting-cli.json ke root project
func SaveMetadata(path string, meta ProjectMetadata) error {
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path+"/gocrafting-cli.json", data, 0600)
}

// LoadMetadata membaca file gocrafting-cli.json (dipakai command generate)
func LoadMetadata() (*ProjectMetadata, error) {
	file, err := os.ReadFile("gocrafting-cli.json")
	if err != nil {
		return nil, err
	}

	var meta ProjectMetadata
	if err := json.Unmarshal(file, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}
