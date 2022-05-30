package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case statusMsg:
		m.status = int(msg)
		// m.view = ""
		return m, nil

	case errMsg:
		m.err = msg
		return m, tea.Quit

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// view tracks of selected playlists
		case "v":
			m.viewList = append(m.viewList, m.view)
			m.view = "trackView"

		// return to previous view with backspace
		case tea.KeyBackspace.String():
			// set the new view to the previous view
			m.view = m.viewList[len(m.viewList)-1]
			// remove the old view
			m.viewList = m.viewList[:len(m.viewList)-1]

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

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			m.choices[m.cursor].selected = !m.choices[m.cursor].selected
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}
