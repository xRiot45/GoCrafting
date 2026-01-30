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

	if uiModel.ProjectName != "" && uiModel.CurrentState > StateInputProjectName {
		viewString += fmt.Sprintf("%s %s: %s\n",
			successSymbol,
			TitleStyle.Render("Project Name"),
			uiModel.ProjectName)
	}

	if uiModel.ModuleName != "" && uiModel.CurrentState > StateInputModuleName {
		viewString += fmt.Sprintf("%s %s: %s\n",
			successSymbol,
			TitleStyle.Render("Module Name "),
			uiModel.ModuleName)
	}

	if uiModel.CurrentState > StateSelectProjectScale {
		var scaleLabel string
		switch uiModel.SelectedOption {
		case 0:
			scaleLabel = "Small (Monolith)"
		case 1:
			scaleLabel = "Medium (Layered)"
		case 2:
			scaleLabel = "Enterprise (Clean Architecture)"
		}
		viewString += fmt.Sprintf("%s %s: %s\n",
			successSymbol,
			TitleStyle.Render("Project Scale"),
			scaleLabel)
	}

	if uiModel.ProjectName != "" && uiModel.CurrentState != StateGenerationDone {
		viewString += "\n"
	}

	switch uiModel.CurrentState {
	case StateInputProjectName:
		viewString += TitleStyle.Render("What is the name of your masterpiece?") + "\n"
		viewString += HintStyle.Render("(e.g., my-project)") + "\n\n"
		uiModel.TextInputComponent.TextStyle = InputStyle
		viewString += "  " + uiModel.TextInputComponent.View() + "\n\n"
		viewString += HintStyle.Render("› press enter to continue")

	case StateInputModuleName:
		viewString += TitleStyle.Render("Define your Go Module name:") + "\n"
		viewString += HintStyle.Render("(e.g., github.com/username/my-project)") + "\n\n"
		uiModel.TextInputComponent.TextStyle = InputStyle
		viewString += "  " + uiModel.TextInputComponent.View() + "\n\n"
		viewString += HintStyle.Render("› press enter to continue")

	case StateSelectProjectScale:
		viewString += TitleStyle.Render("Choose the project scale:") + "\n\n"
		options := []string{
			"Small      (Flat structure, minimal boilerplate, no-fuss)",
			"Medium     (Layered architecture, Docker ready, standard API)",
			"Enterprise (Clean architecture, Full Observability, CI/CD, K8s)",
		}

		for index, optionLabel := range options {
			cursorSymbol := " "
			if uiModel.SelectedOption == index {
				cursorSymbol = "›"
				cursorRendered := lipgloss.NewStyle().Foreground(ColorAccent).Render(cursorSymbol)
				optionRendered := lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true).Render(optionLabel)
				viewString += fmt.Sprintf("  %s %s\n", cursorRendered, optionRendered)
			} else {
				mutedOption := lipgloss.NewStyle().Foreground(ColorMuted).Render(optionLabel)
				viewString += fmt.Sprintf("  %s %s\n", cursorSymbol, mutedOption)
			}
		}
		viewString += "\n" + HintStyle.Render("› use arrow keys to move, enter to select")

	case StateGenerationDone:
		var scaleLabel string
		switch uiModel.SelectedOption {
		case 0:
			scaleLabel = "Small (Monolith)"
		case 1:
			scaleLabel = "Medium (Layered)"
		case 2:
			scaleLabel = "Enterprise (Clean Architecture)"
		}

		viewString += "\n" + lipgloss.NewStyle().
			Foreground(ColorGold).
			Bold(true).
			Render("✨ Ready to forge your masterpiece!") + "\n\n"

		viewString += fmt.Sprintf(
			"Final Summary:\n%s: %s\n%s: %s\n%s: %s\n",
			TitleStyle.Render("  Project "), uiModel.ProjectName,
			TitleStyle.Render("  Module  "), uiModel.ModuleName,
			TitleStyle.Render("  Scale   "), scaleLabel,
		)
	}

	return viewString
}
