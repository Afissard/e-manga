package tui

import tea "github.com/charmbracelet/bubbletea"

func RunTUI() error {
    p := tea.NewProgram(NewModel())

    _, err := p.Run()

    return err
}