package ui

import (
	"log"

	"github.com/Fomiller/mixify/internal/ui/components/base"
	"github.com/Fomiller/mixify/internal/ui/context"
	"github.com/Fomiller/mixify/internal/ui/keys"
	"github.com/Fomiller/mixify/internal/ui/messages"
	"github.com/Fomiller/mixify/internal/ui/views"
	"github.com/Fomiller/mixify/internal/ui/views/combineview"
	"github.com/Fomiller/mixify/internal/ui/views/deleteview"
	"github.com/Fomiller/mixify/internal/ui/views/editview"
	"github.com/Fomiller/mixify/internal/ui/views/mainmenuview"
	"github.com/Fomiller/mixify/internal/ui/views/updateview"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Base base.List
	ctx  context.ProgramContext
	keys keys.KeyMap

	view         views.View
	mainMenuView mainmenuview.Model
	combineView  combineview.Model
	editView     editview.Model
	updateView   updateview.Model
	deleteView   deleteview.Model

	loaded bool
}

func NewModel() Model {
	// init main model values
	m := Model{
		keys:   keys.Keys,
		loaded: false,
		ctx: context.ProgramContext{
			View:        views.MainMenuView,
			ScreenWidth: 50, // Should just be a default and then it will get change almost immediately
		},
	}

	m.mainMenuView = mainmenuview.NewModel(m.ctx)
	// m.combineView = combineview.New(m.ctx)
	// m.editView = editview.New(m.ctx)
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	var output string
	switch m.ctx.View {

	case views.MainMenuView:
		output = m.mainMenuView.View()

	case views.CombineView:
		output = m.combineView.View()

	case views.EditView:
		output = m.editView.View()
	}
	return output
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds    []tea.Cmd
		cmd     tea.Cmd
		viewCmd tea.Cmd
	)
	log.Println("MAIN MSG: ", msg)
	log.Println("CONTEXT: ", m.ctx)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.BackSpace):
			m.view = views.MainMenuView

		// These keys should exit the program.
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

			// case Key.Matches(msg, m.keys.Select):
			// 	m.view = m.list.SelectedItem().(item).view
		}

	case messages.InitMsg:
	// any initial actions can go here

	// might not be needed
	case messages.BackMsg:
		m.view = views.MainMenuView

	case messages.StatusMsg:
		m.Base.Status = int(msg)
		return m, nil

	case messages.ErrMsg:
		m.Base.Err = msg
		return m, tea.Quit

	case tea.WindowSizeMsg:
		// this will update the main ui models context ScreenHeight/Width
		m.onWindowSizeChange(msg)
	}

	// this will sync the view models ctx with the main programs context
	// this is what will adjust the size of all the other models.
	m.syncProgramContext()

	// check current view and update
	viewCmd = m.updateCurrentView(msg)

	cmds = append(
		cmds,
		cmd,
		viewCmd,
	)
	return m, tea.Batch(cmds...)
}

// func (m *Model) switchSelectedView() views.View {
// 	// should this be from ctx?
// 	return m.view
// }

func (m *Model) onWindowSizeChange(msg tea.WindowSizeMsg) {
	m.ctx.ScreenWidth = msg.Width
	m.ctx.ScreenHeight = msg.Height
	// if the main view ever has a sidebar or header etc you,
	// could do a syncMainContentWidth/Height function to subtract the differing sizes
}

func (m *Model) syncProgramContext() {
	m.mainMenuView.UpdateProgramContext(&m.ctx)
	// m.combineView.UpdateProgramContext(&m.ctx)
	// m.editView.UpdateProgramContext(&m.ctx)
	// m.createView.UpdateProgramContext(&m.ctx)
}

func (m *Model) syncMainContentWidth() {
	// m.ctx.MainContentWidth = m.ctx.ScreenWidth - offset
}

func (m *Model) updateCurrentView(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch m.ctx.View {
	case views.MainMenuView:
		m.mainMenuView, cmd = m.mainMenuView.Update(msg)

	case views.CombineView:
		m.combineView, cmd = m.combineView.Update(msg)

	case views.EditView:
		m.editView, cmd = m.editView.Update(msg)

	}
	return cmd
}
