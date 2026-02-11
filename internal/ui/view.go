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
	// 1. Error Handling
	if m.Err != nil {
		return renderError(m.Err)
	}

	// 2. Quit Handling
	if m.IsQuitting {
		return "\n  👋 Aborting process.\n"
	}

	// 3. Render Components
	sidebar := renderSidebar(m)
	mainContent := renderMainContent(m)

	// 4. Layout Assembly (Horizontal Split)
	ui := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, mainContent)

	// 5. Final Rendering with Margin
	return lipgloss.NewStyle().Margin(1, 2).Render(ui)
}

// --- SIDEBAR COMPONENT (Progress Tracker) ---
func renderSidebar(m MainModel) string {
	var s strings.Builder

	// Logo Header
	s.WriteString(lipgloss.NewStyle().Foreground(ColorFocus).Bold(true).Render("GO CRAFTING") + "\n\n")

	// Steps Definition
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

	// Loop to render steps with status indicators
	for _, step := range steps {
		var indicator, label string

		if m.CurrentState == step.state {
			// Current Step (Active)
			indicator = "●"
			label = StepActiveStyle.Render(step.label)
			indicator = StepActiveStyle.Render(indicator)
		} else if m.CurrentState > step.state {
			// Completed Step
			indicator = "✔"
			label = StepDoneStyle.Render(step.label)
			indicator = StepDoneStyle.Render(indicator)
		} else {
			// Pending Step
			indicator = "○"
			label = StepPendingStyle.Render(step.label)
			indicator = StepPendingStyle.Render(indicator)
		}

		s.WriteString(fmt.Sprintf("%s  %s\n", indicator, label))
	}

	return SidebarStyle.Render(s.String())
}

// --- MAIN CONTENT COMPONENT (Interactive Area) ---
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

	// --- SINGLE SELECTIONS (Scale, Template, Framework, DB) ---
	case StateSelectProjectScale, StateSelectTemplate, StateSelectFramework, StateSelectDatabaseDriver:
		var title, subtitle string
		var options []string

		// Setup Context based on state
		provider, _ := generators.GetProvider(m.ProjectScale)
		switch m.CurrentState {
		case StateSelectProjectScale:
			title = "ARCHITECTURE SCALE"
			subtitle = "How big is this project going to be?"
			options = []string{"Small", "Medium", "Enterprise"}
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

		// Render List Options
		for i, opt := range options {
			cursor := " "
			box := "( )"

			// Default Style
			style := lipgloss.NewStyle().Foreground(ColorBlur)
			label := opt

			// LOGIC: Check for Disabled Items (Orange + Lock)
			// Applies to Project Scale and Templates
			isDisabled := false
			if m.CurrentState == StateSelectProjectScale && isDisabledProjectScale(opt) {
				isDisabled = true
			} else if m.CurrentState == StateSelectTemplate && isDisabledTemplate(opt) {
				isDisabled = true
			}

			if isDisabled {
				// Style for Disabled/Coming Soon
				orangeColor := lipgloss.Color("#FF8800")
				style = lipgloss.NewStyle().Foreground(orangeColor).Italic(true)
				label = fmt.Sprintf("%s (Coming Soon)", opt)
				box = "🔒 "
			} else {
				// Style for Active/Normal Items
				if m.SelectedOption == i {
					cursor = "➜"
					box = "(•)"
					style = lipgloss.NewStyle().Foreground(ColorFocus).Bold(true)
				}
			}

			// Render Row: [Cursor] [Box] [Label]
			s.WriteString(style.Render(fmt.Sprintf("%s %s %s", cursor, box, label)) + "\n")
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
				// Highlight overrides color
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

		// Success Command Box
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

// --- HELPER: ERROR RENDERING ---
func renderError(err error) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Bold(true).
		Render(fmt.Sprintf("\n  ❌ CRITICAL FAILURE: %v\n", err))
}
