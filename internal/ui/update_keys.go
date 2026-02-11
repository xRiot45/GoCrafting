package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xRiot45/gocrafting/internal/core"
	"github.com/xRiot45/gocrafting/internal/generators"
)

// handleNavigation mengurus logika Up, Down, dan Space
func (m MainModel) handleNavigation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {

	// --- TOGGLE CHECKBOX (SPASI) ---
	case tea.KeySpace:
		if m.CurrentState == StateSelectAddons {
			idx := m.SelectedOption
			if m.SelectedAddonsIndices[idx] {
				delete(m.SelectedAddonsIndices, idx)
			} else {
				m.SelectedAddonsIndices[idx] = true
			}
		}
		return m, nil

	// --- NAVIGASI ATAS ---
	case tea.KeyUp:
		if m.SelectedOption > 0 {
			m.SelectedOption--
		}
		return m, nil

	// --- NAVIGASI BAWAH ---
	case tea.KeyDown:
		var maxIndex int
		var nextItemName string

		provider, _ := generators.GetProvider(m.ProjectScale)

		switch m.CurrentState {
		case StateSelectProjectScale:
			maxIndex = 2

		case StateSelectTemplate:
			if provider != nil {
				templates := provider.GetTemplates()
				maxIndex = len(templates) - 1

				if m.SelectedOption < maxIndex {
					nextItemName = templates[m.SelectedOption+1]
				}
			}

		case StateSelectFramework:
			if provider != nil {
				maxIndex = len(provider.GetFrameworks(m.SelectedTemplate)) - 1
			}

		case StateSelectDatabaseDriver:
			if provider != nil {
				maxIndex = len(provider.GetDatabaseDrivers(m.SelectedTemplate)) - 1
			}

		case StateSelectAddons:
			maxIndex = len(core.AvailableAddons) - 1
		}

		if m.SelectedOption < maxIndex {
			if m.CurrentState == StateSelectTemplate && isDisabledTemplate(nextItemName) {
				return m, nil
			}

			m.SelectedOption++
		}
		return m, nil
	}

	return m, nil
}
