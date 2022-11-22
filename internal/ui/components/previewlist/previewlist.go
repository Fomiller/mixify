package previewlist

import (
	"context"
	"log"

	"github.com/Fomiller/mixify/internal/auth"
	"github.com/Fomiller/mixify/internal/ui/components/textinput"
	"github.com/Fomiller/mixify/internal/ui/components/track"
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
	Confirm bool
	state   view
	Focused bool
	List    list.Model
	cursor  int
	status  int
	err     error
	name    string
	Width   int
	Height  int
}

func New(msg tea.WindowSizeMsg) Model {
	items := []list.Item{}
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle.Foreground(lipgloss.AdaptiveColor{Light: "#1DB954", Dark: "#1DB954"})
	delegate.Styles.NormalTitle.Foreground(lipgloss.AdaptiveColor{Light: "#1DB925", Dark: "#1DB925"})
	styles.DocStyle.GetFrameSize()
	list := list.New(items, delegate, 60, 50)

	// h, v := docStyle.GetFrameSize()
	// list.SetSize(msg.Width/divisor, msg.Height-v)

	list.KeyMap.NextPage = key.NewBinding(
		key.WithKeys("pgdown", "J"),
	)
	list.KeyMap.PrevPage = key.NewBinding(
		key.WithKeys("pgup", "K"),
	)
	return Model{
		Focused: false,
		List:    list,
		Width:   msg.Width,
		Height:  msg.Height,
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case messages.StatusMsg:
		m.status = int(msg)
		return m, nil

	case messages.ErrMsg:
		m.err = msg
		return m, tea.Quit

	case tea.WindowSizeMsg:
		// h, v := docStyle.GetFrameSize()
		// m.List.SetSize(msg.Width-h, msg.Height-v)

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
	if m.Confirm == true {
		cm := textinput.New()
		return styles.DocStyle.Render(cm.View())
	}
	h, _ := styles.DocStyle.GetFrameSize()
	switch m.Focused {
	case true:
		log.Println("COMBINED WIDTH: ", m.Width)
		return styles.FocusedStyle.Width((m.Width / 3) - h).Render(m.List.View())
	default:
		return styles.DocStyle.Width((m.Width / 3) - h).Render(m.List.View())
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) CreatePlaylist(name string, desc string) error {
	user, err := auth.Client.CurrentUser(context.Background())
	if err != nil {
		return err
	}
	if desc == "" {
		desc = "Created with mixify"
	}
	newPlaylist, err := auth.Client.CreatePlaylistForUser(context.Background(), user.ID, name, desc, true, false)
	if err != nil {
		return err
	}
	var trackIDs []spotify.ID
	tracks := m.List.Items()
	for _, t := range tracks {
		x := t.(track.Track)
		trackIDs = append(trackIDs, x.TrackID)
	}
	// make multiple calls to add tracks if needed, spotify only supports 100 at a time
	if len(trackIDs) > 100 {
		tracks := chunkIDs(trackIDs, 100)
		for _, t := range tracks {
			_, err := auth.Client.AddTracksToPlaylist(context.Background(), newPlaylist.ID, t...)
			if err != nil {
				return err
			}
		}
	} else {
		_, err := auth.Client.AddTracksToPlaylist(context.Background(), newPlaylist.ID, trackIDs...)
		if err != nil {
			return err
		}
	}
	return nil
}

func chunkIDs(slice []spotify.ID, chunkSize int) [][]spotify.ID {
	var chunks [][]spotify.ID
	for {
		if len(slice) == 0 {
			break
		}

		// necessary check to avoid slicing beyond
		// slice capacity
		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}

		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}

	return chunks
}

func (m *Model) SetWidth(width int) {
	m.Width = width
}

func (m *Model) SetHeight(height int) {
	m.Height = height
}
