package playlistlist

import (
	"github.com/Fomiller/mixify/internal/ui/components/base"
	"github.com/Fomiller/mixify/internal/ui/context"
	"github.com/Fomiller/mixify/internal/ui/styles"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zmb3/spotify/v2"
)

type Model struct {
	ctx           *context.ProgramContext
	BaseComponent base.List
	PlaylistList  spotify.SimplePlaylist
	List          list.Model
}

func NewModel(ctx context.ProgramContext) Model {
	return Model{
		ctx:  &ctx,
		List: GetUserPlaylists(),
		BaseComponent: base.List{
			Focused: true,
		},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
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
