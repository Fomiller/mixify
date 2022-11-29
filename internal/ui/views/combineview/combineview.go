package combineview

import (
	"github.com/Fomiller/mixify/internal/ui/commands"
	"github.com/Fomiller/mixify/internal/ui/components/base"
	"github.com/Fomiller/mixify/internal/ui/components/playlist"
	"github.com/Fomiller/mixify/internal/ui/components/playlistlist"
	"github.com/Fomiller/mixify/internal/ui/components/previewlist"
	"github.com/Fomiller/mixify/internal/ui/components/textinput"
	"github.com/Fomiller/mixify/internal/ui/components/track"
	"github.com/Fomiller/mixify/internal/ui/components/tracklist"
	"github.com/Fomiller/mixify/internal/ui/context"
	"github.com/Fomiller/mixify/internal/ui/keys"
	"github.com/Fomiller/mixify/internal/ui/messages"
	"github.com/Fomiller/mixify/internal/ui/views"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type section string

const (
	PlaylistSection section = "PlaylistSection"
	TrackSection    section = "TrackSection"
	PreviewSection  section = "PreviewSection"
	ConfirmSection  section = "ConfirmSection"
)

type Model struct {
	BaseComponent base.List
	ctx           *context.ProgramContext

	playlistlist playlistlist.Model
	tracklist    tracklist.Model
	previewlist  previewlist.Model
	confirm      textinput.Model

	loaded         bool
	CurrentSection section
	Width          int
	Height         int
}

func NewModel(ctx context.ProgramContext) Model {
	m := Model{
		ctx:            &ctx,
		CurrentSection: PlaylistSection,
		loaded:         false,
		Width:          ctx.ScreenWidth,
		Height:         ctx.ScreenHeight,
		previewlist:    previewlist.NewModel(ctx),
		playlistlist:   playlistlist.NewModel(ctx),
		tracklist:      tracklist.NewModel(ctx),
	}
	m.confirm = textinput.NewModel()

	return m
}

func (m Model) ResetModel(ctx *context.ProgramContext) Model {
	newModel := Model{
		ctx:            m.ctx,
		CurrentSection: PlaylistSection,
		previewlist:    previewlist.NewModel(*m.ctx),
		playlistlist:   playlistlist.NewModel(*m.ctx),
		tracklist:      tracklist.NewModel(*m.ctx),
	}
	newModel.confirm = textinput.NewModel()
	return newModel
}

func (m Model) Init() tea.Cmd {
	return commands.GetUserPlaylistsCmd
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case messages.CreatePlaylistSuccessMsg:
		m = m.ResetModel(m.ctx)
		return m, nil

	case messages.CreatePlaylistErrorMsg:
		m.BaseComponent.Err = msg.Err
		return m, tea.Quit

	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch {
		case key.Matches(msg, keys.Keys.Escape):
			if m.ctx.View != views.MainMenuView {
				m.ctx.View = views.MainMenuView
				return m, nil
			}
		}

		switch msg.String() {

		//TODO move this logic into the ui model.
		case "right", "l":
			return m.next(msg)

		//TODO move this logic into the ui model.
		case "left", "h":
			return m.prev(msg)

		case "enter", " ":
			switch m.CurrentSection {
			case PlaylistSection:
				item := m.playlistlist.List.SelectedItem().(playlist.Playlist)
				cursor := m.playlistlist.List.Index()

				if item.BaseComponent.Selected == false {
					item.ToggleSelected()
					m.playlistlist.List.SetItem(cursor, item)
					m.tracklist.InsertTracks(item.Playlist)
					selectedTracks := m.tracklist.GetSelectedTracks()
					m.previewlist.List.SetItems(selectedTracks)
				} else {
					item.ToggleSelected()
					m.playlistlist.List.SetItem(cursor, item)
					m.tracklist.RemoveTracks(item.Playlist.ID)
					selectedTracks := m.tracklist.GetSelectedTracks()
					m.previewlist.List.SetItems(selectedTracks)

				}
				return m, nil

			case TrackSection:
				item := m.tracklist.List.SelectedItem().(track.Track)
				cursor := m.tracklist.List.Index()

				item.ToggleSelected()
				m.tracklist.List.SetItem(cursor, item)
				selectedTracks := m.tracklist.GetSelectedTracks()
				m.previewlist.List.SetItems(selectedTracks)
				return m, nil

			case PreviewSection:
				m.CurrentSection = ConfirmSection
				m.confirm.Tracks = &m.previewlist.List
				return m, nil
			}
		}
	}

	// update nested models based off of state
	switch m.CurrentSection {
	case PlaylistSection:
		m.playlistlist, cmd = m.playlistlist.Update(msg)
		cmds = append(cmds, cmd)

	case TrackSection:
		m.tracklist, cmd = m.tracklist.Update(msg)
		cmds = append(cmds, cmd)

	case PreviewSection:
		m.previewlist, cmd = m.previewlist.Update(msg)
		cmds = append(cmds, cmd)

	case ConfirmSection:
		m.confirm, cmd = m.confirm.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// Main Model view
func (m Model) View() string {
	var output string

	if m.CurrentSection == ConfirmSection {
		output = m.confirm.View()
	} else {
		output = lipgloss.JoinHorizontal(lipgloss.Center, m.playlistlist.View(), m.tracklist.View(), m.previewlist.View())
	}
	return output
}

func (m Model) next(msg tea.Msg) (Model, tea.Cmd) {
	switch m.CurrentSection {
	case PlaylistSection:
		m.playlistlist.BaseComponent.Focused = false
		m.tracklist.BaseComponent.Focused = true
		m.CurrentSection = TrackSection

	case TrackSection:
		m.tracklist.BaseComponent.Focused = false
		m.previewlist.BaseComponent.Focused = true
		m.CurrentSection = PreviewSection
	}
	return m, nil
}

func (m Model) prev(msg tea.Msg) (Model, tea.Cmd) {
	switch m.CurrentSection {
	case PreviewSection:
		m.previewlist.BaseComponent.Focused = false
		m.tracklist.BaseComponent.Focused = true
		m.CurrentSection = TrackSection

	case TrackSection:
		m.playlistlist.BaseComponent.Focused = true
		m.tracklist.BaseComponent.Focused = false
		m.CurrentSection = PlaylistSection
	}
	return m, nil
}

func (m *Model) UpdateProgramContext(ctx *context.ProgramContext) {
	m.ctx = ctx
	m.playlistlist.UpdateProgramContext(m.ctx)
	m.tracklist.UpdateProgramContext(m.ctx)
	m.previewlist.UpdateProgramContext(m.ctx)
}
