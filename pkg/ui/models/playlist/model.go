package playlist

import (
	"github.com/Fomiller/mixify/pkg/ui/models"
	"github.com/Fomiller/mixify/pkg/ui/models/playlist/combined"
	playlistSelect "github.com/Fomiller/mixify/pkg/ui/models/playlist/select"
	"github.com/Fomiller/mixify/pkg/ui/models/playlist/track"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zmb3/spotify/v2"
)

const (
	PLAYLIST_VIEW_1 view = "VIEW_1"
	PLAYLIST_VIEW_2 view = "VIEW_2"
	PLAYLIST_VIEW_3 view = "VIEW_3"
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

	combined       tea.Model
	playlistSelect tea.Model
	track          tea.Model
}

type item struct {
	title string
	desc  string
	ID    spotify.ID
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func New() tea.Model {
	m := Model{
		state:          PLAYLIST_VIEW_1,
		combined:       combined.New(),
		playlistSelect: playlistSelect.New(),
		track:          track.New(),
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return playlistSelect.GetUserPlaylistsCmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// update nested models based off of state
	switch m.state {
	case PLAYLIST_VIEW_1:
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
	}

	switch msg := msg.(type) {

	case models.StatusMsg:
		m.status = int(msg)
		return m, nil

	case models.ErrMsg:
		m.err = msg
		return m, tea.Quit

	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {
		// return to previous view with backspace
		case tea.KeyBackspace.String():
			return m, func() tea.Msg {
				return models.BackMsg(true)
			}

		// These keys should exit the program.
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
				selectModel, _ := m.playlistSelect.(playlistSelect.Model)
				trackModel := m.track.(track.Model)

				item := selectModel.List.SelectedItem().(playlistSelect.Item)
				cursor := selectModel.List.Cursor()

				if item.Selected == false {
					item.ToggleSelected()
					selectModel.List.SetItem(cursor, item)
					trackModel = trackModel.InsertTracks(item.Playlist)
					m.playlistSelect = selectModel
					m.track = trackModel
					return m, nil

				} else {
					item.ToggleSelected()
					selectModel.List.SetItem(cursor, item)
					trackModel = trackModel.RemoveTracks(item.Playlist.ID)
					m.playlistSelect = selectModel
					m.track = trackModel
					return m, nil
				}

				// case PLAYLIST_VIEW_2:
				// trackModel, _ := m.track.(track.Model)
				// i := trackModel.List.SelectedItem().(track.Item)
				// fmt.Print(i)
				// combinedModel := m.combined.(combined.Model)
				// items := append(combinedModel.List.Items(), i)
				// combinedModel.List = list.New(items, list.NewDefaultDelegate(), 0, 0)
				// m.combined = combinedModel
				// return m, nil

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

	output = lipgloss.JoinHorizontal(lipgloss.Top, m.playlistSelect.View(), m.track.View(), m.combined.View())

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
