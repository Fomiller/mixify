package ui

import tea "github.com/charmbracelet/bubbletea"

// PLAYLIST
type playlistModel struct {
	choices []ListItem
	cursor  int
	status  int
	err     error
	state   string
	view    view
}

func newPlaylistModel() tea.Model {
	m := playlistModel{}
	m.view = PLAYLIST
	return m
}

func (m playlistModel) Init() tea.Cmd {
	return nil
}

func (m playlistModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	return m, cmd
}

func (m playlistModel) View() string {
	return ""
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

type track struct {
	name string
}

type Playlists struct {
	list []Playlist
}

type detail interface {
	Name() string
	// Description() string
	// Tracks() []string
}

type playlistDetail interface {
	detail
	Description() string
	Tracks() []string
}

type trackDetail interface {
	Name() string
}

func (p Playlist) Name() string        { return p.name }
func (p Playlist) Description() string { return p.description }
func (p Playlist) Tracks() []string    { return p.tracks }

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
