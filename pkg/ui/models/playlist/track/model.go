package track

import (
	"context"
	"log"

	"github.com/Fomiller/mixify/pkg/auth"
	"github.com/Fomiller/mixify/pkg/ui/models"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zmb3/spotify/v2"
)

type view string

type Model struct {
	state        view
	Focused      bool
	List         list.Model
	PlaylistList []*spotify.SimplePlaylist
	cursor       int
	status       int
	err          error
	Name         string
}

type Item struct {
	title    string
	desc     string
	ID       spotify.ID
	Selected bool
}

func (i Item) Title() string {
	if i.Selected == true {
		return selectedItemStyle.Render(i.title)
	} else {
		return i.title
	}
}
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.title }

func New() Model {
	items := []list.Item{}
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle.Foreground(lipgloss.AdaptiveColor{Light: "#1DB954", Dark: "#1DB954"})
	delegate.Styles.NormalTitle.Foreground(lipgloss.AdaptiveColor{Light: "#3FB925", Dark: "#3FB925"})

	trackList := list.New(items, delegate, 60, 50)
	trackList.KeyMap.NextPage = key.NewBinding(
		key.WithKeys("pgdown", "J"),
	)
	trackList.KeyMap.PrevPage = key.NewBinding(
		key.WithKeys("pgup", "K"),
	)

	return Model{Focused: false, List: trackList}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case models.StatusMsg:
		m.status = int(msg)
		return m, nil

	case models.ErrMsg:
		m.err = msg
		return m, tea.Quit

	case tea.WindowSizeMsg:
		// h, v := docStyle.GetFrameSize()
		m.List.SetSize(msg.Width/3, msg.Height)

	// Is it a key press?
	case tea.KeyMsg:
		switch msg.String() {
		// return to previous view with backspace
		case tea.KeyBackspace.String():
			return m, func() tea.Msg {
				return models.BackMsg(true)
			}

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		}
	}
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	switch m.Focused {
	case true:
		return focusedStyle.Render(m.List.View())
	default:
		return docStyle.Render(m.List.View())
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) PopulateTracks() {
	var tracks []spotify.PlaylistTrack
	var items []list.Item
	// get all tracks in each list
	for _, p := range m.PlaylistList {
		tracklist, err := auth.Client.GetPlaylistTracks(context.Background(), p.ID)
		if err != nil {
			log.Fatal(err)
		}

		// combine all tracks into one list
		for _, t := range tracklist.Tracks {
			tracks = append(tracks, t)
		}
	}
	// create items out of master track list
	for _, t := range tracks {
		items = append(items, Item{title: t.Track.Name, desc: t.Track.Album.Name, ID: t.Track.ID})
	}
	trackList := list.New(items, list.NewDefaultDelegate(), 60, 50)
	trackList.KeyMap.NextPage = key.NewBinding(
		key.WithKeys("pgdown", "J"),
	)
	trackList.KeyMap.PrevPage = key.NewBinding(
		key.WithKeys("pgup", "K"),
	)
	m.List = trackList
}
