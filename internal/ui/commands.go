package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	runner "github.com/xRiot45/gocrafting/internal/commands"
	"github.com/xRiot45/gocrafting/internal/generator"
)

// FilesCreatedMsg indicates that the project files have been successfully generated.
type FilesCreatedMsg struct{}

// DepsInstalledMsg indicates that all dependencies have been downloaded.
type DepsInstalledMsg struct{}

// ProjectFormattedMsg indicates that 'go fmt' has completed.
type ProjectFormattedMsg struct{}

// InstallErrorMsg conveys an error that occurred during the installation process.
type InstallErrorMsg error

// Cmd 1: Create File Structure
func generateFilesCmd(config generator.ProjectConfig) tea.Cmd {
	return func() tea.Msg {
		// UX: Tambah jeda sedikit biar spinner kelihatan muter
		time.Sleep(time.Millisecond * 800)

		if err := generator.Forge(config); err != nil {
			return InstallErrorMsg(err)
		}
		return FilesCreatedMsg{}
	}
}

// Cmd 2: Download Dependencies
func installDepsCmd(path string, config generator.ProjectConfig) tea.Cmd {
	return func() tea.Msg {
		// Jeda sedikit sebelum mulai download
		time.Sleep(time.Millisecond * 500)

		var packages []string
		if config.Persistence == "SQLite" || config.Persistence == "sqlite" {
			packages = append(packages, "modernc.org/sqlite")
		}

		// 1. Go Get
		if len(packages) > 0 {
			if err := runner.GoGet(path, packages...); err != nil {
				return InstallErrorMsg(err)
			}
		}

		// 2. Go Mod Tidy
		if err := runner.RunGoModTidy(path); err != nil {
			return InstallErrorMsg(err)
		}

		return DepsInstalledMsg{}
	}
}

// Cmd 3: Format Code
func formatCodeCmd(path string) tea.Cmd {
	return func() tea.Msg {
		// UX: Tambah jeda biar user sempat baca "Polishing code..."
		time.Sleep(time.Second * 1)

		if err := runner.RunGoFmt(path); err != nil {
			return InstallErrorMsg(err)
		}
		return ProjectFormattedMsg{}
	}
}
