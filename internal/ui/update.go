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

	// --- HANDLER ASYNC (Tetap Sama) ---
	case FilesCreatedMsg:
		uiModel.InstallMsg = "Downloading dependencies..."
		cmds = append(cmds, uiModel.Progress.SetPercent(0.3))

		// Recreate config for context
		config := core.ProjectConfig{
			ProjectName:      uiModel.ProjectName,
			ModuleName:       uiModel.ModuleName,
			ProjectScale:     uiModel.ProjectScale,
			SelectedTemplate: uiModel.SelectedTemplate,
			Persistence:      uiModel.Persistence,
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
	// ----------------------------------

	case spinner.TickMsg:
		var cmd tea.Cmd
		uiModel.Spinner, cmd = uiModel.Spinner.Update(msg)
		cmds = append(cmds, cmd)

	case progress.FrameMsg:
		progressModel, cmd := uiModel.Progress.Update(msg)
		uiModel.Progress = progressModel.(progress.Model)
		cmds = append(cmds, cmd)

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
				// ... (Logic Input Name Tetap Sama)
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
				// ... (Logic Input Module Tetap Sama)
				val := uiModel.TextInputComponent.Value()
				if val == "" {
					val = "github.com/username/" + uiModel.ProjectName
				}
				uiModel.ModuleName = val
				uiModel.TextInputComponent.Blur()
				uiModel.CurrentState = StateSelectProjectScale
				return uiModel, nil

			// 1. SELECT SCALE (Titik Penentuan)
			case StateSelectProjectScale:
				scales := []string{"Small", "Medium", "Enterprise"}
				if uiModel.SelectedOption >= len(scales) {
					return uiModel, nil
				}

				selectedScale := scales[uiModel.SelectedOption]

				// VALIDASI: Cek apakah scale ini sudah ada implementasinya?
				_, err := features.GetProvider(selectedScale)
				if err != nil {
					// Jika error (misal Medium belum ada), tampilkan error di UI atau ignore
					// Disini kita return nil agar user tetap di menu ini
					uiModel.Err = err
					return uiModel, nil
				}

				// Jika valid, simpan scale dan lanjut
				uiModel.ProjectScale = selectedScale
				uiModel.Err = nil // Clear error jika ada sebelumnya
				uiModel.SelectedOption = 0
				uiModel.CurrentState = StateSelectTemplate
				return uiModel, nil

			// 2. SELECT TEMPLATE (Dinamic via Interface)
			case StateSelectTemplate:
				provider, _ := features.GetProvider(uiModel.ProjectScale)
				// UI TIDAK PEDULI INI SMALL ATAU MEDIUM
				// UI cuma panggil: provider.GetTemplates()
				options := provider.GetTemplates()

				uiModel.SelectedTemplate = options[uiModel.SelectedOption]
				uiModel.SelectedOption = 0
				uiModel.CurrentState = StateSelectPersistence
				return uiModel, nil

			// 3. SELECT PERSISTENCE (Dinamic via Interface)
			case StateSelectPersistence:
				provider, _ := features.GetProvider(uiModel.ProjectScale)
				dbOptions := provider.GetPersistenceOptions()

				if uiModel.SelectedOption >= len(dbOptions) {
					uiModel.SelectedOption = 0
				}
				uiModel.Persistence = dbOptions[uiModel.SelectedOption]

				// TRIGGER INSTALLATION
				uiModel.CurrentState = StateInstalling
				uiModel.InstallMsg = "Forging project files..."

				config := core.ProjectConfig{
					ProjectName:      uiModel.ProjectName,
					ModuleName:       uiModel.ModuleName,
					ProjectScale:     uiModel.ProjectScale,
					SelectedTemplate: uiModel.SelectedTemplate,
					Persistence:      uiModel.Persistence,
				}

				cmds = append(cmds, uiModel.Spinner.Tick)
				cmds = append(cmds, generateFilesCmd(config)) // Command generic
				return uiModel, tea.Batch(cmds...)
			}

		case tea.KeyUp:
			if uiModel.SelectedOption > 0 {
				uiModel.SelectedOption--
			}

		case tea.KeyDown:
			var maxIndex int

			// LOGIC MAX INDEX (Sangat Bersih Sekarang)
			switch uiModel.CurrentState {
			case StateSelectProjectScale:
				maxIndex = 2

			case StateSelectTemplate:
				// Tanya Provider: "Punya berapa template?"
				if provider, err := features.GetProvider(uiModel.ProjectScale); err == nil {
					maxIndex = len(provider.GetTemplates()) - 1
				}

			case StateSelectPersistence:
				// Tanya Provider: "Punya berapa opsi DB?"
				if provider, err := features.GetProvider(uiModel.ProjectScale); err == nil {
					maxIndex = len(provider.GetPersistenceOptions()) - 1
				}
			}

			if uiModel.SelectedOption < maxIndex {
				uiModel.SelectedOption++
			}
		}
	}

	if uiModel.CurrentState == StateInputProjectName || uiModel.CurrentState == StateInputModuleName {
		var cmd tea.Cmd
		uiModel.TextInputComponent, cmd = uiModel.TextInputComponent.Update(message)
		cmds = append(cmds, cmd)
	}

	return uiModel, tea.Batch(cmds...)
}
