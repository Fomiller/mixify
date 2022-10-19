package ui

import (
	"github.com/Fomiller/mixify/pkg/ui/models"
	"github.com/Fomiller/mixify/pkg/ui/models/playlist"
	"github.com/Fomiller/mixify/pkg/ui/models/playlist/track"
	"github.com/charmbracelet/bubbles/list"
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
	list    list.Model
	choices []playlist.ListItem
	cursor  int // which item our cursor is pointing at, This could be pulled into a nested model?
	status  int
	err     error

	playlist tea.Model
	track    tea.Model
}

type item struct {
	title, desc string
	model       tea.Model
	view        view
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func New() Model {
	// init main model values
	m := Model{
		state:    MAIN,
		playlist: playlist.New(),
		track:    playlist.New(),
	}
	items := []list.Item{
		item{view: "PLAYLIST", title: "PLAYLIST", desc: "create playlists"},
		item{view: "TRACK", title: "TRACKS", desc: "edit tracks"},
	}
	m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	switch m.state {

	case "playlist":
		return m.playlist.View()

	case "track":
		return m.track.View()

	default:
		return MainMenuView(m)
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
		model, newCmd := m.playlist.Update(msg)
		// assert returned interface into struct
		playlistModel, ok := model.(playlist.Model)
		if !ok {
			panic("could not perfom assertion on playlist model")
		}
		// set cmd to the returned cmd
		cmd = newCmd
		// update the stored model
		m.playlist = playlistModel

	case TRACK:
		// return a new updated model and a cmd
		model, newCmd := m.track.Update(msg)
		// assert returned interface into struct
		trackModel, ok := model.(track.Model)
		if !ok {
			panic("could not perfom assertion on track model")
		}
		// set cmd to the returned cmd
		cmd = newCmd
		// update the stored model
		m.playlist = trackModel

	// if the state is MAIN
	default:
		switch msg := msg.(type) {
		case statusMsg:
			m.status = int(msg)
			return m, nil

		case errMsg:
			m.err = msg
			return m, tea.Quit

		case tea.WindowSizeMsg:
			h, v := docStyle.GetFrameSize()
			m.list.SetSize(msg.Width-h, msg.Height-v)

		// Is it a key press?
		case tea.KeyMsg:
			switch msg.String() {
			// return to previous view with backspace
			case tea.KeyBackspace.String():
				m.state = MAIN

			// These keys should exit the program.
			case "ctrl+c", "q":
				return m, tea.Quit

			// The "enter" key and the spacebar (a literal space) toggle
			// the selected state for the item that the cursor is pointing at.
			case "enter", " ":
				m.state = PLAYLIST
			}
		}
		m.list, cmd = m.list.Update(msg)

	}
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)

}

func MainMenuView(m Model) string {
	return docStyle.Render(m.list.View())
}
