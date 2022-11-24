package tracklist

import (
	"github.com/Fomiller/mixify/internal/ui/components/base"
	"github.com/Fomiller/mixify/internal/ui/context"
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
	ctx           *context.ProgramContext
	BaseComponent base.List
	List          list.Model
	PlaylistList  []*spotify.SimplePlaylist
}

func (m Model) Init() tea.Cmd {
	return nil
}

func NewModel(ctx context.ProgramContext) Model {
	items := []list.Item{}
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle.Foreground(lipgloss.AdaptiveColor{Light: "#1DB954", Dark: "#1DB954"})
	delegate.Styles.NormalTitle.Foreground(lipgloss.AdaptiveColor{Light: "#3FB925", Dark: "#3FB925"})

	newList := list.New(items, delegate, 60, 50)
	newList.KeyMap.NextPage = key.NewBinding(key.WithKeys("pgdown", "J"))
	newList.KeyMap.PrevPage = key.NewBinding(key.WithKeys("pgup", "K"))

	return Model{
		ctx:  &ctx,
		List: newList,
		BaseComponent: base.List{
			Focused: false,
		},
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case messages.StatusMsg:
		m.BaseComponent.Status = int(msg)
		return m, cmd

	case messages.ErrMsg:
		m.BaseComponent.Err = msg
		return m, tea.Quit

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
	divisor := 3
	h, _ := styles.DocStyle.GetFrameSize()
	switch m.BaseComponent.Focused {
	case true:
		return styles.FocusedStyle.Width((m.ctx.ScreenWidth / divisor) - h).Render(m.List.View())
	default:
		return styles.DocStyle.Width((m.ctx.ScreenWidth / divisor) - h).Render(m.List.View())
	}
}

func (m *Model) UpdateProgramContext(ctx *context.ProgramContext) {
	m.ctx = ctx
	// not sure this is how I want to do this
	m.SetSize()
}

func (m *Model) SetSize() {
	divisor := 3
	h, v := styles.DocStyle.GetFrameSize()
	m.List.SetSize((m.ctx.ScreenWidth/divisor)-h, m.ctx.ScreenHeight-v)
}
