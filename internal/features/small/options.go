// Package small implements the logic and configuration for "Small" scale projects.
package small

// GetTemplates returns available templates for Small scale
func GetTemplates() []string {
	return []string{
		"simple-api",
		"fast-http",
		"cli-tool",
		"bot-starter",
	}
}

// GetPersistence returns available database options
func GetPersistence() []string {
	return []string{
		"None",
		"SQLite",
	}
}
