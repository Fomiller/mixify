package mainmenuview

import (
	"log"

	"github.com/Fomiller/mixify/internal/ui/components/base"
	"github.com/Fomiller/mixify/internal/ui/context"
	"github.com/Fomiller/mixify/internal/ui/messages"
	"github.com/Fomiller/mixify/internal/ui/styles"
	"github.com/Fomiller/mixify/internal/ui/views"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	BaseComponent base.List
	ctx           *context.ProgramContext
	list          list.Model
}

type item struct {
	title, desc string
	model       tea.Model
	view        views.View
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func NewModel(ctx context.ProgramContext) Model {
	// init main model values
	m := Model{
		ctx: &ctx,
	}

	items := []list.Item{
		item{
			view:  views.CombineView,
			title: "Combine Playlists",
			desc:  "Combine your favorite playlists into one.",
		},
		item{
			view:  views.EditView,
			title: "Edit Playlists",
			desc:  "Edit your favorite playlists by adding and removing tracks.",
		},
	}

	log.Println("main menu width: ", ctx.ScreenWidth)
	log.Println("main menu height: ", ctx.ScreenHeight)
	m.list = list.New(items, list.NewDefaultDelegate(), ctx.ScreenWidth, ctx.ScreenHeight)
	log.Println(m.list.Items())
	return m
}

func (m Model) Init() tea.Cmd {
	var cmd tea.Cmd
	return tea.Batch(cmd)
}

func (m Model) View() string {
	return styles.DocStyle.Render(m.list.View())
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	log.Println("mainmenu being updated")
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case messages.BackMsg:
		m.ctx.View = views.MainMenuView

	case messages.StatusMsg:
		m.BaseComponent.Status = int(msg)
		return m, nil

	case messages.ErrMsg:
		m.BaseComponent.Err = msg
		return m, tea.Quit

	// Is it a key press?
	case tea.KeyMsg:
		switch msg.String() {
		// return to previous view with backspace
		case tea.KeyBackspace.String():
			m.ctx.View = views.MainMenuView

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter", " ":
			m.ctx.View = m.list.SelectedItem().(item).view
		}
	}

	m.list, cmd = m.list.Update(msg)

	cmds = append(
		cmds,
		cmd,
	)

	return m, tea.Batch(cmds...)
}

func (m *Model) UpdateProgramContext(ctx *context.ProgramContext) {
	m.ctx = ctx
	// not sure this is how I want to do this
	m.SetSize()
}

func (m *Model) SetSize() {
	divisor := 3
	h, v := styles.DocStyle.GetFrameSize()
	m.list.SetSize((m.ctx.ScreenWidth/divisor)-h, m.ctx.ScreenHeight-v)
}
