package playlistlist

import (
	"context"
	"log"

	"github.com/Fomiller/mixify/internal/auth"
	"github.com/Fomiller/mixify/internal/ui/components/basecomponents"
	"github.com/Fomiller/mixify/internal/ui/components/playlist"
	"github.com/Fomiller/mixify/internal/ui/messages"
	"github.com/Fomiller/mixify/internal/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	emoji "github.com/tmdvs/Go-Emoji-Utils"
	"github.com/zmb3/spotify/v2"
)

type view string

var (
	Focused   bool = true
	Unfocused bool = false
)

type Model struct {
	BaseList     basecomponents.List
	PlaylistList spotify.SimplePlaylist
	List         list.Model
	Name         string
}

func UpdateUserPlaylists() list.Model {
	var items []list.Item

	spotifyUserPlaylists, err := auth.Client.CurrentUsersPlaylists(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range spotifyUserPlaylists.Playlists {
		item := playlist.Playlist{
			PlaylistTitle: emoji.RemoveAll(p.Name),
			Desc:          p.Description,
			Item: basecomponents.Item{
				Selected: false,
			},
			Playlist: p,
		}
		items = append(items, item)
	}
	// TODO make this height and width dynamic for now it works
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle.Foreground(lipgloss.AdaptiveColor{Light: "#1DB954", Dark: "#1DB954"})
	delegate.Styles.NormalTitle.Foreground(lipgloss.AdaptiveColor{Light: "#1DB925", Dark: "#1DB925"})
	playlistList := list.New(items, delegate, 60, 50)
	playlistList.KeyMap.NextPage = key.NewBinding(
		key.WithKeys("pgdown", "J"),
	)
	playlistList.KeyMap.PrevPage = key.NewBinding(
		key.WithKeys("pgup", "K"),
	)
	return playlistList
}

func New(msg tea.WindowSizeMsg) Model {
	list := UpdateUserPlaylists()
	return Model{
		List: list,
		BaseList: basecomponents.List{
			Focused: true,
			Width:   msg.Width,
			Height:  msg.Height,
		},
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case messages.StatusMsg:
		m.BaseList.Status = int(msg)
		return m, nil

	case messages.ErrMsg:
		m.BaseList.Err = msg
		return m, tea.Quit

	case tea.WindowSizeMsg:

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

		case "right", "l":
			return m, cmd

		}
	}
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	h, _ := styles.DocStyle.GetFrameSize()
	switch m.BaseList.Focused {
	case true:
		return styles.FocusedStyle.Width((m.BaseList.Width / 3) - h).Render(m.List.View())
	default:
		return styles.DocStyle.Width((m.BaseList.Width / 3) - h).Render(m.List.View())
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) MoveToNext() tea.Msg {
	return nil
}

func (m *Model) SetWidth(width int) {
	m.BaseList.Width = width
}

func (m *Model) SetHeight(height int) {
	m.BaseList.Height = height
}
