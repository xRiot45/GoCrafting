// Package small implements the logic and configuration for "Small" scale projects.
package small

// GetTemplates returns available templates for Small scale
func GetTemplates() []string {
	return []string{
		"simple-api",
		"fast-http",
		"cli-tool",
		"telegram-bot-starter",
	}
}

// GetFrameworks returns available frameworks based on the selected template
func GetFrameworks(template string) []string {
	switch template {
	case "fast-http":
		return []string{
			"Fiber",
			"Gin",
		}

	case "simple-api":
		return []string{}
	case "cli-tool":
		return []string{}
	case "bot-starter":
		return []string{}
	default:
		return []string{}
	}
}

// GetDatabaseDrivers returns available database options
func GetDatabaseDrivers(template string) []string {
	if template == "cli-tool" || template == "telegram-bot-starter" {
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
