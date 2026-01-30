package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xRiot45/gocrafting/internal/generator"
)

// Update handles messages (user input) and changes the model state.
func (uiModel MainModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var command tea.Cmd

	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			uiModel.IsQuitting = true
			return uiModel, tea.Quit

		case tea.KeyEnter:
			switch uiModel.CurrentState {
			case StateInputProjectName:
				userInput := uiModel.TextInputComponent.Value()
				if userInput == "" {
					return uiModel, nil
				}

				uiModel.ProjectName = userInput
				uiModel.CurrentState = StateInputModuleName

				uiModel.TextInputComponent.Reset()
				uiModel.TextInputComponent.Placeholder = "github.com/username/" + userInput
				return uiModel, nil

			case StateInputModuleName:
				userInput := uiModel.TextInputComponent.Value()
				if userInput == "" {
					userInput = "github.com/username/" + uiModel.ProjectName
				}
				uiModel.ModuleName = userInput

				uiModel.TextInputComponent.Blur()
				uiModel.CurrentState = StateSelectProjectScale
				return uiModel, nil

			case StateSelectProjectScale:
				uiModel.CurrentState = StateGenerationDone

				var projectScale string
				switch uiModel.SelectedOption {
				case 0:
					projectScale = "small"
				case 1:
					projectScale = "medium"
				case 2:
					projectScale = "enterprise"
				}

				config := generator.ProjectConfig{
					ProjectName:  uiModel.ProjectName,
					ModuleName:   uiModel.ModuleName,
					ProjectScale: projectScale,
				}

				if err := generator.Forge(config); err != nil {
					return uiModel, tea.Quit
				}

				return uiModel, tea.Quit
			}

		case tea.KeyUp:
			if uiModel.CurrentState == StateSelectProjectScale && uiModel.SelectedOption > 0 {
				uiModel.SelectedOption--
			}

		case tea.KeyDown:
			if uiModel.CurrentState == StateSelectProjectScale && uiModel.SelectedOption < 2 {
				uiModel.SelectedOption++
			}
		}
	}

	// Update text input components only if in typing mode Update text input components only if in typing mode
	if uiModel.CurrentState == StateInputProjectName || uiModel.CurrentState == StateInputModuleName {
		uiModel.TextInputComponent, command = uiModel.TextInputComponent.Update(message)
	}

	return uiModel, command
}
