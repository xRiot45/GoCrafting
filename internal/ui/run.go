package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Start memulai TUI GoCrafting.
func Start(initialName string) error {
	m := NewMainModel()

	if initialName != "" {
		m.ProjectName = initialName
		m.TextInputComponent.SetValue(initialName)
	}

	// Jalankan Bubble Tea
	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}

// NewMainModel creates and returns a new MainModel instance.
func NewMainModel() MainModel {
	return InitialModel()
}
