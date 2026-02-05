// Package small provides the implementation for small-scale project generation.
package small

import (
	"github.com/xRiot45/gocrafting/internal/core"
)

// Provider adalah struct yang mengimplementasikan core.FeatureProvider
type Provider struct{}

// NewProvider mengembalikan instance provider baru
func NewProvider() Provider {
	return Provider{}
}

// GetTemplates (Implementasi Interface)
func (p Provider) GetTemplates() []string {
	return GetTemplates()
}

// GetFrameworks (Implementasi Interface)
func (p Provider) GetFrameworks(template string) []string {
	return GetFrameworks(template)
}

// GetDatabaseDrivers (Implementasi Interface)
func (p Provider) GetDatabaseDrivers(template string) []string {
	return GetDatabaseDrivers(template)
}

// Generate (Implementasi Interface)
func (p Provider) Generate(config core.ProjectConfig) error {
	return Generate(config)
}
