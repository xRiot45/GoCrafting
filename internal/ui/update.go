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
			// Logic for input project name
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

			// Logic for input module name
			case StateInputModuleName:
				userInput := uiModel.TextInputComponent.Value()
				if userInput == "" {
					userInput = "github.com/username/" + uiModel.ProjectName
				}
				uiModel.ModuleName = userInput

				uiModel.TextInputComponent.Blur()
				uiModel.CurrentState = StateSelectProjectScale
				return uiModel, nil

			// Logic for select project scale
			case StateSelectProjectScale:
				if uiModel.SelectedOption != 0 {
					return uiModel, nil
				}

				scales := []string{"Small", "Medium", "Enterprise"}

				uiModel.ProjectScale = scales[uiModel.SelectedOption]

				if uiModel.SelectedOption == 0 {
					uiModel.CurrentState = StateSelectProjectSmallTemplate
					uiModel.SelectedOption = 0
					return uiModel, nil
				}

				return uiModel, nil

			// Logic for select project template
			case StateSelectProjectSmallTemplate:
				templates := []string{"simple-api", "fast-http", "cli-tool", "telegram-bot-starter"}
				uiModel.SelectedTemplate = templates[uiModel.SelectedOption]

				uiModel.CurrentState = StateSelectSmallPersistence
				uiModel.SelectedOption = 0
				return uiModel, nil

			// Logic for select persistence for small project
			case StateSelectSmallPersistence:
				dbOptions := []string{"none", "sqlite"}

				if uiModel.SelectedOption >= len(dbOptions) {
					uiModel.SelectedOption = 0
				}

				uiModel.Persistence = dbOptions[uiModel.SelectedOption]

				config := generator.ProjectConfig{
					ProjectName:      uiModel.ProjectName,
					ModuleName:       uiModel.ModuleName,
					ProjectScale:     "small",
					SelectedTemplate: uiModel.SelectedTemplate,
					Persistence:      uiModel.Persistence,
				}

				if err := generator.Forge(config); err != nil {
					uiModel.Err = err
					return uiModel, nil
				}

				uiModel.CurrentState = StateGenerationDone
				return uiModel, nil
			}

		case tea.KeyUp:
			if uiModel.SelectedOption > 0 {
				uiModel.SelectedOption--
			}

		case tea.KeyDown:
			var maxIndex int
			switch uiModel.CurrentState {
			case StateSelectProjectScale:
				maxIndex = 2
			case StateSelectProjectSmallTemplate:
				maxIndex = 3
			case StateSelectSmallPersistence:
				maxIndex = 1
			default:
				maxIndex = 0
			}

			if uiModel.SelectedOption < maxIndex {
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
