package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// TRACK
type trackModel struct {
	choices []ListItem
	cursor  int
	status  int
	err     error
	state   string
	view    view
	name    string
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
	switch msg := msg.(type) {

	case statusMsg:
		m.status = int(msg)
		return m, nil

	case errMsg:
		m.err = msg
		return m, tea.Quit

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// view tracks of selected playlists
		// case "v":
		// 	var newChoices []ListItem
		// 	for _, choice := range m.choices {
		// 		if choice.selected {
		// 			choice, ok := choice.detail.(Playlist)
		// 			if ok {
		// 				for _, value := range choice.tracks {
		// 					item := ListItem{
		// 						selected: false,
		// 						detail:   track{name: value},
		// 					}
		// 					newChoices = append(newChoices, item)
		// 				}
		// 			}

		// 		}
		// 	}

		// 	m.choices = newChoices
		// 	m.view = "playlist"
		// 	// m.viewList = append(m.viewList, m.view)
		// 	m.cursor = 0

		// return to previous view with backspace
		case tea.KeyBackspace.String():
			return m, func() tea.Msg {
				return backMsg(true)
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

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			m.view = m.choices[m.cursor].detail.(view)
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m trackModel) Name() view {
	return m.view
}

func (m trackModel) View() string {
	var output string

	for i, choice := range m.choices {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if choice.selected {
			checked = "x" // selected!
		}

		choice, ok := choice.detail.(track)
		if ok {
			// for _, track := range choice {
			// Is the cursor pointing at this choice?
			// Render the row
			output += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.name)
		}
	}
	output += fmt.Sprintf("STATE: %s\n", m.view)
	return output
}
