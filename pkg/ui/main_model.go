package ui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"net/http"
)

type view string

const (
	MENU     view = "menu"
	PLAYLIST view = "playlist"
	TRACK    view = "track"
	TEST     view = "test"
)

// MAIN MODEL
type mainModel struct {
	view     view
	views    map[view]tea.Model
	viewList []view //list of previous views, this could be a linked list ?

	choices   []ListItem       // items on the to-do list
	selected  map[int]Playlist // which to-do items are selected
	playlists map[int]Playlist // which to-do items are selected

	cursor int // which to-do list item our cursor is pointing at, This could be pulled into a nested model?
	status int
	err    error
	state  string
}

func New() mainModel {
	var m mainModel
	for _, v := range PlaylistList.list {
		item := ListItem{
			selected: false,
			detail:   v,
		}
		m.choices = append(m.choices, item)
	}

	m.view = MENU
	m.views = make(map[view]tea.Model)

	m.views[MENU] = newMenuModel()
	m.views[TEST] = newTestModel()
	m.views[PLAYLIST] = newTestModel()
	m.views[TRACK] = newTestModel()

	m.selected = make(map[int]Playlist)
	return m
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) View() string {
	switch m.view {

	case "menu":
		return menuView(m)

	case "playlist":
		return playlistView(m)

	case "track":
		return trackView(m)

	case "test":
		return menuView(m)

	default:
		return choiceView(m)
	}
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			var newChoices []ListItem
			for _, choice := range m.choices {
				if choice.selected {
					choice, ok := choice.detail.(Playlist)
					if ok {
						for _, value := range choice.tracks {
							item := ListItem{
								selected: false,
								detail:   track{name: value},
							}
							newChoices = append(newChoices, item)
						}
					}

				}
			}

			m.choices = newChoices
			m.view = "playlist"
			m.viewList = append(m.viewList, m.view)
			m.cursor = 0

		// return to previous view with backspace
		case tea.KeyBackspace.String():
			// OLD CODE FROM PRE REVIVAL
			// set the new view to the previous view
			// m.viewList = m.viewList[:len(m.viewList)-1]
			// m.view = m.viewList[len(m.viewList)-1]
			// m.view = "choiceView"
			// remove the old view

			// // reset choices
			// if m.view == "choiceView" {
			// 	m = NewModel()
			// }
			m.view = MENU

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

func playlistView(m mainModel) string {
	var output string
	state := m.views[PLAYLIST].(playlistModel)

	output += fmt.Sprintf("STATE: %s\n", state.view)
	return output
}

func trackView(m mainModel) string {
	var output string
	state := m.views[TRACK].(trackModel)

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
	output += fmt.Sprintf("STATE: %s\n", state.view)
	return output
}

func choiceView(m mainModel) string {
	var output string
	state := m.views[TEST].(testModel)

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if choice.selected {
			checked = "x" // selected!
		}

		// Render the row
		choice, ok := choice.detail.(Playlist)
		if ok {
			output += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.Name())
		}
	}
	output += fmt.Sprintf("STATE: %s\n", state.view)

	// The footer
	output += "\nPress q to quit.\n"
	return output
}

func serverView(m mainModel) string {
	var output string
	state := m.views[TRACK].(trackModel)

	// Send the UI for rendering
	output = fmt.Sprintf("Checking %s ... %s", url, m.view)

	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: Q%v\n\n", m.err)
	}

	if m.status > 0 {
		output += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	}
	output += fmt.Sprintf("STATE: %s\n", state.view)
	return output
}

func menuView(m mainModel) string {
	var output string
	state := m.views[MENU].(menuModel)

	// Send the UI for rendering
	output = fmt.Sprintf("Checking MENU %s ... %v", url, state.tag)

	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: Q%v\n\n", m.err)
	}

	if m.status > 0 {
		output += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	}
	output += fmt.Sprintf("STATE: %s\n", state.view)
	return output
}

func testView(m mainModel) string {
	var output string
	state := m.views[TEST].(testModel)
	// Send the UI for rendering
	output = fmt.Sprintf("Checking %s ... %s", url, m.view)

	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: Q%v\n\n", m.err)
	}

	if m.status > 0 {
		output += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	}
	output += fmt.Sprintf("STATE: %s\n", state.view)
	return output
}
