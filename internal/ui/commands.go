package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/xRiot45/gocrafting/internal/core"
	"github.com/xRiot45/gocrafting/internal/features/small"
	"github.com/xRiot45/gocrafting/internal/runner"
)

// --- MESSAGE TYPES ---
// (Ini memperbaiki error undefined: FilesCreatedMsg, dll)

// FilesCreatedMsg indicates that the project files have been successfully generated.
type FilesCreatedMsg struct{}

// DepsInstalledMsg indicates that all dependencies have been downloaded.
type DepsInstalledMsg struct{}

// ProjectFormattedMsg indicates that 'go fmt' has completed.
type ProjectFormattedMsg struct{}

// InstallErrorMsg conveys an error that occurred during the installation process.
type InstallErrorMsg error

// --- COMMAND FUNCTIONS ---

// Cmd 1: Generate Files (Memanggil Logic di Features)
func generateFilesCmd(config core.ProjectConfig) tea.Cmd {
	return func() tea.Msg {
		// UX: Tambah jeda sedikit biar spinner kelihatan muter
		time.Sleep(time.Millisecond * 800)

		var err error

		// ROUTING: Panggil generator sesuai Scale
		switch config.ProjectScale {
		case "Small":
			err = small.Generate(config)
		case "Medium":
			// err = medium.Generate(config) // Nanti
		default:
			// Default fallback jika scale tidak dikenali
			err = small.Generate(config)
		}

		if err != nil {
			return InstallErrorMsg(err)
		}
		return FilesCreatedMsg{}
	}
}

// Cmd 2: Download Dependencies (Dummy visual, karena logic asli ada di step 1)
func installDepsCmd(_ string, _ core.ProjectConfig) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(time.Millisecond * 1000)
		return DepsInstalledMsg{}
	}
}

// Cmd 3: Format Code (Memanggil Runner)
// (Ini memperbaiki error undefined: formatCodeCmd)
func formatCodeCmd(path string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(time.Second * 1)

		if err := runner.RunGoFmt(path); err != nil {
			return InstallErrorMsg(err)
		}
		return ProjectFormattedMsg{}
	}
}
