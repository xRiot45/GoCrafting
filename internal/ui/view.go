// Package ui contains the terminal UI views and rendering logic.
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/xRiot45/gocrafting/internal/core"
	"github.com/xRiot45/gocrafting/internal/generators"
)

// View renders the entire UI based on the current state of the MainModel.
func (m MainModel) View() string {
	if m.Err != nil {
		return renderError(m.Err)
	}
	if m.IsQuitting {
		return "\n  👋 Aborting process.\n"
	}

	// 1. Render Sidebar (Progress Tracker)
	sidebar := renderSidebar(m)

	// 2. Render Main Window (Active Step)
	mainContent := renderMainContent(m)

	// 3. Gabungkan Keduanya (Layout Horizontal)
	// Output: [ Sidebar | Main Content ]
	ui := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, mainContent)

	// Tambahkan margin luar agar tidak nempel pinggir terminal
	return lipgloss.NewStyle().Margin(1, 2).Render(ui)
}

// --- SIDEBAR COMPONENT ---
func renderSidebar(m MainModel) string {
	var s strings.Builder

	// Logo Kecil di atas Sidebar
	s.WriteString(lipgloss.NewStyle().Foreground(ColorFocus).Bold(true).Render("GO CRAFTING") + "\n\n")

	// Daftar Step
	steps := []struct {
		state SessionState
		label string
	}{
		{StateInputProjectName, "Identity"},
		{StateInputModuleName, "Go Module"},
		{StateSelectProjectScale, "Architecture"},
		{StateSelectTemplate, "Template"},
		{StateSelectFramework, "Framework"},
		{StateSelectDatabaseDriver, "Database"},
		{StateSelectAddons, "Add-ons"},
		{StateInstalling, "Installation"},
	}

	for _, step := range steps {
		var indicator, label string

		if m.CurrentState == step.state {
			// Step Aktif
			indicator = "●"
			label = StepActiveStyle.Render(step.label)
			indicator = StepActiveStyle.Render(indicator)
		} else if m.CurrentState > step.state {
			// Step Selesai
			indicator = "✔"
			label = StepDoneStyle.Render(step.label)
			indicator = StepDoneStyle.Render(indicator)
		} else {
			// Step Belum
			indicator = "○"
			label = StepPendingStyle.Render(step.label)
			indicator = StepPendingStyle.Render(indicator)
		}

		s.WriteString(fmt.Sprintf("%s  %s\n", indicator, label))
	}

	return SidebarStyle.Render(s.String())
}

// --- MAIN CONTENT COMPONENT ---
func renderMainContent(m MainModel) string {
	var s strings.Builder

	switch m.CurrentState {

	// --- INPUTS ---
	case StateInputProjectName:
		s.WriteString(HeaderStyle.Render("PROJECT IDENTITY"))
		s.WriteString("\n\n")
		s.WriteString("Choose a name for your new service.\n")
		s.WriteString(DescStyle.Render("Lowercase, hyphens allowed (e.g. payment-service)"))
		s.WriteString("\n\n")
		s.WriteString(m.TextInputComponent.View())

	case StateInputModuleName:
		s.WriteString(HeaderStyle.Render("MODULE PATH"))
		s.WriteString("\n\n")
		s.WriteString("Define the Go module name.\n")
		s.WriteString(DescStyle.Render("Usually github.com/user/project"))
		s.WriteString("\n\n")
		s.WriteString(m.TextInputComponent.View())

	// --- SELECTION (SINGLE) ---
	case StateSelectProjectScale, StateSelectTemplate, StateSelectFramework, StateSelectDatabaseDriver:
		var title, subtitle string
		var options []string

		// Setup Context based on state
		provider, _ := generators.GetProvider(m.ProjectScale)
		switch m.CurrentState {
		case StateSelectProjectScale:
			title = "ARCHITECTURE SCALE"
			subtitle = "How big is this project going to be?"
			options = []string{"Small (Monolith/Script)", "Medium (Standard Service)", "Enterprise (Clean Arch/Microservice)"}
		case StateSelectTemplate:
			title = "TEMPLATE VARIATION"
			subtitle = fmt.Sprintf("Available templates for %s scale:", m.ProjectScale)
			if provider != nil {
				options = provider.GetTemplates()
			}
		case StateSelectFramework:
			title = "HTTP FRAMEWORK"
			subtitle = "Select the backbone of your REST API:"
			if provider != nil {
				options = provider.GetFrameworks(m.SelectedTemplate)
			}
		case StateSelectDatabaseDriver:
			title = "DATA PERSISTENCE"
			subtitle = "Select your primary database driver:"
			if provider != nil {
				options = provider.GetDatabaseDrivers(m.SelectedTemplate)
			}
		}

		s.WriteString(HeaderStyle.Render(title))
		s.WriteString("\n\n")
		s.WriteString(subtitle + "\n\n")

		// Render List
		for i, opt := range options {
			cursor := " "
			box := "( )"
			style := lipgloss.NewStyle().Foreground(ColorBlur)

			if m.SelectedOption == i {
				cursor = "›"
				box = "(•)"
				style = lipgloss.NewStyle().Foreground(ColorFocus).Bold(true)
			}

			// Format: › (•) Option Label
			s.WriteString(style.Render(fmt.Sprintf("%s %s %s", cursor, box, opt)) + "\n")
		}

		s.WriteString("\n" + DescStyle.Render("Use Arrow Keys to move • Enter to select"))

	// --- MULTI SELECT (ADDONS) ---
	case StateSelectAddons:
		s.WriteString(HeaderStyle.Render("SYSTEM MODULES"))
		s.WriteString("\n\n")
		s.WriteString("Select additional components to install:\n\n")

		for i, addon := range core.AvailableAddons {
			cursor := " "
			box := "[ ]"
			style := lipgloss.NewStyle().Foreground(ColorBlur)

			// Logic Highlight & Checked
			isHighlighted := m.SelectedOption == i
			isChecked := m.SelectedAddonsIndices[i]

			if isChecked {
				box = "[x]"
				style = lipgloss.NewStyle().Foreground(ColorSuccess)
			}

			if isHighlighted {
				cursor = "›"

				style = style.Bold(true).Foreground(ColorText)
				if isChecked {
					style = style.Foreground(ColorSuccess).Bold(true)
				}
			}

			s.WriteString(style.Render(fmt.Sprintf("%s %s %s", cursor, box, addon.Label)) + "\n")
		}
		s.WriteString("\n" + DescStyle.Render("Space: Toggle • Enter: Confirm"))

	// --- INSTALLING ---
	case StateInstalling:
		s.WriteString(HeaderStyle.Render("FABRICATING"))
		s.WriteString("\n\n")
		s.WriteString(fmt.Sprintf("%s %s\n\n", m.Spinner.View(), m.InstallMsg))

		// Progress Bar container
		s.WriteString(m.Progress.View())

	// --- DONE ---
	case StateGenerationDone:
		s.WriteString(lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true).Render("✔ DEPLOYMENT SUCCESSFUL"))
		s.WriteString("\n\n")

		// Kotak info perintah selanjutnya
		cmdBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBlur).
			Padding(1, 2).
			Render(fmt.Sprintf("cd %s\nmake run", m.ProjectName))

		s.WriteString("Get started with:\n" + cmdBox)
		s.WriteString("\n\n" + DescStyle.Render("Press Enter to close."))
	}

	return MainContentStyle.Render(s.String())
}

func renderError(err error) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Bold(true).
		Render(fmt.Sprintf("\n  ❌ CRITICAL FAILURE: %v\n", err))
}
