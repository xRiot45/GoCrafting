// Package ui provides types for the GoCrafting TUI application.
package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// View renders the UI view to a string.
func (uiModel MainModel) View() string {
	if uiModel.IsQuitting {
		return "Crafting cancelled.\n"
	}

	viewString := LogoStyle.Render(logoASCII) + "\n\n"
	successSymbol := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render("✔")

	// Render Header Status (Project & Module Name)
	if uiModel.ProjectName != "" && uiModel.CurrentState > StateInputProjectName {
		viewString += fmt.Sprintf("%s %s: %s\n", successSymbol, TitleStyle.Render("Project Name"), uiModel.ProjectName)
	}
	if uiModel.ModuleName != "" && uiModel.CurrentState > StateInputModuleName {
		viewString += fmt.Sprintf("%s %s: %s\n", successSymbol, TitleStyle.Render("Module Name "), uiModel.ModuleName)
	}

	switch uiModel.CurrentState {
	case StateInputProjectName:
		viewString += TitleStyle.Render("What is the name of your masterpiece?") + "\n"
		viewString += HintStyle.Render("(e.g., my-project)") + "\n\n"
		viewString += "  " + uiModel.TextInputComponent.View() + "\n\n"
		viewString += HintStyle.Render("› press enter to continue")

	case StateInputModuleName:
		viewString += TitleStyle.Render("Define your Go Module name:") + "\n"
		viewString += HintStyle.Render("(e.g., github.com/username/my-project)") + "\n\n"
		viewString += "  " + uiModel.TextInputComponent.View() + "\n\n"
		viewString += HintStyle.Render("› press enter to continue")

	case StateSelectProjectScale:
		viewString += TitleStyle.Render("Choose the project scale:") + "\n\n"

		options := []struct {
			label    string
			desc     string
			disabled bool
		}{
			{"Small", "Flat structure, minimal boilerplate, no-fuss", false},
			{"Medium", "Layered architecture, Docker ready, standard API", true},
			{"Enterprise", "Clean architecture, Full Observability, CI/CD, K8s", true},
		}

		for index, opt := range options {
			if opt.disabled {
				label := lipgloss.NewStyle().Foreground(ColorMuted).Render(fmt.Sprintf("  %s", opt.label))
				desc := lipgloss.NewStyle().Foreground(ColorMuted).Render(fmt.Sprintf("- %s", opt.desc))
				tag := lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Render("(Coming Soon)")

				viewString += fmt.Sprintf(" %s %s %s\n", label, desc, tag)
			} else {
				if uiModel.SelectedOption == index {
					cursor := lipgloss.NewStyle().Foreground(ColorAccent).Render("  ›")
					label := lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true).Render(opt.label)
					desc := lipgloss.NewStyle().Foreground(ColorPrimary).Render(opt.desc)

					viewString += fmt.Sprintf("%s %s - %s\n", cursor, label, desc)
				} else {
					viewString += fmt.Sprintf("    %s - %s\n", opt.label, opt.desc)
				}
			}
		}
		viewString += "\n" + HintStyle.Render("› currently only Small scale is available for forging")

	case StateGenerationDone:
		viewString += "\n" + lipgloss.NewStyle().Foreground(ColorGold).Bold(true).Render("✨ Ready to forge your masterpiece!") + "\n\n"
		viewString += fmt.Sprintf("Final Summary:\n  Project: %s\n  Module:  %s\n  Scale:   Small\n",
			uiModel.ProjectName, uiModel.ModuleName)
	}

	return viewString
}
