package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Form struct {
	title string

	fields []Field
	inputs []textinput.Model

	focus int
}

func NewForm(title string, fields []Field) *Form {
	f := &Form{
		title:  title,
		fields: fields,
	}

	f.inputs = make([]textinput.Model, len(fields))

	for i, field := range fields {
		input := textinput.New()
		input.Placeholder = field.Placeholder
		input.Prompt = field.Name + ": "

		if i == 0 {
			input.Focus()
		}

		f.inputs[i] = input
	}

	return f
}

func (f *Form) Update(msg tea.Msg) (*Form, tea.Cmd, bool) {

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "tab", "down":

			f.inputs[f.focus].Blur()

			f.focus = (f.focus + 1) % len(f.inputs)

			f.inputs[f.focus].Focus()

			return f, nil, false

		case "shift+tab", "up":

			f.inputs[f.focus].Blur()

			f.focus--

			if f.focus < 0 {
				f.focus = len(f.inputs) - 1
			}

			f.inputs[f.focus].Focus()

			return f, nil, false

		case "enter":

			if f.focus == len(f.inputs)-1 {
				return f, nil, true
			}

			f.inputs[f.focus].Blur()
			f.focus++
			f.inputs[f.focus].Focus()

			return f, nil, false
		}
	}

	var cmd tea.Cmd

	f.inputs[f.focus], cmd = f.inputs[f.focus].Update(msg)

	return f, cmd, false
}

func (f *Form) View() string {

	var b strings.Builder

	b.WriteString(f.title)
	b.WriteString("\n\n")

	for _, input := range f.inputs {
		b.WriteString(input.View())
		b.WriteRune('\n')
	}

	b.WriteString("\nTab: next field")
	b.WriteString("\nEnter: validate")

	return b.String()
}

func (f *Form) Values() map[string]string {

	values := make(map[string]string)

	for i, field := range f.fields {
		values[field.Name] = f.inputs[i].Value()
	}

	return values
}
