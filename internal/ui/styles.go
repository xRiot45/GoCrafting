// Package ui provides UI constants, styling, and terminal interface logic
// for the GoCrafting CLI using the Lip Gloss library.
package ui

import "github.com/charmbracelet/lipgloss"

// logoASCII contains the ASCII art representation of the GoCrafting logo.
const logoASCII = `
  ____             ____            __ _   _                 
 / ___|  ___      / ___|_ __ __ _ / _| |_(_)_ __   __ _ 
| |  _  / _ \    | |   | '__/ _' | |_| __| | '_ \ / _' |
| |_| | (_) |    | |___| | | (_| |  _| |_| | | | | (_| |
 \____| \___/     \____|_|  \__,_|_|  \__|_|_| |_|\__, |
                                                  |___/ 
`

var (
	// ColorPrimary is the royal blue color used for main branding.
	ColorPrimary = lipgloss.Color("#566ef4")

	// ColorSecondary is a snow white color for primary text readability.
	ColorSecondary = lipgloss.Color("#FAFAFA")

	// ColorAccent is a neon cyan color for highlighting inputs and active states.
	ColorAccent = lipgloss.Color("#00f2ff")

	// ColorGold is a gold color for success messages or artisan-level elements.
	ColorGold = lipgloss.Color("#FFD700")

	// ColorMuted is a dimmed gray color for secondary hints and help text.
	ColorMuted = lipgloss.Color("#626262")
)

var (
	// LogoStyle defines the visual style for the main ASCII logo.
	LogoStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	// TitleStyle defines the visual style for questions or section headers.
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Bold(true)

	// InputStyle defines the visual style for user-provided text inputs.
	InputStyle = lipgloss.NewStyle().
			Foreground(ColorAccent)

	// HintStyle defines the visual style for small help hints at the bottom.
	HintStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Italic(true)
)
