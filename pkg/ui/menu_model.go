package ui

import tea "github.com/charmbracelet/bubbletea"

// MENU
type menuModel struct {
	choices []ListItem
	tag     string
	cursor  int
	status  int
	err     error
	state   string
	view    view
}

func newMenuModel() menuModel {
	m := menuModel{}
	m.tag = "my cool tag"
	m.view = MENU
	return m
}

func (m menuModel) getTag() string {
	return m.tag
}

func (m menuModel) Init() tea.Cmd {
	return nil
}

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	return m, cmd
}

func (m menuModel) View() string {
	return ""
}
