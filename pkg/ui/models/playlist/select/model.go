package playlistSelect

import (
	"context"
	"log"

	"github.com/Fomiller/mixify/pkg/auth"
	"github.com/Fomiller/mixify/pkg/ui/models"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zmb3/spotify/v2"
)

type view string

var (
	Focused   bool = true
	Unfocused bool = false
)

type Model struct {
	state   view
	Focused bool
	list    list.Model
	Plist   *spotify.SimplePlaylistPage
	cursor  int
	status  int
	err     error
	name    string
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func New() Model {
	var items []list.Item
	spotifyUserPlaylists, err := auth.Client.CurrentUsersPlaylists(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range spotifyUserPlaylists.Playlists {
		items = append(items, item{title: p.Name, desc: p.Description})
	}
	// TODO make this height and width dynamic for now it works
	playlistList := list.New(items, list.NewDefaultDelegate(), 60, 50)
	playlistList.KeyMap.NextPage = key.NewBinding(
		key.WithKeys("pgdown", "J"),
	)
	playlistList.KeyMap.PrevPage = key.NewBinding(
		key.WithKeys("pgup", "K"),
	)

	return Model{Focused: true, list: playlistList}
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
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

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

		case "right", "l":
			return m, cmd
		}

	}
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	switch m.Focused {
	case true:
		return focusedStyle.Render(m.list.View())
	default:
		return docStyle.Render(m.list.View())
	}
}

func (m Model) Init() tea.Cmd {
	// fmt.Println("SELECT INIT CALLED")
	// return GetUserPlaylistsCmd
	return nil
}

func (m *Model) MoveToNext() tea.Msg {
	return nil
}
