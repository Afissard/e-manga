package tui

import (
	"e-manga/internal/command"

	tea "github.com/charmbracelet/bubbletea"
)

func NewMangaForm() *Form {
	return NewForm(
		"New Manga",
		[]Field{
			{"Name", "", FieldText},
			{"Author", "", FieldText},
			{"Summary", "", FieldText},
			{"Cover", "", FieldText},
			{"Target", "", FieldText},
			{"URL", "", FieldText},
			{"Left To Right (true/false)", "false", FieldBool},
		},
	)
}

func ProcessForm() *Form {
	return NewForm(
		"Process Manga",
		[]Field{
			{"Manga", "", FieldText},
			{"Target", "none", FieldText},
		},
	)
}

func UpdateMetadataForm() *Form {
	return NewForm(
		"Update Metadata",
		[]Field{
			{"Manga", "", FieldText},
			{"Author", "", FieldText},
			{"Summary", "", FieldText},
			{"Cover", "", FieldText},
			{"Target", "", FieldText},
			{"URL", "", FieldText},
			{"Left To Right (true/false)", "false", FieldBool},
		},
	)
}

func (m model) runSelectedCommand() tea.Cmd {
	switch m.form.title {

	case "New Manga":

		values := m.form.Values()

		opts := command.NewMangaOptions{
			Name:        values["Name"],
			Author:      values["Author"],
			Summary:     values["Summary"],
			Cover:       values["Cover"],
			Target:      values["Target"],
			URL:         values["URL"],
			LeftToRight: values["Left To Right (true/false)"] == "true",
		}

		return func() tea.Msg {
			return commandFinishedMsg{
				err: command.NewManga(opts),
			}
		}

	case "Process Manga":

		values := m.form.Values()

		opts := command.ProcessMangaOptions{
			Manga:  values["Manga"],
			Target: values["Target"],
		}

		return func() tea.Msg {
			return commandFinishedMsg{
				err: command.ProcessManga(opts),
			}
		}

	case "Update Metadata":

		values := m.form.Values()

		opts := command.UpdateMangaMetadataOptions{
			Manga:       values["Manga"],
			Author:      values["Author"],
			Summary:     values["Summary"],
			Cover:       values["Cover"],
			Target:      values["Target"],
			URL:         values["URL"],
			LeftToRight: values["Left To Right (true/false)"] == "true",
		}

		return func() tea.Msg {
			return commandFinishedMsg{
				err: command.UpdateMangaMetadata(opts),
			}
		}
	}

	return nil
}
