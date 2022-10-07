package ui

import tea "github.com/charmbracelet/bubbletea"

// TRACK
type trackModel struct {
	choices []ListItem
	cursor  int
	status  int
	err     error
	state   string
	view    view
}

func newTrackModel() tea.Model {
	m := trackModel{}
	m.view = TRACK
	return m
}

func (m trackModel) Init() tea.Cmd {
	return nil
}

func (m trackModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	return m, cmd
}

func (m trackModel) View() string {
	return ""
}
