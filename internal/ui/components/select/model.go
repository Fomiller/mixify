package playlistSelect

import (
	"context"
	"log"

	"github.com/Fomiller/mixify/internal/auth"
	"github.com/Fomiller/mixify/internal/ui/messages"
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
	state        view
	Focused      bool
	List         list.Model
	PlaylistList spotify.SimplePlaylist
	cursor       int
	status       int
	err          error
	name         string
	Width        int
	Height       int
}

type Item struct {
	title    string
	desc     string
	ID       spotify.ID
	Playlist spotify.SimplePlaylist
	Selected bool
}

func (i *Item) ToggleSelected() {
	i.Selected = !i.Selected
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

func UpdateUserPlaylists() list.Model {
	var items []list.Item

	spotifyUserPlaylists, err := auth.Client.CurrentUsersPlaylists(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range spotifyUserPlaylists.Playlists {
		item := Item{
			title:    emoji.RemoveAll(p.Name),
			desc:     p.Description,
			Selected: false,
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
		Focused: true,
		List:    list,
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
	h, _ := docStyle.GetFrameSize()
	switch m.Focused {
	case true:
		return focusedStyle.Width((m.Width / 3) - h).Render(m.List.View())
	default:
		return docStyle.Width((m.Width / 3) - h).Render(m.List.View())
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) MoveToNext() tea.Msg {
	return nil
}

func (m *Model) SetWidth(width int) {
	m.Width = width
}

func (m *Model) SetHeight(height int) {
	m.Height = height
}
