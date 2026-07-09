package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) updateFinished(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {

	case tea.KeyMsg:
		return m, tea.Quit
	}

	return m, nil
}

func (m model) finishedView() string {

	if m.err != nil {
		return fmt.Sprintf(
			"❌ Command failed.\n\n%s\n\nPress any key to quit.",
			m.err,
		)
	}

	return "✅ Command completed successfully!\n\nPress any key to quit."
}