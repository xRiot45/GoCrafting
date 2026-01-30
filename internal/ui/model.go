// Package ui provides types for the GoCrafting TUI application.
package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// InitialModel sets up the initial state of the application when it is first run.
func InitialModel() MainModel {
	inputComponent := textinput.New()
	inputComponent.Placeholder = "my-awesome-project"
	inputComponent.Focus()
	inputComponent.CharLimit = 156
	inputComponent.Width = 40

	return MainModel{
		CurrentState:       StateInputProjectName,
		TextInputComponent: inputComponent,
		SelectedOption:     0,
		IsQuitting:         false,
	}
}

// Init is run by Bubble Tea when the program starts.
func (uiModel MainModel) Init() tea.Cmd {
	return textinput.Blink
}
