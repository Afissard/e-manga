package tui

type menuItem struct {
    title string
    desc  string
}

func (m menuItem) Title() string {
    return m.title
}

func (m menuItem) Description() string {
    return m.desc
}

func (m menuItem) FilterValue() string {
    return m.title
}