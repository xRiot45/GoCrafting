package core

import (
	"encoding/json"
	"os"
)

// ProjectMeta merepresentasikan isi file gocrafting-cli.json
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

// LoadProjectMeta membaca file JSON dari root folder
func LoadProjectMeta() (*ProjectMeta, error) {
	file, err := os.ReadFile("gocrafting-cli.json")
	if err != nil {
		return nil, err
	}

	var meta ProjectMeta
	if err := json.Unmarshal(file, &meta); err != nil {
		return nil, err
	}

	return &meta, nil
}
