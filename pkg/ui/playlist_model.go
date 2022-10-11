package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type backMsg bool

type playlistModel struct {
	choices []ListItem
	cursor  int
	status  int
	err     error
	state   string
	view    view
	name    string
}

func newPlaylistModel() tea.Model {
	m := playlistModel{
		view: PLAYLIST,
	}

	for _, v := range PlaylistList.list {
		item := ListItem{
			selected: false,
			detail:   v,
		}
		m.choices = append(m.choices, item)
	}

	return m
}

func (m playlistModel) Init() tea.Cmd {
	return nil
}

func (m playlistModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.choices[m.cursor].selected = !m.choices[m.cursor].selected
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m playlistModel) Name() view {
	return m.view
}

func (m playlistModel) View() string {
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
		if choice.selected {
			checked = "x" // selected!
		}

		// Render the row
		choice := choice.detail.(Playlist)
		output += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.Name())

	}

	// The footer
	output += "\nPress q to quit.\n"
	return output
}

type ListItem struct {
	selected bool
	detail   interface{}
}

type Playlist struct {
	name        string
	description string
	tracks      []string
}

func (p Playlist) Name() string        { return p.name }
func (p Playlist) Description() string { return p.description }
func (p Playlist) Tracks() []string    { return p.tracks }

type track struct {
	name string
}

type Playlists struct {
	list []Playlist
}

type detail interface {
	Name() string
}

type playlistDetail interface {
	detail
	Description() string
	Tracks() []string
}

type trackDetail interface {
	Name() string
}

// func (d detail) FilterValue() string { return d.name }

var PlaylistList = Playlists{
	list: []Playlist{
		{
			name:        "playlist_01",
			description: "raggae music",
			tracks:      []string{"raggae 1", "raggae 2", "raggae 3", "raggae 4", "raggae 5", "raggae 6", "raggae 7", "raggae 8", "raggae 9", "raggae 10"},
		},
		{
			name:        "playlist_02",
			description: "chill music",
			tracks:      []string{"chill 1", "chill 2", "chill 3", "chill 4", "chill 5", "chill 6", "chill 7", "chill 8", "chill 9", "chill 10"},
		},
		{
			name:        "playlist_03",
			description: "rap music",
			tracks:      []string{"rap 1", "rap 2", "rap 3", "rap 4", "rap 5", "rap 6", "rap 7", "rap 8", "rap 9", "rap 10"},
		},
		{
			name:        "playlist_04",
			description: "EDM music",
			tracks:      []string{"EDM 1", "EDM 2", "EDM 3", "EDM 4", "EDM 5", "EDM 6", "EDM 7", "EDM 8", "EDM 9", "EDM 10"},
		},
		{
			name:        "playlist_05",
			description: "classical music",
			tracks:      []string{"classical 1", "classical 2", "classical 3", "classical 4", "classical 5", "classical 6", "classical 7", "classical 8", "classical 9", "classical 10"},
		},
	},
}
