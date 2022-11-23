package ui

import (
	"log"

	"github.com/Fomiller/mixify/internal/ui/components/base"
	"github.com/Fomiller/mixify/internal/ui/context"
	"github.com/Fomiller/mixify/internal/ui/messages"
	"github.com/Fomiller/mixify/internal/ui/styles"
	"github.com/Fomiller/mixify/internal/ui/views/combineview"
	"github.com/Fomiller/mixify/internal/ui/views/deleteview"
	"github.com/Fomiller/mixify/internal/ui/views/editview"
	"github.com/Fomiller/mixify/internal/ui/views/mainmenuview"
	"github.com/Fomiller/mixify/internal/ui/views/updateview"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type view string

const (
	MENU     view = "menu"
	PLAYLIST view = "playlist"
	TRACK    view = "track"
	TEST     view = "test"
	CHOICE   view = "choice"
	MAIN     view = "main"
)

// MAIN MODEL
type Model struct {
	Base base.List
	list list.Model
	ctx  context.ProgramContext

	mainMenuView mainmenuview.Model
	combineView  combineview.Model
	editView     editview.Model
	updateView   updateview.Model
	deleteView   deleteview.Model

	state  view
	view   view
	loaded bool
}

type item struct {
	title, desc string
	model       tea.Model
	view        view
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func New() Model {
	// init main model values
	m := Model{
		state: MAIN,
		// playlist: playlist.New(),
		// track:    playlist.New(),
		loaded: false,
		ctx:    context.ProgramContext{},
	}
	items := []list.Item{
		item{view: PLAYLIST, title: "PLAYLIST", desc: "create playlists"},
		item{view: TRACK, title: "TRACKS", desc: "edit tracks"},
	}

	m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	switch m.state {

	case "playlist":
		return m.combineView.View()

	case "track":
		return m.editView.View()

	default:
		return MainMenuView(m)
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg.(type) {
	case messages.BackMsg:
		m.state = MAIN
	}

	// handle the update functions for views other then the main menu
	switch m.state {

	// if the state is PLAYLIST
	case PLAYLIST:
		// return a new updated model and a cmd
		model, newCmd := m.combineView.Update(msg)
		// assert returned interface into struct
		playlistModel, ok := model.(combineview.Model)
		if !ok {
			panic("could not perfom assertion on playlist model")
		}
		// set cmd to the returned cmd
		cmd = newCmd
		// update the stored model
		m.combineView = playlistModel

	// if the state is TRACK
	case TRACK:
		// return a new updated model and a cmd
		model, newCmd := m.editView.Update(msg)
		// assert returned interface into struct
		editViewModel, ok := model.(editview.Model)
		if !ok {
			panic("could not perfom assertion on track model")
		}
		// set cmd to the returned cmd
		cmd = newCmd
		// update the stored model
		m.editView = editViewModel

	// if the state is MAIN
	default:
		switch msg := msg.(type) {
		case messages.StatusMsg:
			m.Base.Status = int(msg)
			return m, nil

		case messages.ErrMsg:
			m.Base.Err = msg
			return m, tea.Quit

		case tea.WindowSizeMsg:
			log.Println("WINDOW main")
			if !m.loaded {
				m.ctx.ScreenHeight = msg.Height
				m.ctx.ScreenWidth = msg.Width
				h, v := styles.DocStyle.GetFrameSize()
				m.list.SetSize(msg.Width-h, msg.Height-v)
				m.combineView = combineview.New(msg)
				m.editView = editview.New(msg)
				m.loaded = true
			}
			// _, v := docStyle.GetFrameSize()
			// m.list.SetSize(msg.Width/2, msg.Height-v)

		// Is it a key press?
		case tea.KeyMsg:
			switch msg.String() {
			// return to previous view with backspace
			case tea.KeyBackspace.String():
				m.state = MAIN

			// These keys should exit the program.
			case "ctrl+c", "q":
				return m, tea.Quit

			// The "enter" key and the spacebar (a literal space) toggle
			// the selected state for the item that the cursor is pointing at.
			case "enter", " ":
				m.state = m.list.SelectedItem().(item).view
			}
		}
		m.list, cmd = m.list.Update(msg)

	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)

}

func MainMenuView(m Model) string {
	return styles.DocStyle.Render(m.list.View())
}
