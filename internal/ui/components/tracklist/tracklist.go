package tracklist

import (
	"log"

	"github.com/Fomiller/mixify/internal/ui/components/base"
	"github.com/Fomiller/mixify/internal/ui/messages"
	"github.com/Fomiller/mixify/internal/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zmb3/spotify/v2"
)

type view string

type Model struct {
	Base         base.List
	List         list.Model
	PlaylistList []*spotify.SimplePlaylist
}

func New(msg tea.WindowSizeMsg) Model {
	items := []list.Item{}
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle.Foreground(lipgloss.AdaptiveColor{Light: "#1DB954", Dark: "#1DB954"})
	delegate.Styles.NormalTitle.Foreground(lipgloss.AdaptiveColor{Light: "#3FB925", Dark: "#3FB925"})

	newList := list.New(items, delegate, 60, 50)
	newList.KeyMap.NextPage = key.NewBinding(key.WithKeys("pgdown", "J"))
	newList.KeyMap.PrevPage = key.NewBinding(key.WithKeys("pgup", "K"))

	return Model{
		Base: base.List{
			Focused: false,
			Width:   msg.Width,
			Height:  msg.Height,
		},
		List: newList,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case messages.StatusMsg:
		m.Base.Status = int(msg)
		return m, cmd

	case messages.ErrMsg:
		m.Base.Err = msg
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
	h, _ := styles.DocStyle.GetFrameSize()
	switch m.Base.Focused {
	case true:
		log.Println("TRACK WIDTH: ", m.Base.Width)
		return styles.FocusedStyle.Width((m.Base.Width / 3) - h).Render(m.List.View())
	default:
		return styles.DocStyle.Width((m.Base.Width / 3) - h).Render(m.List.View())
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) SetWidth(width int) {
	m.Base.Width = width
}

func (m *Model) SetHeight(height int) {
	m.Base.Height = height
}
