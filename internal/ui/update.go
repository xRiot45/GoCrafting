package ui

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// Update adalah gerbang utama (Main Entry Point)
func (uiModel MainModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := message.(type) {

	// =====================================
	// 1. SYSTEM & ASYNC MESSAGES
	// =====================================

	case FilesCreatedMsg:
		uiModel.InstallMsg = "Downloading dependencies..."
		cmds = append(cmds, uiModel.Progress.SetPercent(0.3))
		// Reconstruct config ada di helpers.go
		config := uiModel.reconstructConfig()
		cmds = append(cmds, installDepsCmd(uiModel.ProjectName, config))

	case DepsInstalledMsg:
		uiModel.InstallMsg = "Polishing code with go fmt..."
		cmds = append(cmds, uiModel.Progress.SetPercent(0.8))
		cmds = append(cmds, formatCodeCmd(uiModel.ProjectName))

	case ProjectFormattedMsg:
		uiModel.InstallMsg = "Done!"
		cmds = append(cmds, uiModel.Progress.SetPercent(1.0))
		uiModel.CurrentState = StateGenerationDone
		return uiModel, tea.Batch(cmds...)

	case InstallErrorMsg:
		uiModel.Err = msg
		return uiModel, tea.Quit

	case spinner.TickMsg:
		var cmd tea.Cmd
		uiModel.Spinner, cmd = uiModel.Spinner.Update(msg)
		cmds = append(cmds, cmd)

	case progress.FrameMsg:
		progressModel, cmd := uiModel.Progress.Update(msg)
		uiModel.Progress = progressModel.(progress.Model)
		cmds = append(cmds, cmd)

	// =====================================
	// 2. KEYBOARD INTERACTION
	// =====================================

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			uiModel.IsQuitting = true
			return uiModel, tea.Quit

		// Delegasi ke update_keys.go
		case tea.KeyUp, tea.KeyDown, tea.KeySpace:
			return uiModel.handleNavigation(msg)

		// Delegasi ke update_flow.go
		case tea.KeyEnter:
			return uiModel.handleEnter()
		}
	}

	// Update komponen Text Input jika sedang aktif
	if uiModel.CurrentState == StateInputProjectName || uiModel.CurrentState == StateInputModuleName {
		var cmd tea.Cmd
		uiModel.TextInputComponent, cmd = uiModel.TextInputComponent.Update(message)
		cmds = append(cmds, cmd)
	}

	return uiModel, tea.Batch(cmds...)
}
