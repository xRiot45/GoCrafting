// Package templates handles the embedded static files for project generation.
package templates

import (
	"embed"
)

// ProjectTemplates holds the embedded filesystem containing all project templates (small, medium, etc).
//
//go:embed all:*
var ProjectTemplates embed.FS
