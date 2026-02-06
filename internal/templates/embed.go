// Package templates handles the embedded static files for project generation.
package templates

import (
	"embed"
)

// FS contains all the embedded templates for different project scales and features.
//
//go:embed all:*
var FS embed.FS
