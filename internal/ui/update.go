package ui

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xRiot45/gocrafting/internal/core"
	"github.com/xRiot45/gocrafting/internal/features/small"
)

// Update handles messages (user input) and changes the model state.
func (uiModel MainModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := message.(type) {

	// =========================================================
	// 1. HANDLER ASYNC PROCESS (Generate -> Install -> Format)
	// =========================================================

	// Step 1 Selesai: File Berhasil Dibuat
	case FilesCreatedMsg:
		uiModel.InstallMsg = "Downloading dependencies..."
		// Update progress ke 30%
		progressCmd := uiModel.Progress.SetPercent(0.3)
		cmds = append(cmds, progressCmd)

		// Lanjut ke Step 2: Download Dependencies
		// Kita buat config ulang atau passing dari sebelumnya (disini kita recreate sederhana)
		config := core.ProjectConfig{
			ProjectName:      uiModel.ProjectName,
			ModuleName:       uiModel.ModuleName,
			ProjectScale:     uiModel.ProjectScale,
			SelectedTemplate: uiModel.SelectedTemplate,
			Persistence:      uiModel.Persistence,
		}
		// Jalankan command install (Visual Only, karena logic asli sudah di features)
		cmds = append(cmds, installDepsCmd(uiModel.ProjectName, config))

	// Step 2 Selesai: Dependencies Berhasil Diinstall
	case DepsInstalledMsg:
		uiModel.InstallMsg = "Polishing code with go fmt..."
		// Update progress ke 80%
		progressCmd := uiModel.Progress.SetPercent(0.8)
		cmds = append(cmds, progressCmd)

		// Lanjut ke Step 3: Formatting Code
		cmds = append(cmds, formatCodeCmd(uiModel.ProjectName))

	// Step 3 Selesai: Formatting Selesai (FINAL)
	case ProjectFormattedMsg:
		uiModel.InstallMsg = "Done!"
		// Update progress ke 100%
		progressCmd := uiModel.Progress.SetPercent(1.0)
		cmds = append(cmds, progressCmd)

		// Pindah ke State Selesai untuk menampilkan Summary
		uiModel.CurrentState = StateGenerationDone

		// Opsional: Jika ingin langsung keluar otomatis setelah selesai, uncomment baris ini:
		// return uiModel, tea.Quit

		return uiModel, tea.Batch(cmds...)

	// Jika Terjadi Error di tengah proses
	case InstallErrorMsg:
		uiModel.Err = msg
		return uiModel, tea.Quit

	// =========================================================
	// 2. HANDLER KOMPONEN UI (Spinner & Progress Bar)
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
			// Jika sudah selesai generated, Enter akan keluar aplikasi
			if uiModel.CurrentState == StateGenerationDone {
				return uiModel, tea.Quit
			}

			switch uiModel.CurrentState {
			// Logic for input project name
			case StateInputProjectName:
				userInput := uiModel.TextInputComponent.Value()
				if userInput == "" {
					return uiModel, nil
				}
				uiModel.ProjectName = userInput
				uiModel.CurrentState = StateInputModuleName
				uiModel.TextInputComponent.Reset()
				uiModel.TextInputComponent.Placeholder = "github.com/username/" + userInput
				return uiModel, nil

			// Logic for input module name
			case StateInputModuleName:
				userInput := uiModel.TextInputComponent.Value()
				if userInput == "" {
					userInput = "github.com/username/" + uiModel.ProjectName
				}
				uiModel.ModuleName = userInput
				uiModel.TextInputComponent.Blur()
				uiModel.CurrentState = StateSelectProjectScale
				return uiModel, nil

			// Logic for select project scale
			case StateSelectProjectScale:
				// Saat ini hanya Small yang aktif (index 0)
				// Jika pilih Medium/Enterprise (index > 0), kita block dulu atau biarkan
				if uiModel.SelectedOption != 0 {
					return uiModel, nil
				}

				scales := []string{"Small", "Medium", "Enterprise"}
				uiModel.ProjectScale = scales[uiModel.SelectedOption]

				// Reset option untuk menu selanjutnya
				uiModel.SelectedOption = 0

				// Pindah ke pemilihan Template (Generic State)
				uiModel.CurrentState = StateSelectTemplate
				return uiModel, nil

			// Logic for select project template (DINAMIS)
			case StateSelectTemplate:
				var options []string

				// Cek Scale apa yang aktif, ambil opsi dari feature terkait
				switch uiModel.ProjectScale {
				case "Small":
					options = small.GetTemplates()
				case "Medium":
					// options = medium.GetTemplates()
					options = []string{"Standard API"} // Placeholder
				default:
					options = []string{"Default"}
				}

				// Simpan pilihan template
				uiModel.SelectedTemplate = options[uiModel.SelectedOption]

				// Pindah ke pemilihan Persistence
				uiModel.SelectedOption = 0
				uiModel.CurrentState = StateSelectPersistence
				return uiModel, nil

			// Logic for select persistence (TRIGGER INSTALL)
			case StateSelectPersistence:
				var dbOptions []string

				// Cek Scale apa yang aktif, ambil opsi DB dari feature terkait
				switch uiModel.ProjectScale {
				case "Small":
					dbOptions = small.GetPersistence()
				case "Medium":
					// dbOptions = medium.GetPersistence()
					dbOptions = []string{"None"} // Placeholder
				default:
					dbOptions = []string{"None"}
				}

				// Safety check index
				if uiModel.SelectedOption >= len(dbOptions) {
					uiModel.SelectedOption = 0
				}
				uiModel.Persistence = dbOptions[uiModel.SelectedOption]

				// --- MULAI INSTALLASI ---
				uiModel.CurrentState = StateInstalling
				uiModel.InstallMsg = "Forging project files..."

				config := core.ProjectConfig{
					ProjectName:      uiModel.ProjectName,
					ModuleName:       uiModel.ModuleName,
					ProjectScale:     uiModel.ProjectScale, // Gunakan dynamic scale
					SelectedTemplate: uiModel.SelectedTemplate,
					Persistence:      uiModel.Persistence,
				}

				// Mulai Spinner dan Jalankan Command Generate Files
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

			// Hitung Max Index secara dinamis berdasarkan State & Scale
			switch uiModel.CurrentState {
			case StateSelectProjectScale:
				maxIndex = 2 // Small, Medium, Enterprise

			case StateSelectTemplate:
				// Hitung jumlah template yang tersedia
				if uiModel.ProjectScale == "Small" {
					maxIndex = len(small.GetTemplates()) - 1
				} else {
					maxIndex = 0 // Placeholder medium
				}

			case StateSelectPersistence:
				// Hitung jumlah opsi DB yang tersedia
				if uiModel.ProjectScale == "Small" {
					maxIndex = len(small.GetPersistence()) - 1
				} else {
					maxIndex = 0
				}

			default:
				maxIndex = 0
			}

			if uiModel.SelectedOption < maxIndex {
				uiModel.SelectedOption++
			}
		}
	}

	// Update text input components only if in typing mode
	if uiModel.CurrentState == StateInputProjectName || uiModel.CurrentState == StateInputModuleName {
		var cmd tea.Cmd
		uiModel.TextInputComponent, cmd = uiModel.TextInputComponent.Update(message)
		cmds = append(cmds, cmd)
	}

	return uiModel, tea.Batch(cmds...)
}
