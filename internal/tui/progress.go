package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type commandFinishedMsg struct {
	err error
}

func (m model) updateRunning(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case commandFinishedMsg:
		m.err = msg.err
		m.state = FinishedState

		// Animate the bar to 100%
		return m, m.progress.SetPercent(1.0)
	}

	var cmd tea.Cmd

	updated, cmd := m.progress.Update(msg)
	m.progress = updated.(progress.Model)

	return m, cmd
}

func (m model) runningView() string {
	return fmt.Sprintf(
		"Running command...\n\n%s\n\nPlease wait...",
		m.progress.View(),
	)
}
