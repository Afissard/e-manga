package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	state State
	menu  list.Model

	form *Form

	progress progress.Model


	err error
}

func NewModel() model {
	items := []list.Item{
		menuItem{
			title: "New Manga",
			desc:  "Create a new manga project",
		},
		menuItem{
			title: "Process Manga",
			desc:  "Generate CBZ files",
		},
		menuItem{
			title: "Update Metadata",
			desc:  "Modify metadata",
		},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)

	l.Title = "Manga Manager"

	p := progress.New()

	return model{
		state:    MenuState,
		menu:     l,
		progress: p,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case MenuState:
		return m.updateMenu(msg)
	case FormState:
		return m.updateForm(msg)
	case RunningState:
		return m.updateRunning(msg)
	case FinishedState:
		return m.updateFinished(msg)
	default:
		return m, nil
	}
}

func (m model) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.menu.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			selected := m.menu.SelectedItem().(menuItem)

			switch selected.title {
			case "New Manga":
				m.form = NewMangaForm()

			case "Process Manga":
				m.form = ProcessForm()

			case "Update Metadata":
				m.form = UpdateMetadataForm()
			}

			m.state = FormState
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.menu, cmd = m.menu.Update(msg)

	return m, cmd
}

func (m model) updateForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		done bool
		cmd  tea.Cmd
	)

	m.form, cmd, done = m.form.Update(msg)

	if done {
		//values := m.form.Values()

		m.state = RunningState

		return m, m.runSelectedCommand()
	}

	return m, cmd
}

/*
func (m model) updateRunning(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case commandFinishedMsg:

		m.err = msg.err

		cmd := m.progress.SetPercent(1.0)

		m.state = FinishedState

		return m, cmd
	}

	var cmd tea.Cmd
	m.progress, cmd = m.progress.Update(msg)

	return m, cmd
}
*/
/*
func (m model) updateFinished(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg.(type) {

	case tea.KeyMsg:
		return m, tea.Quit
	}

	return m, nil
}
*/

func (m model) View() string {
	switch m.state {

	case MenuState:
		return m.menu.View()

	case FormState:
		return m.form.View()

	case RunningState:
		return m.runningView()

	case FinishedState:
		return m.finishedView()
	}

	return ""
}
