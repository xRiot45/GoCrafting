// Package ui provides terminal-based user interface (TUI) logic using Bubble Tea.
package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the state of the GoCrafting TUI application.
type Model struct {
	textInput textinput.Model
	quitting  bool
}

// InitialModel initializes the application's UI model with default settings.
func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "my-awesome-project"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 30

	return Model{
		textInput: ti,
		quitting:  false,
	}
}

// Init defines the initial command to be run when the application starts.
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles incoming messages and updates the model's state accordingly.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit
		case tea.KeyEnter:
			m.quitting = true
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the current state of the UI into a string to be displayed in the terminal.
func (m Model) View() string {
	if m.quitting {
		return fmt.Sprintf("\n🎨 %s: %s\n",
			lipgloss.NewStyle().Foreground(ColorSecondary).Render("Forging your masterpiece"),
			lipgloss.NewStyle().Foreground(ColorAccent).Bold(true).Render(m.textInput.Value()),
		)
	}

	// Build the view string using styles defined in styles.go
	s := LogoStyle.Render(logoASCII) + "\n\n"

	s += TitleStyle.Render("What is the name of your masterpiece?") + "\n\n"

	// Apply accent color specifically to the input text
	m.textInput.TextStyle = lipgloss.NewStyle().Foreground(ColorAccent)
	s += "  " + m.textInput.View() + "\n\n"

	s += HintStyle.Render("  › press enter to forge, esc to quit") + "\n"

	return s
}
