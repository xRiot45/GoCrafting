// Package ui provides types for the GoCrafting TUI application.
package ui

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SessionState represents the currently active step in the application.
type SessionState int

const (
	// StateInputProjectName is the initial stage where user enters the directory name.
	StateInputProjectName SessionState = iota

	// StateInputModuleName is the stage where user enters the module name.
	StateInputModuleName

	// StateSelectProjectScale is the stage where user selects the project scale.
	StateSelectProjectScale

	// StateSelectTemplate is the stage where user selects the project template.
	StateSelectTemplate

	// StateSelectFramework is the stage where user selects the web framework.
	StateSelectFramework

	// StateSelectDatabaseDriver is the stage where user selects the database type.
	StateSelectDatabaseDriver

	// StateInstalling is the stage where dependencies are being installed.
	StateInstalling

	// StateGenerationDone is the final stage where the project is generated.
	StateGenerationDone
)

// MainModel is the main struct that stores all TUI application data.
type MainModel struct {
	// Data Fields
	ProjectName       string
	ModuleName        string
	ProjectScale      string
	SelectedTemplate  string
	SelectedFramework string
	DatabaseDriver    string

	// State & UI Fields
	SelectedOption     int
	CurrentState       SessionState
	IsQuitting         bool
	Err                error
	TextInputComponent textinput.Model
	Progress           progress.Model
	Spinner            spinner.Model
	InstallMsg         string
}

// New initializes and returns a new MainModel with default components.
// (Ini pengganti fungsi InitialModel Anda sebelumnya, tapi lebih lengkap)
func New() MainModel {
	// 1. Setup Text Input (Sesuai settingan lama Anda)
	ti := textinput.New()
	ti.Placeholder = "my-awesome-project"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40

	// 2. Setup Progress Bar
	prog := progress.New(progress.WithDefaultGradient())

	// 3. Setup Spinner
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return MainModel{
		TextInputComponent: ti,
		Progress:           prog,
		Spinner:            spin,
		CurrentState:       StateInputProjectName,
		SelectedOption:     0,
		IsQuitting:         false,
	}
}

// Init is the first function executed by Bubble Tea.
// It initializes the spinner tick and cursor blink.
func (uiModel MainModel) Init() tea.Cmd {
	return textinput.Blink
}
