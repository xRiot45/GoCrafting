package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/xRiot45/gocrafting/internal/core"
	"github.com/xRiot45/gocrafting/internal/generators"
)

// View renders the UI view to a string.
func (uiModel MainModel) View() string {
	if uiModel.Err != nil {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true).
			Render(fmt.Sprintf("\n❌ Error Encountered: %v\n\nPress Ctrl+C to exit.", uiModel.Err))
	}

	if uiModel.IsQuitting {
		return "Crafting cancelled.\n"
	}

	viewString := LogoStyle.Render(logoASCII) + "\n\n"
	successSymbol := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render("✔")

	// --- HEADER SECTION ---
	if uiModel.ProjectName != "" && uiModel.CurrentState > StateInputProjectName {
		viewString += fmt.Sprintf("%s %s: %s\n", successSymbol, TitleStyle.Render("Project Name"), uiModel.ProjectName)
	}
	if uiModel.ModuleName != "" && uiModel.CurrentState > StateInputModuleName {
		viewString += fmt.Sprintf("%s %s: %s\n", successSymbol, TitleStyle.Render("Module Name "), uiModel.ModuleName)
	}
	if uiModel.CurrentState > StateSelectProjectScale {
		viewString += fmt.Sprintf("%s %s: %s\n", successSymbol, TitleStyle.Render("Project Scale"), uiModel.ProjectScale)
	}
	if uiModel.SelectedTemplate != "" && uiModel.CurrentState > StateSelectTemplate {
		viewString += fmt.Sprintf("%s %s: %s\n", successSymbol, TitleStyle.Render("Template    "), uiModel.SelectedTemplate)
	}
	if uiModel.SelectedFramework != "" && uiModel.CurrentState > StateSelectFramework {
		viewString += fmt.Sprintf("%s %s: %s\n", successSymbol, TitleStyle.Render("Framework   "), uiModel.SelectedFramework)
	}
	if uiModel.SelectedDatabaseDriver != "" && uiModel.CurrentState > StateSelectDatabaseDriver {
		viewString += fmt.Sprintf("%s %s: %s\n", successSymbol, TitleStyle.Render("Database Driver "), uiModel.SelectedDatabaseDriver)
	}
	if len(uiModel.SelectedAddonsIndices) > 0 && uiModel.CurrentState > StateSelectAddons {
		var selectedAddons []string
		for idx := range uiModel.SelectedAddonsIndices {
			selectedAddons = append(selectedAddons, AddonList[idx])
		}
		viewString += fmt.Sprintf("%s %s: %s\n", successSymbol, TitleStyle.Render("Add-ons     "), strings.Join(selectedAddons, ", "))
	}

	if uiModel.CurrentState > StateInputProjectName {
		viewString += "\n"
	}

	// --- INTERACTIVE SECTION ---
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
			label, desc string
			disabled    bool
		}{
			{"Small", "Flat structure, minimal boilerplate", false},
			{"Medium", "Layered architecture, Docker ready", true},
			{"Enterprise", "Clean architecture, K8s ready", true},
		}
		for index, opt := range options {
			if opt.disabled {
				viewString += fmt.Sprintf("    %s - %s %s\n", opt.label, opt.desc, "(Coming Soon)")
			} else {
				cursor := " "
				if uiModel.SelectedOption == index {
					cursor = "›"
				}
				viewString += fmt.Sprintf("  %s %s - %s\n", cursor, opt.label, opt.desc)
			}
		}

	case StateSelectTemplate:
		viewString += TitleStyle.Render(fmt.Sprintf("Choose %s Template:", uiModel.ProjectScale)) + "\n\n"

		var options []string
		if provider, err := generators.GetProvider(uiModel.ProjectScale); err == nil {
			options = provider.GetTemplates()
		}

		for index, label := range options {
			cursor := " "
			txt := label
			if uiModel.SelectedOption == index {
				cursor = lipgloss.NewStyle().Foreground(ColorAccent).Render("›")
				txt = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true).Render(label)
			}
			viewString += fmt.Sprintf("  %s %s\n", cursor, txt)
		}

	case StateSelectFramework:
		viewString += TitleStyle.Render("Select Web Framework:") + "\n\n"

		var options []string
		if provider, err := generators.GetProvider(uiModel.ProjectScale); err == nil {
			options = provider.GetFrameworks(uiModel.SelectedTemplate)
		}

		for index, label := range options {
			cursor := " "
			txt := label
			if uiModel.SelectedOption == index {
				cursor = lipgloss.NewStyle().Foreground(ColorAccent).Render("›")
				txt = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true).Render(label)
			}
			viewString += fmt.Sprintf("  %s %s\n", cursor, txt)
		}

	case StateSelectDatabaseDriver:
		viewString += TitleStyle.Render("Select Persistence (Database):") + "\n\n"

		var options []string
		if provider, err := generators.GetProvider(uiModel.ProjectScale); err == nil {
			options = provider.GetDatabaseDrivers(uiModel.SelectedTemplate)
		}

		for index, label := range options {
			cursor := " "
			txt := label
			if uiModel.SelectedOption == index {
				cursor = lipgloss.NewStyle().Foreground(ColorAccent).Render("›")
				txt = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true).Render(label)
			}
			viewString += fmt.Sprintf("  %s %s\n", cursor, txt)
		}

	case StateSelectAddons:
		var sb strings.Builder

		// Header
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4DC4")).Bold(true).Render("Select Additional Features") + "\n")
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#5C5C5C")).Italic(true).Render("Press <SPACE> to toggle, <ENTER> to confirm") + "\n\n")

		// Loop Options dari Core
		for index, addon := range core.AvailableAddons {
			cursor := "  "
			checkbox := "[ ]"

			// Style Default (Item tidak aktif/tidak terpilih)
			itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
			cursorStyle := lipgloss.NewStyle()
			checkStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))

			// 1. Logika Kursor (Posisi yang sedang disorot)
			if uiModel.SelectedOption == index {
				cursor = "› "
				cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4DC4")).Bold(true)
				itemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
				checkStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
			}

			// 2. Logika Checkbox (Apakah item sudah dicentang?)
			if uiModel.SelectedAddonsIndices[index] {
				checkbox = "[✔]"
				checkStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true)

				if uiModel.SelectedOption != index {
					itemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#DDDDDD"))
				}
			}

			row := fmt.Sprintf("%s%s %s\n",
				cursorStyle.Render(cursor),
				checkStyle.Render(checkbox),
				itemStyle.Render(addon.Label),
			)

			sb.WriteString(row)
		}

		viewString += sb.String()

	case StateInstalling:
		viewString += "\n" + TitleStyle.Render("🚀 Initiating Launch Sequence...") + "\n\n"
		viewString += fmt.Sprintf(" %s %s\n\n", uiModel.Spinner.View(), uiModel.InstallMsg)
		viewString += " " + uiModel.Progress.View() + "\n\n"

	case StateGenerationDone:
		viewString += "\n" + lipgloss.NewStyle().Foreground(ColorGold).Bold(true).Render("✨ Project successfully forged!") + "\n"
		viewString += fmt.Sprintf("\n   cd %s\n   go run .\n\n", uiModel.ProjectName)
		viewString += HintStyle.Render("Press Enter to exit.")
	}

	return viewString
}
