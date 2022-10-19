package ui

import (
	"fmt"

	"github.com/Fomiller/mixify/pkg/ui/models"
	"github.com/Fomiller/mixify/pkg/ui/models/playlist"
	"github.com/Fomiller/mixify/pkg/ui/models/playlist/track"
	tea "github.com/charmbracelet/bubbletea"
)

type view string

const (
	MENU     view = "menu"
	PLAYLIST view = "playlist"
	TRACK    view = "track"
	TEST     view = "test"
	CHOICE   view = "choice"
	MAIN     view = "main"
)

// MAIN MODEL
type Model struct {
	state   view
	view    view
	views   map[view]tea.Model
	choices []playlist.ListItem // items on the to-do list
	cursor  int                 // which to-do list item our cursor is pointing at, This could be pulled into a nested model?
	status  int
	err     error
}

func New() Model {
	// init main model values
	m := Model{
		state: MAIN,
		views: map[view]tea.Model{
			PLAYLIST: playlist.New(),
			TRACK:    track.New(),
		},
	}

	// init choices
	for i := range m.views {
		item := playlist.ListItem{
			Selected: false,
			Detail:   i,
		}
		m.choices = append(m.choices, item)
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	switch m.state {

	case "playlist":
		return m.views[PLAYLIST].View()

	case "track":
		return m.views[TRACK].View()

	case "test":
		return m.views[TEST].View()

	default:
		return MainMenu(m)
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg.(type) {
	case models.BackMsg:
		m.state = MAIN
	}

	// handle the update functions for views other then the main menu
	switch m.state {

	case PLAYLIST:
		// return a new updated model and a cmd
		model, newCmd := m.views[PLAYLIST].Update(msg)
		// assert returned interface into struct
		playlistModel, ok := model.(playlist.Model)
		if !ok {
			panic("could not perfom assertion on playlist model")
		}
		// set cmd to the returned cmd
		cmd = newCmd
		// update the stored model
		m.views[PLAYLIST] = playlistModel

	case TRACK:
		// return a new updated model and a cmd
		model, newCmd := m.views[TRACK].Update(msg)
		// assert returned interface into struct
		trackModel, ok := model.(track.Model)
		if !ok {
			panic("could not perfom assertion on track model")
		}
		// set cmd to the returned cmd
		cmd = newCmd
		// update the stored model
		m.views[TRACK] = trackModel

	// if the state is MAIN
	default:
		switch msg := msg.(type) {
		case statusMsg:
			m.status = int(msg)
			return m, nil

		case errMsg:
			m.err = msg
			return m, tea.Quit

		// Is it a key press?
		case tea.KeyMsg:
			switch msg.String() {
			// return to previous view with backspace
			case tea.KeyBackspace.String():
				m.state = MAIN

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
				m.state = m.choices[m.cursor].Detail.(view)
			}
		}

	}
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)

}

func MainMenu(m Model) string {
	var output string

	// Iterate over our choices and create menu items
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if choice.Selected {
			checked = "x" // selected!
		}

		// Render the row
		choice := choice.Detail
		output += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)

	}

	// The footer
	output += "\nPress q to quit.\n"
	return output
}
