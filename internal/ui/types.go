// Package ui provides types for the GoCrafting TUI application.
package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
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

	// StateGenerationDone is the final stage where the project is generated.
	StateGenerationDone
)

// MainModel is the main struct that stores all TUI application data.
type MainModel struct {
	CurrentState       SessionState    // Stores the user's current position
	TextInputComponent textinput.Model // Text input component
	ProjectName        string          // Project name data
	ModuleName         string          // Module name data
	SelectedOption     int             // Menu option index (0 or 1)
	IsQuitting         bool            // Status of whether the user wants to exit
}
