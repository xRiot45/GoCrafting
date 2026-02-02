// Package ui provides types for the GoCrafting TUI application.
package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// View renders the UI view to a string.
func (uiModel MainModel) View() string {
	// 1. TANGKAP ERROR: Jika ada error, tampilkan di sini!
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

	// --- HEADER SECTION (Progress Tracking) ---

	// 1. Project Name
	if uiModel.ProjectName != "" && uiModel.CurrentState > StateInputProjectName {
		viewString += fmt.Sprintf("%s %s: %s\n",
			successSymbol,
			TitleStyle.Render("Project Name"),
			uiModel.ProjectName)
	}

	// 2. Module Name
	if uiModel.ModuleName != "" && uiModel.CurrentState > StateInputModuleName {
		viewString += fmt.Sprintf("%s %s: %s\n",
			successSymbol,
			TitleStyle.Render("Module Name "),
			uiModel.ModuleName)
	}

	// 3. Project Scale
	if uiModel.CurrentState > StateSelectProjectScale {
		viewString += fmt.Sprintf("%s %s: %s\n",
			successSymbol,
			TitleStyle.Render("Project Scale"),
			uiModel.ProjectScale)
	}

	// 4. Project Template
	if uiModel.SelectedTemplate != "" && uiModel.CurrentState > StateSelectProjectSmallTemplate {
		viewString += fmt.Sprintf("%s %s: %s\n",
			successSymbol,
			TitleStyle.Render("Template    "),
			uiModel.SelectedTemplate)
	}

	// 5. Persistence
	// Tampilkan jika sudah dipilih (Persistence tidak kosong) dan state sudah lewat atau sedang install
	if uiModel.Persistence != "" && uiModel.CurrentState > StateSelectSmallPersistence {
		viewString += fmt.Sprintf("%s %s: %s\n",
			successSymbol,
			TitleStyle.Render("Persistence "),
			uiModel.Persistence)
	}

	// Tambahkan jarak jika header sudah ada isinya
	if uiModel.CurrentState > StateInputProjectName {
		viewString += "\n"
	}

	// --- INTERACTIVE SECTION ---

	switch uiModel.CurrentState {
	// View for TUI input project name
	case StateInputProjectName:
		viewString += TitleStyle.Render("What is the name of your masterpiece?") + "\n"
		viewString += HintStyle.Render("(e.g., my-project)") + "\n\n"
		viewString += "  " + uiModel.TextInputComponent.View() + "\n\n"
		viewString += HintStyle.Render("› press enter to continue")

	// View for TUI input module name
	case StateInputModuleName:
		viewString += TitleStyle.Render("Define your Go Module name:") + "\n"
		viewString += HintStyle.Render("(e.g., github.com/username/my-project)") + "\n\n"
		viewString += "  " + uiModel.TextInputComponent.View() + "\n\n"
		viewString += HintStyle.Render("› press enter to continue")

	// View for TUI select project scale
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

	// View for TUI select project Small template
	case StateSelectProjectSmallTemplate:
		viewString += TitleStyle.Render("Choose Small Project Template:") + "\n\n"

		options := []string{
			"Simple REST API (net/http)",
			"Fast HTTP Server (Fiber/Gin)",
			"CLI Tool Template (Cobra)",
			"Telegram Bot Starter",
		}

		for index, label := range options {
			cursor := " "
			txt := label
			if uiModel.SelectedOption == index {
				cursor = lipgloss.NewStyle().Foreground(ColorAccent).Render("›")
				txt = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true).Render(label)
			}
			viewString += fmt.Sprintf("  %s [%d] %s\n", cursor, index+1, txt)
		}
		viewString += "\n" + HintStyle.Render("› use arrow keys to move, enter to select")

	// View for TUI select persistence
	case StateSelectSmallPersistence:
		viewString += TitleStyle.Render("Select Persistence (Database):") + "\n\n"

		options := []struct {
			label string
			desc  string
		}{
			{"None", "In-memory data, resets on restart"},
			{"SQLite", "Lightweight file-based database (Recommended)"},
		}

		for index, opt := range options {
			cursor := " "
			label := opt.label
			if uiModel.SelectedOption == index {
				cursor = lipgloss.NewStyle().Foreground(ColorAccent).Render("›")
				label = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true).Render(opt.label)
			}
			viewString += fmt.Sprintf("  %s %-6s - %s\n", cursor, label, HintStyle.Render(opt.desc))
		}
		viewString += "\n" + HintStyle.Render("› use arrow keys to move, enter to select")

	// --- NEW: State for Installing Progress ---
	case StateInstalling:
		viewString += "\n" + TitleStyle.Render("🚀 Initiating Launch Sequence...") + "\n\n"

		// Tampilkan Spinner & Pesan Status (misal: "Downloading dependencies...")
		viewString += fmt.Sprintf(" %s %s\n\n", uiModel.Spinner.View(), uiModel.InstallMsg)

		// Tampilkan Progress Bar
		viewString += " " + uiModel.Progress.View() + "\n\n"

		viewString += HintStyle.Render("Please wait while we craft your masterpiece...")

	// --- UPDATED: Final Success View ---
	case StateGenerationDone:
		// Pesan Sukses Besar
		viewString += "\n" + lipgloss.NewStyle().Foreground(ColorGold).Bold(true).Render("✨  Project successfully forged!") + "\n"
		viewString += lipgloss.NewStyle().Foreground(ColorMuted).Render("   You are ready to go:") + "\n\n"

		// Instruksi selanjutnya
		viewString += fmt.Sprintf("   cd %s\n", uiModel.ProjectName)
		viewString += "   go run .\n\n"

		viewString += HintStyle.Render("Press Enter to exit.")
	}

	return viewString
}
