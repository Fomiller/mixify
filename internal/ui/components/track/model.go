package track

import (
	"context"
	"fmt"
	"log"

	"github.com/Fomiller/mixify/internal/auth"
	"github.com/Fomiller/mixify/internal/ui/messages"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zmb3/spotify/v2"
)

type view string

type Model struct {
	Width        int
	Height       int
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
	ItemTitle  string
	Desc       string
	TrackID    spotify.ID
	PlaylistID spotify.ID
	Selected   bool
}

func (i Item) Title() string {
	if i.Selected == true {
		return selectedItemStyle.Render(i.ItemTitle)
	} else {
		return i.ItemTitle
	}
}
func (i Item) Description() string { return i.Desc }
func (i Item) FilterValue() string { return i.ItemTitle }

func (i *Item) ToggleSelected() {
	i.Selected = !i.Selected
}

func New(msg tea.WindowSizeMsg) Model {
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

	return Model{
		Focused: false,
		List:    trackList,
		Width:   msg.Width,
		Height:  msg.Height,
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case messages.StatusMsg:
		m.status = int(msg)
		return m, nil

	case messages.ErrMsg:
		m.err = msg
		return m, tea.Quit

	case tea.WindowSizeMsg:
		// // h, v := docStyle.GetFrameSize()
		// m.List.SetSize(msg.Width/3, msg.Height)

	// Is it a key press?
	case tea.KeyMsg:
		switch msg.String() {
		// return to previous view with backspace
		case tea.KeyBackspace.String():
			return m, func() tea.Msg {
				return messages.BackMsg(true)
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
	h, _ := docStyle.GetFrameSize()
	switch m.Focused {
	case true:
		log.Println("TRACK WIDTH: ", m.Width)
		return focusedStyle.Width((m.Width / 3) - h).Render(m.List.View())
	default:
		return docStyle.Width((m.Width / 3) - h).Render(m.List.View())
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// this needs to be optimized
func (m Model) InsertTracks(playlist spotify.SimplePlaylist) Model {
	// limit := 100
	offset := 0
	for {
		result, err := auth.Client.GetPlaylistTracks(context.Background(), playlist.ID, spotify.Offset(offset))
		if err != nil {
			panic(err)
		}

		for _, t := range result.Tracks {
			m.List.InsertItem(len(m.List.Items())+1, Item{
				ItemTitle:  t.Track.Name,
				Desc:       fmt.Sprintf("%v:%v", playlist.Name, playlist.ID),
				TrackID:    t.Track.ID,
				PlaylistID: playlist.ID,
				Selected:   true,
			})
		}

		if result.Next != "" {
			offset = offset + len(result.Tracks)
			continue
		} else {
			break
		}

	}

	return m
}

func (m Model) RemoveTracks(playlistID spotify.ID) Model {
	newList := []list.Item{}
	for _, t := range m.List.Items() {
		track, ok := t.(Item)
		if !ok {
			panic("could not assert list.Item to type Item")
		}
		if track.PlaylistID != playlistID {
			newList = append(newList, track)
		}
	}
	m.List.SetItems(newList)
	return m
}

func (m Model) GetSelectedTracks() []list.Item {
	selected := []list.Item{}
	for _, t := range m.List.Items() {
		track, ok := t.(Item)
		if !ok {
			panic("could not assert list.Item to type Item")
		}
		if track.Selected != false {
			selected = append(selected, track)
		}
	}
	m.List.SetItems(selected)
	return m.List.Items()
}

func (m *Model) SetWidth(width int) {
	m.Width = width
}

func (m *Model) SetHeight(height int) {
	m.Height = height
}
