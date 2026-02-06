package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xRiot45/gocrafting/internal/generators"
)

// handleEnter mengurus logika perpindahan State saat tombol Enter ditekan
func (m MainModel) handleEnter() (tea.Model, tea.Cmd) {
	var _ []tea.Cmd

	if m.CurrentState == StateGenerationDone {
		return m, tea.Quit
	}

	switch m.CurrentState {

	// STEP 1: Input Project Name
	case StateInputProjectName:
		val := m.TextInputComponent.Value()
		if val == "" {
			return m, nil
		}
		m.ProjectName = val
		m.CurrentState = StateInputModuleName
		m.TextInputComponent.Reset()
		m.TextInputComponent.Placeholder = "github.com/username/" + val
		return m, nil

	// STEP 2: Input Module Name
	case StateInputModuleName:
		val := m.TextInputComponent.Value()
		if val == "" {
			val = "github.com/username/" + m.ProjectName
		}
		m.ModuleName = val
		m.TextInputComponent.Blur()
		m.CurrentState = StateSelectProjectScale
		return m, nil

	// STEP 3: Select Scale
	case StateSelectProjectScale:
		scales := []string{"Small", "Medium", "Enterprise"}
		selected := scales[m.SelectedOption]

		if _, err := generators.GetProvider(selected); err != nil {
			m.Err = err
			return m, nil
		}

		m.ProjectScale = selected
		m.Err = nil
		m.SelectedOption = 0
		m.CurrentState = StateSelectTemplate
		return m, nil

	// STEP 4: Select Template
	case StateSelectTemplate:
		provider, _ := generators.GetProvider(m.ProjectScale)
		m.SelectedTemplate = provider.GetTemplates()[m.SelectedOption]
		m.SelectedOption = 0

		// Cek Framework
		if len(provider.GetFrameworks(m.SelectedTemplate)) > 0 {
			m.CurrentState = StateSelectFramework
			return m, nil
		}

		m.SelectedFramework = "None"
		// Logic GOTO diganti fungsi call internal
		return m.checkDatabaseStep()

	// STEP 5: Select Framework
	case StateSelectFramework:
		provider, _ := generators.GetProvider(m.ProjectScale)
		frameworks := provider.GetFrameworks(m.SelectedTemplate)
		m.SelectedFramework = frameworks[m.SelectedOption]
		m.SelectedOption = 0

		return m.checkDatabaseStep()

	// STEP 6: Select Database
	case StateSelectDatabaseDriver:
		provider, _ := generators.GetProvider(m.ProjectScale)
		dbs := provider.GetDatabaseDrivers(m.SelectedTemplate)
		m.SelectedDatabaseDriver = dbs[m.SelectedOption]

		m.CurrentState = StateSelectAddons
		m.SelectedOption = 0
		return m, nil

	// STEP 7: Select Addons (Final)
	case StateSelectAddons:
		return m.triggerInstall()
	}

	return m, nil
}

// --- SUB-FLOW FUNCTIONS (Pengganti GOTO) ---

// checkDatabaseStep menentukan apakah perlu masuk menu database atau skip ke addons
func (m MainModel) checkDatabaseStep() (tea.Model, tea.Cmd) {
	provider, _ := generators.GetProvider(m.ProjectScale)
	dbOptions := provider.GetDatabaseDrivers(m.SelectedTemplate)

	// Jika ada opsi database
	if len(dbOptions) > 0 {
		m.CurrentState = StateSelectDatabaseDriver
		return m, nil
	}

	// Jika tidak (Skip Database)
	m.SelectedDatabaseDriver = "None"
	m.CurrentState = StateSelectAddons
	m.SelectedOption = 0
	return m, nil
}

// triggerInstall memulai proses instalasi
func (m MainModel) triggerInstall() (tea.Model, tea.Cmd) {
	m.CurrentState = StateInstalling
	m.InstallMsg = "Forging project files..."

	config := m.reconstructConfig()

	// Return commands
	return m, tea.Batch(
		m.Spinner.Tick,
		generateFilesCmd(config),
	)
}
