package combined

import (
	"fmt"

	"github.com/Fomiller/mixify/pkg/ui/models"
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
	name    string
}

type ListItem struct {
	Selected bool
	Detail   interface{}
}

func New() Model {
	return Model{}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	fmt.Println("combined")
	switch msg := msg.(type) {

	case models.StatusMsg:
		m.status = int(msg)
		return m, nil

	case models.ErrMsg:
		m.err = msg
		return m, tea.Quit

	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {
		// return to previous view with backspace
		case tea.KeyBackspace.String():
			return m, func() tea.Msg {
				return models.BackMsg(true)
			}

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "down" and "j" keys move the cursor down
		// case "right", "l":
		// return m.next(msg)

		// case "left", "h":
		// return m.prev(msg)

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			m.choices[m.cursor].Selected = !m.choices[m.cursor].Selected
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, cmd
}

func (m Model) View() string {
	var output string

	output = " combined view "

	// The footer
	output += "\nPress q to quit.\n"
	return output
}

func (m Model) Init() tea.Cmd {
	return nil
}
