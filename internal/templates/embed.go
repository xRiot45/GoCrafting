// Package templates handles the embedded filesystem for project blueprints.
package templates

import "embed"

// FS holds all template files within this directory.
//
//go:embed all:*
var FS embed.FS
