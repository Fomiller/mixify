package playlistSelect

import (
	"context"
	"fmt"
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
	items := []list.Item{
		item{title: "Cats", desc: "Usually"},
		item{title: "Plantasia, the album", desc: "My plants love it too"},
		item{title: "Pour over coffee", desc: "It takes forever to make though"},
		item{title: "VR", desc: "Virtual reality...what is there to say?"},
		item{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
		item{title: "Nutella", desc: "It's good on toast"},
		item{title: "Bitter melon", desc: "It cools you down"},
		item{title: "Nice socks", desc: "And by that I mean socks without holes"},
		item{title: "Eight hours of sleep", desc: "I had this once"},
		item{title: "Cats", desc: "Usually"},
		item{title: "Plantasia, the album", desc: "My plants love it too"},
		item{title: "Pour over coffee", desc: "It takes forever to make though"},
		item{title: "VR", desc: "Virtual reality...what is there to say?"},
	}
	// TODO make this height and width dynamic for now it works
	playlistList := list.New(items, list.NewDefaultDelegate(), 60, 50)
	playlistList.KeyMap.NextPage = key.NewBinding(
		key.WithKeys("pgdown", "J"),
	)
	playlistList.KeyMap.PrevPage = key.NewBinding(
		key.WithKeys("pgup", "K"),
	)
	playlist, err := auth.Client.CurrentUsersPlaylists(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return Model{Focused: true, list: playlistList, Plist: playlist}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	fmt.Println(msg)
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
	fmt.Println(m.Plist.Playlists)
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
