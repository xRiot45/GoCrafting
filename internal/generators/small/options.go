// Package small implements the logic and configuration for "Small" scale projects.
package small

// GetTemplates returns available templates for Small scale
func GetTemplates() []string {
	return []string{
		"Simple API",
		"Fast HTTP",
		"CLI Tool",
		"Telegram Bot Starter",
	}
}

// GetFrameworks returns available frameworks based on the selected template
func GetFrameworks(template string) []string {
	switch template {
	case "Fast HTTP":
		return []string{
			"Fiber",
			"Gin",
		}

	case "Simple API":
		return []string{}
	case "CLI Tool":
		return []string{}
	case "Telegram Bot Starter":
		return []string{}
	default:
		return []string{}
	}
}

// GetDatabaseDrivers returns available database options
func GetDatabaseDrivers(template string) []string {
	if template == "CLI Tool" || template == "Telegram Bot Starter" {
		return []string{
			"None",
		}
	}

	return []string{
		"None",
		"SQLite",
		"MySQL",
		"PostgreSQL",
	}
}
