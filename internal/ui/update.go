package ui

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xRiot45/gocrafting/internal/core"
	"github.com/xRiot45/gocrafting/internal/features"
)

// Update handles messages and system events, updating the TUI state and returning commands.
func (uiModel MainModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := message.(type) {

	// =========================================================
	// 1. HANDLER ASYNC PROCESS
	// =========================================================

	case FilesCreatedMsg:
		uiModel.InstallMsg = "Downloading dependencies..."
		cmds = append(cmds, uiModel.Progress.SetPercent(0.3))

		// Recreate config for context (Dependency Installation)
		config := core.ProjectConfig{
			ProjectName:            uiModel.ProjectName,
			ModuleName:             uiModel.ModuleName,
			ProjectScale:           uiModel.ProjectScale,
			SelectedTemplate:       uiModel.SelectedTemplate,
			SelectedFramework:      uiModel.SelectedFramework,
			SelectedDatabaseDriver: uiModel.DatabaseDriver,
		}
		cmds = append(cmds, installDepsCmd(uiModel.ProjectName, config))

	case DepsInstalledMsg:
		uiModel.InstallMsg = "Polishing code with go fmt..."
		cmds = append(cmds, uiModel.Progress.SetPercent(0.8))
		cmds = append(cmds, formatCodeCmd(uiModel.ProjectName))

	case ProjectFormattedMsg:
		uiModel.InstallMsg = "Done!"
		cmds = append(cmds, uiModel.Progress.SetPercent(1.0))
		uiModel.CurrentState = StateGenerationDone
		return uiModel, tea.Batch(cmds...)

	case InstallErrorMsg:
		uiModel.Err = msg
		return uiModel, tea.Quit

	// =========================================================
	// 2. HANDLER KOMPONEN UI
	// =========================================================

	case spinner.TickMsg:
		var cmd tea.Cmd
		uiModel.Spinner, cmd = uiModel.Spinner.Update(msg)
		cmds = append(cmds, cmd)

	case progress.FrameMsg:
		progressModel, cmd := uiModel.Progress.Update(msg)
		uiModel.Progress = progressModel.(progress.Model)
		cmds = append(cmds, cmd)

	// =========================================================
	// 3. HANDLER INPUT KEYBOARD
	// =========================================================

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			uiModel.IsQuitting = true
			return uiModel, tea.Quit

		case tea.KeyEnter:
			if uiModel.CurrentState == StateGenerationDone {
				return uiModel, tea.Quit
			}

			switch uiModel.CurrentState {
			case StateInputProjectName:
				val := uiModel.TextInputComponent.Value()
				if val == "" {
					return uiModel, nil
				}
				uiModel.ProjectName = val
				uiModel.CurrentState = StateInputModuleName
				uiModel.TextInputComponent.Reset()
				uiModel.TextInputComponent.Placeholder = "github.com/username/" + val
				return uiModel, nil

			case StateInputModuleName:
				val := uiModel.TextInputComponent.Value()
				if val == "" {
					val = "github.com/username/" + uiModel.ProjectName
				}
				uiModel.ModuleName = val
				uiModel.TextInputComponent.Blur()
				uiModel.CurrentState = StateSelectProjectScale
				return uiModel, nil

			case StateSelectProjectScale:
				scales := []string{"Small", "Medium", "Enterprise"}
				if uiModel.SelectedOption >= len(scales) {
					return uiModel, nil
				}

				selectedScale := scales[uiModel.SelectedOption]

				// Validasi Provider (Mencegah memilih fitur yang belum ada)
				if _, err := features.GetProvider(selectedScale); err != nil {
					uiModel.Err = err
					return uiModel, nil
				}

				uiModel.ProjectScale = selectedScale
				uiModel.Err = nil
				uiModel.SelectedOption = 0
				uiModel.CurrentState = StateSelectTemplate
				return uiModel, nil

			// --- LOGIC FLOW PINTAR (Smart Skipping) ---

			// 1. SELECT TEMPLATE
			case StateSelectTemplate:
				provider, _ := features.GetProvider(uiModel.ProjectScale)
				options := provider.GetTemplates()

				uiModel.SelectedTemplate = options[uiModel.SelectedOption]
				uiModel.SelectedOption = 0

				// CEK: Apakah butuh Framework?
				frameworks := provider.GetFrameworks(uiModel.SelectedTemplate)
				if len(frameworks) > 0 {
					uiModel.CurrentState = StateSelectFramework
					return uiModel, nil
				}

				// JIKA TIDAK: Skip Framework -> Cek Database
				uiModel.SelectedFramework = "None"
				goto CHECK_DATABASE

			// 2. SELECT FRAMEWORK
			case StateSelectFramework:
				provider, _ := features.GetProvider(uiModel.ProjectScale)
				frameworkOptions := provider.GetFrameworks(uiModel.SelectedTemplate)

				uiModel.SelectedFramework = frameworkOptions[uiModel.SelectedOption]
				uiModel.SelectedOption = 0

				// Setelah pilih framework -> Cek Database
				goto CHECK_DATABASE

			// 3. SELECT DATABASE
			case StateSelectDatabaseDriver:
				provider, _ := features.GetProvider(uiModel.ProjectScale)
				// Gunakan template untuk filter opsi DB
				dbOptions := provider.GetDatabaseDrivers(uiModel.SelectedTemplate)

				uiModel.DatabaseDriver = dbOptions[uiModel.SelectedOption]

				// Setelah pilih DB -> Install
				goto TRIGGER_INSTALL
			}

			// LABEL: Logic pengecekan database (Reusable)
		CHECK_DATABASE:
			{
				provider, _ := features.GetProvider(uiModel.ProjectScale)
				// Kita cek opsi database berdasarkan template yang dipilih
				// (Ingat update GetDatabaseDrivers agar menerima parameter template)
				dbOptions := provider.GetDatabaseDrivers(uiModel.SelectedTemplate)

				if len(dbOptions) > 0 {
					uiModel.CurrentState = StateSelectDatabaseDriver
					return uiModel, nil
				}

				// Jika opsi kosong (misal CLI Tool), set None dan langsung Install
				uiModel.DatabaseDriver = "None"
				goto TRIGGER_INSTALL
			}

			// LABEL: Logic Memulai Instalasi (Reusable)
		TRIGGER_INSTALL:
			{
				uiModel.CurrentState = StateInstalling
				uiModel.InstallMsg = "Forging project files..."

				config := core.ProjectConfig{
					ProjectName:            uiModel.ProjectName,
					ModuleName:             uiModel.ModuleName,
					ProjectScale:           uiModel.ProjectScale,
					SelectedTemplate:       uiModel.SelectedTemplate,
					SelectedFramework:      uiModel.SelectedFramework,
					SelectedDatabaseDriver: uiModel.DatabaseDriver,
				}

				cmds = append(cmds, uiModel.Spinner.Tick)
				cmds = append(cmds, generateFilesCmd(config))
				return uiModel, tea.Batch(cmds...)
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

			case StateSelectTemplate:
				if provider, err := features.GetProvider(uiModel.ProjectScale); err == nil {
					maxIndex = len(provider.GetTemplates()) - 1
				}

			case StateSelectFramework:
				if provider, err := features.GetProvider(uiModel.ProjectScale); err == nil {
					options := provider.GetFrameworks(uiModel.SelectedTemplate)
					maxIndex = len(options) - 1
				}

			case StateSelectDatabaseDriver:
				if provider, err := features.GetProvider(uiModel.ProjectScale); err == nil {
					// PENTING: Pass SelectedTemplate ke sini
					options := provider.GetDatabaseDrivers(uiModel.SelectedTemplate)
					maxIndex = len(options) - 1
				}
			}

			if uiModel.SelectedOption < maxIndex {
				uiModel.SelectedOption++
			}
		}
	}

	// Update Text Input hanya di state tertentu
	if uiModel.CurrentState == StateInputProjectName || uiModel.CurrentState == StateInputModuleName {
		var cmd tea.Cmd
		uiModel.TextInputComponent, cmd = uiModel.TextInputComponent.Update(message)
		cmds = append(cmds, cmd)
	}

	return uiModel, tea.Batch(cmds...)
}
