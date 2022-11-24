package combineview

import (
	"log"

	"github.com/Fomiller/mixify/internal/ui/commands"
	"github.com/Fomiller/mixify/internal/ui/components/base"
	"github.com/Fomiller/mixify/internal/ui/components/playlist"
	"github.com/Fomiller/mixify/internal/ui/components/playlistlist"
	"github.com/Fomiller/mixify/internal/ui/components/previewlist"
	"github.com/Fomiller/mixify/internal/ui/components/textinput"
	"github.com/Fomiller/mixify/internal/ui/components/track"
	"github.com/Fomiller/mixify/internal/ui/components/tracklist"
	"github.com/Fomiller/mixify/internal/ui/context"
	"github.com/Fomiller/mixify/internal/ui/messages"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	PLAYLIST_VIEW_1 view = "VIEW_1"
	PLAYLIST_VIEW_2 view = "VIEW_2"
	PLAYLIST_VIEW_3 view = "VIEW_3"
	PLAYLIST_VIEW_4 view = "VIEW_4"
)

type view string

type Model struct {
	BaseComponent base.List
	ctx           *context.ProgramContext

	playlistlist playlistlist.Model
	tracklist    tracklist.Model
	previewlist  previewlist.Model
	confirm      textinput.Model

	loaded bool
	state  view
	Width  int
	Height int
}

func NewModel(ctx context.ProgramContext) Model {
	m := Model{
		ctx:          &ctx,
		state:        PLAYLIST_VIEW_1,
		loaded:       false,
		Width:        ctx.ScreenWidth,
		Height:       ctx.ScreenHeight,
		previewlist:  previewlist.NewModel(ctx),
		playlistlist: playlistlist.NewModel(ctx),
		tracklist:    tracklist.NewModel(ctx),
		confirm:      textinput.NewModel(),
	}

	return m
}

func (m Model) ResetModel(ctx *context.ProgramContext) Model {
	return Model{
		ctx:          m.ctx,
		state:        PLAYLIST_VIEW_1,
		previewlist:  previewlist.NewModel(*m.ctx),
		playlistlist: playlistlist.NewModel(*m.ctx),
		tracklist:    tracklist.NewModel(*m.ctx),
		confirm:      textinput.NewModel(),
	}
}

