package ui

import (
	"fmt"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

type testModel struct {
	choices []ListItem
	cursor  int
	status  int
	err     error
	state   string
	view    view
	name    string
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

func (m testModel) Name() view {
	return m.view
}

func (m testModel) View() string {
	var output string
	// state := m.views[TEST].(testModel)
	// Send the UI for rendering
	output = fmt.Sprintf("Checking %s ... %s", url, m.view)

	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: Q%v\n\n", m.err)
	}

	if m.status > 0 {
		output += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	}
	output += fmt.Sprintf("STATE: %s\n", m.view)
	return output
}
