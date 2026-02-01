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

	// StateSelectProjectSmallTemplate is the stage where user selects projec template based on small project scale.
	StateSelectProjectSmallTemplate

	// StateSelectSmallPersistence is the stage where user selects persistance based on small project scale.
	StateSelectSmallPersistence

	// StateGenerationDone is the final stage where the project is generated.
	StateGenerationDone
)

// MainModel is the main struct that stores all TUI application data.
type MainModel struct {
	CurrentState       SessionState    // Stores the user's current position
	TextInputComponent textinput.Model // Text input component
	ProjectName        string          // Project name data
	ModuleName         string          // Module name data
	ProjectScale       string          // Project scale
	SelectedOption     int             // Menu option index (0 or 1)
	SelectedTemplate   string          // Selected project template
	Persistence        string          // Selected persistence
	IsQuitting         bool            // Status of whether the user wants to exit
	Err                error           // Error message
}
