package ui

import tea "github.com/charmbracelet/bubbletea"

type testModel struct {
	choices []ListItem
	cursor  int
	status  int
	err     error
	state   string
	view    view
}

func newTestModel() tea.Model {
	m := testModel{}
	m.view = TEST
	return m
}

func (m testModel) Init() tea.Cmd {
	return nil
}

func (m testModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	return m, cmd
}

func (m testModel) View() string {
	return ""
}