func (m Model) Init() tea.Cmd {
	return commands.GetUserPlaylistsCmd
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	// update nested models based off of state
	switch m.state {
	case PLAYLIST_VIEW_1:
		m.playlistlist, cmd = m.playlistlist.Update(msg)
		cmds = append(cmds, cmd)

	case PLAYLIST_VIEW_2:
		m.tracklist, cmd = m.tracklist.Update(msg)
		cmds = append(cmds, cmd)

	case PLAYLIST_VIEW_3:
		m.previewlist, cmd = m.previewlist.Update(msg)
		cmds = append(cmds, cmd)

	case PLAYLIST_VIEW_4:
		m.confirm, cmd = m.confirm.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case messages.CreatePlaylistMsg:
		name := m.confirm.Inputs[0].Value()
		desc := m.confirm.Inputs[1].Value()
		err := m.previewlist.CreatePlaylist(name, desc)
		if err != nil {
			log.Fatal(err)
		}
		m = m.ResetModel(m.ctx)
		return m, nil

	case messages.ResetStateMsg:
		m.state = PLAYLIST_VIEW_1
		return m, nil

	case messages.StatusMsg:
		m.BaseComponent.Status = int(msg)
		return m, nil

	case messages.ErrMsg:
		m.BaseComponent.Err = msg
		return m, tea.Quit

	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {
		// return to previous view with backspace
		case tea.KeyBackspace.String():
			switch m.state {
			case PLAYLIST_VIEW_4:
			// override backspace to allow for text input
			default:
				return m, func() tea.Msg {
					return messages.BackMsg(true)
				}
			}

		// These keys should exit the program.
		case "esc":
			switch m.state {
			case PLAYLIST_VIEW_4:
				m.state = PLAYLIST_VIEW_3
			}
			return m, nil

		case "ctrl+c", "q":
			return m, tea.Quit

		// The "down" and "j" keys move the cursor down
		case "right", "l":
			return m.next(msg)

		case "left", "h":
			return m.prev(msg)

		case "enter", " ":
			switch m.state {
			case PLAYLIST_VIEW_1:
				item := m.playlistlist.List.SelectedItem().(playlist.Playlist)
				cursor := m.playlistlist.List.Index()

				if item.BaseComponent.Selected == false {
					item.ToggleSelected()
					m.playlistlist.List.SetItem(cursor, item)
					m.tracklist = m.tracklist.InsertTracks(item.Playlist)
					selectedTracks := m.tracklist.GetSelectedTracks()
					m.previewlist.List.SetItems(selectedTracks)
				} else {
					item.ToggleSelected()
					m.playlistlist.List.SetItem(cursor, item)
					m.tracklist = m.tracklist.RemoveTracks(item.Playlist.ID)
					selectedTracks := m.tracklist.GetSelectedTracks()
					m.previewlist.List.SetItems(selectedTracks)

				}
				return m, nil

			case PLAYLIST_VIEW_2:
				item := m.tracklist.List.SelectedItem().(track.Track)
				cursor := m.tracklist.List.Index()

				item.ToggleSelected()
				m.tracklist.List.SetItem(cursor, item)
				selectedTracks := m.tracklist.GetSelectedTracks()
				m.previewlist.List.SetItems(selectedTracks)
				return m, nil

			case PLAYLIST_VIEW_3:
				m.state = PLAYLIST_VIEW_4
				return m, nil
			}
		}
	}
	return m, tea.Batch(cmds...)
}

// Main Model view
func (m Model) View() string {
	var output string

	if m.state == PLAYLIST_VIEW_4 {
		output = m.confirm.View()
	} else {
		output = lipgloss.JoinHorizontal(lipgloss.Center, m.playlistlist.View(), m.tracklist.View(), m.previewlist.View())
	}
	return output
}

func (m Model) next(msg tea.Msg) (Model, tea.Cmd) {
	switch m.state {
	case PLAYLIST_VIEW_1:
		m.playlistlist.BaseComponent.Focused = false
		m.tracklist.BaseComponent.Focused = true
		m.state = PLAYLIST_VIEW_2

	case PLAYLIST_VIEW_2:
		m.tracklist.BaseComponent.Focused = false
		m.previewlist.BaseComponent.Focused = true
		m.state = PLAYLIST_VIEW_3
	}
	return m, nil
}

func (m Model) prev(msg tea.Msg) (Model, tea.Cmd) {
	switch m.state {
	case PLAYLIST_VIEW_3:
		m.previewlist.BaseComponent.Focused = false
		m.tracklist.BaseComponent.Focused = true
		m.state = PLAYLIST_VIEW_2

	case PLAYLIST_VIEW_2:
		m.playlistlist.BaseComponent.Focused = true
		m.tracklist.BaseComponent.Focused = false
		m.state = PLAYLIST_VIEW_1
	}
	return m, nil
}

func (m *Model) loadModels(msg tea.WindowSizeMsg) {
	m.state = PLAYLIST_VIEW_1
	// m.previewlist = previewlist.New(msg)
	// m.playlistlist = playlistlist.New(msg)
	// m.tracklist = tracklist.New(msg)
	m.confirm = textinput.NewModel()
}

func (m *Model) UpdateProgramContext(ctx *context.ProgramContext) {
	m.ctx = ctx
	// not sure this is how I want to do this
	m.SetSize()
}

func (m *Model) SetSize() {
	// m.SetWidth()
	// m.SetHeight()
	m.playlistlist.UpdateProgramContext(m.ctx)
	m.tracklist.UpdateProgramContext(m.ctx)
	m.previewlist.UpdateProgramContext(m.ctx)
}
