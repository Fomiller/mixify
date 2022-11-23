package playlistlist

import (
	"github.com/Fomiller/mixify/internal/ui/components/base"
	"github.com/Fomiller/mixify/internal/ui/messages"
	"github.com/Fomiller/mixify/internal/ui/styles"
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
	Base         base.List
	PlaylistList spotify.SimplePlaylist
	List         list.Model
	Name         string
}

func New(msg tea.WindowSizeMsg) Model {
	return Model{
		List: GetUserPlaylists(),
		Base: base.List{
			Focused: true,
			Width:   msg.Width,
			Height:  msg.Height,
		},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
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
	switch m.Base.Focused {
	case true:
		return styles.FocusedStyle.Width((m.Base.Width / 3) - h).Render(m.List.View())
	default:
		return styles.DocStyle.Width((m.Base.Width / 3) - h).Render(m.List.View())
	}
}

func (m *Model) SetWidth(width int) {
	m.Base.Width = width
}

func (m *Model) SetHeight(height int) {
	m.Base.Height = height
}
