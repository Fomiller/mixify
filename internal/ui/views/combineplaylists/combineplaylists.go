package combineplaylists

import (
	"log"

	"github.com/Fomiller/mixify/internal/ui/components/combined"
	"github.com/Fomiller/mixify/internal/ui/components/confirm"
	playlistSelect "github.com/Fomiller/mixify/internal/ui/components/select"
	"github.com/Fomiller/mixify/internal/ui/components/track"
	"github.com/Fomiller/mixify/internal/ui/context"
	"github.com/Fomiller/mixify/internal/ui/messages"
	"github.com/Fomiller/mixify/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zmb3/spotify/v2"
)

const (
	PLAYLIST_VIEW_1 view = "VIEW_1"
	PLAYLIST_VIEW_2 view = "VIEW_2"
	PLAYLIST_VIEW_3 view = "VIEW_3"
	PLAYLIST_VIEW_4 view = "VIEW_4"
)

type view string

type Model struct {
	state   view
	focused view
	cursor  int
	status  int
	err     error
	view    view
	name    string
	ctx     context.ProgramContext

	combined       tea.Model
	playlistSelect tea.Model
	track          tea.Model
	confirm        tea.Model

	loaded bool
	Width  int
	Height int
}

type item struct {
	title string
	desc  string
	ID    spotify.ID
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func New(msg tea.WindowSizeMsg) tea.Model {
	m := Model{
		state:          PLAYLIST_VIEW_1,
		loaded:         false,
		Width:          msg.Width,
		Height:         msg.Height,
		combined:       combined.New(msg),
		playlistSelect: playlistSelect.New(msg),
		track:          track.New(msg),
		confirm:        confirm.New(),
	}

	return m
}

func (m Model) ResetModel() tea.Model {
	return Model{
		combined:       combined.New(tea.WindowSizeMsg{Width: m.Width, Height: m.Height}),
		playlistSelect: playlistSelect.New(tea.WindowSizeMsg{Width: m.Width, Height: m.Height}),
		track:          track.New(tea.WindowSizeMsg{Width: m.Width, Height: m.Height}),
		confirm:        confirm.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return playlistSelect.GetUserPlaylistsCmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// update nested models based off of state
	switch m.state {
	case PLAYLIST_VIEW_1:
		log.Println("bing")
		// return a new updated model and a cmd
		model, newCmd := m.playlistSelect.Update(msg)
		// assert returned interface into struct
		playlistSelectModel, ok := model.(playlistSelect.Model)
		if !ok {
			panic("could not perfom assertion on playlist select model")
		}
		// set cmd to the returned cmd
		cmd = newCmd
		// update the stored model
		m.playlistSelect = playlistSelectModel

	case PLAYLIST_VIEW_2:
		// return a new updated model and a cmd
		model, newCmd := m.track.Update(msg)
		// assert returned interface into struct
		trackSelectModel, ok := model.(track.Model)
		if !ok {
			panic("could not perfom assertion on track select model")
		}
		// set cmd to the returned cmd
		cmd = newCmd
		// update the stored model
		m.track = trackSelectModel

	case PLAYLIST_VIEW_3:
		// return a new updated model and a cmd
		model, newCmd := m.combined.Update(msg)
		// assert returned interface into struct
		combinedPlaylistModel, ok := model.(combined.Model)
		if !ok {
			panic("could not perfom assertion on track select model")
		}
		// set cmd to the returned cmd
		cmd = newCmd
		// update the stored model
		m.combined = combinedPlaylistModel

	case PLAYLIST_VIEW_4:
		// return a new updated model and a cmd
		model, newCmd := m.confirm.Update(msg)
		// assert returned interface into struct
		confirmModel, ok := model.(confirm.Model)
		if !ok {
			panic("could not perfom assertion on confirm input model")
		}
		// set cmd to the returned cmd
		cmd = newCmd
		// update the stored model
		m.confirm = confirmModel
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := styles.DocStyle.GetFrameSize()
		m.setModelSize(msg, h, v)

	case messages.CreatePlaylistMsg:
		promptModel := m.confirm.(confirm.Model)
		combinedModel := m.combined.(combined.Model)
		name := promptModel.Inputs[0].Value()
		desc := promptModel.Inputs[1].Value()
		err := combinedModel.CreatePlaylist(name, desc)
		if err != nil {
			log.Fatal(err)
		}
		m = m.ResetModel().(Model)
		return m, nil

	case messages.ResetStateMsg:
		m.state = PLAYLIST_VIEW_1
		return m, nil

	case messages.StatusMsg:
		m.status = int(msg)
		return m, nil

	case messages.ErrMsg:
		m.err = msg
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
				log.Println("bang")
				selectModel, _ := m.playlistSelect.(playlistSelect.Model)
				trackModel := m.track.(track.Model)
				combinedModel := m.combined.(combined.Model)

				item := selectModel.List.SelectedItem().(playlistSelect.Item)
				cursor := selectModel.List.Index()

				if item.Selected == false {
					item.ToggleSelected()
					selectModel.List.SetItem(cursor, item)
					trackModel = trackModel.InsertTracks(item.Playlist)
					selectedTracks := trackModel.GetSelectedTracks()
					combinedModel.List.SetItems(selectedTracks)

					m.playlistSelect = selectModel
					m.track = trackModel
					m.combined = combinedModel
					return m, nil

				} else {
					item.ToggleSelected()
					selectModel.List.SetItem(cursor, item)
					trackModel = trackModel.RemoveTracks(item.Playlist.ID)
					selectedTracks := trackModel.GetSelectedTracks()
					combinedModel.List.SetItems(selectedTracks)

					m.playlistSelect = selectModel
					m.track = trackModel
					m.combined = combinedModel
					return m, nil
				}

			case PLAYLIST_VIEW_2:
				trackModel := m.track.(track.Model)
				item := trackModel.List.SelectedItem().(track.Item)
				cursor := trackModel.List.Index()

				// if item.Selected == false {
				item.ToggleSelected()
				trackModel.List.SetItem(cursor, item)
				combinedModel := m.combined.(combined.Model)
				selectedTracks := trackModel.GetSelectedTracks()
				combinedModel.List.SetItems(selectedTracks)
				m.track = trackModel
				m.combined = combinedModel
				return m, nil

			case PLAYLIST_VIEW_3:
				m.state = PLAYLIST_VIEW_4
				return m, nil
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, cmd
}

// Main Model view
func (m Model) View() string {
	var output string

	if m.state == PLAYLIST_VIEW_4 {
		output = m.confirm.View()
	} else {
		output = lipgloss.JoinHorizontal(lipgloss.Center, m.playlistSelect.View(), m.track.View(), m.combined.View())
	}
	return output
}

func (m Model) next(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.state == PLAYLIST_VIEW_1 {
		m.state = PLAYLIST_VIEW_2
		// focus track model
		t, ok := m.track.(track.Model)
		if !ok {
			panic("something went wrong")
		}
		t.Focused = true
		m.track = t

		// unfocus playlistselect model
		p, ok := m.playlistSelect.(playlistSelect.Model)
		if !ok {
			panic("something went wrong")
		}
		p.Focused = false
		m.playlistSelect = p

	} else if m.state == PLAYLIST_VIEW_2 {
		m.state = PLAYLIST_VIEW_3
		c, ok := m.combined.(combined.Model)
		if !ok {
			panic("something went wrong")
		}
		c.Focused = true
		m.combined = c

		// unfocus playlistselect model
		t, ok := m.track.(track.Model)
		if !ok {
			panic("some went wrong")
		}
		t.Focused = false
		m.track = t

	} else {
		return m, nil
	}

	return m, nil
}

func (m Model) prev(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.state == PLAYLIST_VIEW_3 {
		m.state = PLAYLIST_VIEW_2

		// focus track model
		t, ok := m.track.(track.Model)
		if !ok {
			panic("something went wrong")
		}
		t.Focused = true
		m.track = t
		// unfocus combiined model
		c, ok := m.combined.(combined.Model)
		if !ok {
			panic("something went wrong")
		}
		c.Focused = false
		m.combined = c

	} else if m.state == PLAYLIST_VIEW_2 {
		m.state = PLAYLIST_VIEW_1
		// focus playlistselect model
		p, ok := m.playlistSelect.(playlistSelect.Model)
		if !ok {
			panic("something went wrong")
		}
		p.Focused = true
		m.playlistSelect = p

		// unfocus track model
		t, ok := m.track.(track.Model)
		if !ok {
			panic("something went wrong")
		}
		t.Focused = false
		m.track = t

	} else {
		return m, cmd
	}

	return m, cmd
}

func (m *Model) setModelSize(msg tea.WindowSizeMsg, h int, v int) {
	divisor := 3
	// log.Println(msg)
	combinedModel := m.combined.(combined.Model)
	combinedModel.SetHeight(msg.Height)
	combinedModel.SetWidth(msg.Width)

	trackModel := m.track.(track.Model)
	trackModel.SetHeight(msg.Height)
	trackModel.SetWidth(msg.Width)

	selectModel := m.playlistSelect.(playlistSelect.Model)
	selectModel.SetHeight(msg.Height)
	selectModel.SetWidth(msg.Width)

	// log.Print(msg)
	// log.Printf("select - w:%v h:%v", selectModel.List.Width(), selectModel.List.Height())
	// log.Printf("Combined - w:%v h:%v", combinedModel.List.Width(), combinedModel.List.Height())
	// log.Printf("Track - w:%v h:%v", trackModel.List.Width(), trackModel.List.Height())
	// log.Printf("------------------------------")
	combinedModel.List.SetSize((msg.Width/divisor)-h, msg.Height-v)
	trackModel.List.SetSize((msg.Width/divisor)-h, msg.Height-v)
	selectModel.List.SetSize((msg.Width/divisor)-h, msg.Height-v)
	// log.Printf("select - w:%v h:%v", selectModel.List.Width(), selectModel.List.Height())
	// log.Printf("Combined - w:%v h:%v", combinedModel.List.Width(), combinedModel.List.Height())
	// log.Printf("Track - w:%v h:%v", trackModel.List.Width(), trackModel.List.Height())
	// log.Printf("------------------------------")
	// log.Printf("------------------------------")

	m.combined = combinedModel
	m.track = trackModel
	m.playlistSelect = selectModel
}

func (m *Model) loadModels(msg tea.WindowSizeMsg) {
	m.state = PLAYLIST_VIEW_1
	m.combined = combined.New(msg)
	m.playlistSelect = playlistSelect.New(msg)
	m.track = track.New(msg)
	m.confirm = confirm.New()
}
