package playlist

import (
	"github.com/Fomiller/mixify/pkg/ui/models"
	"github.com/Fomiller/mixify/pkg/ui/models/playlist/combined"
	playlistSelect "github.com/Fomiller/mixify/pkg/ui/models/playlist/select"
	"github.com/Fomiller/mixify/pkg/ui/models/playlist/track"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type view string

const (
	PLAYLIST_VIEW_1 view = "VIEW_1"
	PLAYLIST_VIEW_2 view = "VIEW_2"
	PLAYLIST_VIEW_3 view = "VIEW_3"
)

type backMsg bool

type Model struct {
	state   view
	focused bool
	choices []ListItem
	cursor  int
	status  int
	err     error
	view    view
	name    string

	combined       tea.Model
	playlistSelect tea.Model
	track          tea.Model
}

func New() tea.Model {
	m := Model{
		state:          PLAYLIST_VIEW_1,
		combined:       combined.New(),
		playlistSelect: playlistSelect.New(),
		track:          track.New(),
	}

	for _, v := range PlaylistList.list {
		item := ListItem{
			Selected: false,
			Detail:   v,
		}
		m.choices = append(m.choices, item)
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// update nested models based off of state
	switch m.state {
	case PLAYLIST_VIEW_1:
		// return a new updated model and a cmd
		model, newCmd := m.playlistSelect.Update(msg)
		// assert returned interface into struct
		playlistSelectModel, ok := model.(Model)
		if !ok {
			panic("could not perfom assertion on playlist select model")
		}
		// set cmd to the returned cmd
		cmd = newCmd
		// update the stored model
		m.playlistSelect = playlistSelectModel

	case PLAYLIST_VIEW_2:
		// return a new updated model and a cmd
		model, newCmd := m.track.Update(msg)
		// assert returned interface into struct
		trackSelectModel, ok := model.(track.Model)
		if !ok {
			panic("could not perfom assertion on track select model")
		}
		// set cmd to the returned cmd
		cmd = newCmd
		// update the stored model
		m.track = trackSelectModel

	case PLAYLIST_VIEW_3:
		// return a new updated model and a cmd
		model, newCmd := m.combined.Update(msg)
		// assert returned interface into struct
		combinedPlaylistModel, ok := model.(combined.Model)
		if !ok {
			panic("could not perfom assertion on track select model")
		}
		// set cmd to the returned cmd
		cmd = newCmd
		// update the stored model
		m.combined = combinedPlaylistModel
	}

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

		// The "down" and "j" keys move the cursor down
		case "right", "l":
			return m.next(msg)

		case "left", "h":
			return m.prev(msg)

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

// Main Model view
func (m Model) View() string {
	var output string

	output = lipgloss.JoinHorizontal(0.2, m.playlistSelect.View(), m.track.View(), m.combined.View())

	// The footer
	output += "\nPress q to quit.\n"
	return output
}

func (m Model) next(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.state == PLAYLIST_VIEW_1 {
		m.state = PLAYLIST_VIEW_2

	} else if m.state == PLAYLIST_VIEW_2 {
		m.state = PLAYLIST_VIEW_3

	} else {
		return m, cmd
	}

	return m, cmd
}

func (m Model) prev(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.state == PLAYLIST_VIEW_3 {
		m.state = PLAYLIST_VIEW_2

	} else if m.state == PLAYLIST_VIEW_2 {
		m.state = PLAYLIST_VIEW_1

	} else {
		return m, cmd
	}

	return m, cmd
}

type ListItem struct {
	Selected bool
	Detail   interface{}
}

type Playlist struct {
	name        string
	description string
	tracks      []string
}

func (p Playlist) Name() string        { return p.name }
func (p Playlist) Description() string { return p.description }
func (p Playlist) Tracks() []string    { return p.tracks }

// type track struct {
// 	name string
// }

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