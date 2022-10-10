package ui

import (
	"fmt"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

// MENU
type menuModel struct {
	choices []ListItem
	tag     string
	cursor  int
	status  int
	err     error
	state   string
	view    view
	name    string
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

func (m menuModel) Name() view {
	return m.view
}

func (m menuModel) View() string {
	var output string

	// Send the UI for rendering
	output = fmt.Sprintf("Checking MENU %s ... %v", url, m.tag)

	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: Q%v\n\n", m.err)
	}

	if m.status > 0 {
		output += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	}
	output += fmt.Sprintf("STATE: %s\n", m.view)
	return output
}
