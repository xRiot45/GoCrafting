// Package ui contains styles and themes for the terminal UI components.
package ui

import "github.com/charmbracelet/lipgloss"

// --- PALETTE (Professional DevOps Theme) ---
var (
	ColorFocus   = lipgloss.Color("#3B82F6") // Bright Blue (Active)
	ColorBlur    = lipgloss.Color("#4B5563") // Gray (Inactive)
	ColorSuccess = lipgloss.Color("#10B981") // Emerald Green (Done)
	ColorText    = lipgloss.Color("#F3F4F6") // Whiteish
	ColorBg      = lipgloss.Color("#111827") // Dark Background (Optional if terminal is dark)
)

// --- LAYOUT STYLES ---
var (
	// Sidebar (Kiri): Lebar tetap, border kanan
	SidebarStyle = lipgloss.NewStyle().
			Width(25).
			PaddingRight(2).
			Border(lipgloss.NormalBorder(), false, true, false, false). // Border kanan saja
			BorderForeground(ColorBlur)

	// Main Content (Kanan): Padding kiri
	MainContentStyle = lipgloss.NewStyle().
				PaddingLeft(3).
				PaddingTop(1)

	// Step (Sidebar Items)
	StepPendingStyle = lipgloss.NewStyle().Foreground(ColorBlur)
	StepActiveStyle  = lipgloss.NewStyle().Foreground(ColorFocus).Bold(true)
	StepDoneStyle    = lipgloss.NewStyle().Foreground(ColorSuccess)

	// Header di Main Window
	HeaderStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			Background(ColorFocus).
			Bold(true).
			Padding(0, 1).
			MarginBottom(1)

	// Helper Text
	DescStyle = lipgloss.NewStyle().Foreground(ColorBlur).Italic(true)
)
