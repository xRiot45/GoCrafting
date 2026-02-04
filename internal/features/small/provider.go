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
	// Memanggil fungsi dari options.go yang sudah ada
	return GetTemplates()
}

// GetPersistenceOptions (Implementasi Interface)
func (p Provider) GetPersistenceOptions() []string {
	// Memanggil fungsi dari options.go yang sudah ada
	return GetPersistence()
}

// Generate (Implementasi Interface)
func (p Provider) Generate(config core.ProjectConfig) error {
	// Memanggil fungsi dari service.go yang sudah ada
	return Generate(config)
}
