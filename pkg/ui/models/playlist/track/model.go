package track

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type view string

type Model struct {
	state   view
	focused bool
	choices []ListItem
	cursor  int
	status  int
	err     error
	Name    string
}

func New() Model {
	return Model{}
}

type ListItem struct {
	selected bool
	detail   interface{}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	fmt.Println("track select")
	return m, nil
}

func (m Model) View() string {
	var output string

	output = " track "

	// The footer
	output += "\nPress q to quit.\n"
	return output
}

func (m Model) Init() tea.Cmd {
	return nil
}
