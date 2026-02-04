// Package core defines the shared interfaces and data structures for the application.
package core

// FeatureProvider adalah kontrak yang harus dipenuhi oleh setiap fitur scale (Small, Medium, Enterprise).
// UI akan menggunakan interface ini untuk mengambil data dinamis.
type FeatureProvider interface {
	// Mengembalikan daftar template yang tersedia untuk scale ini
	GetTemplates() []string

	// Mengembalikan daftar opsi database/persistence untuk scale ini
	GetPersistenceOptions() []string

	// Menjalankan logika generate project
	Generate(config ProjectConfig) error
}
