package combined

import (
	"context"
	"fmt"

	"github.com/Fomiller/mixify/pkg/auth"
	"github.com/Fomiller/mixify/pkg/ui/models"
	"github.com/Fomiller/mixify/pkg/ui/models/playlist/track"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zmb3/spotify/v2"
)

type view string

type Model struct {
	state   view
	Focused bool
	List    list.Model
	cursor  int
	status  int
	err     error
	name    string
}

type Item struct {
	title      string
	desc       string
	TrackID    spotify.ID
	PlaylistID spotify.ID
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.title }

func New() Model {
	items := []list.Item{}
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle.Foreground(lipgloss.AdaptiveColor{Light: "#1DB954", Dark: "#1DB954"})
	delegate.Styles.NormalTitle.Foreground(lipgloss.AdaptiveColor{Light: "#1DB925", Dark: "#1DB925"})
	combinedList := list.New(items, delegate, 60, 50)
	combinedList.KeyMap.NextPage = key.NewBinding(
		key.WithKeys("pgdown", "J"),
	)
	combinedList.KeyMap.PrevPage = key.NewBinding(
		key.WithKeys("pgup", "K"),
	)
	return Model{Focused: false, List: combinedList}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case models.StatusMsg:
		m.status = int(msg)
		return m, nil

	case models.ErrMsg:
		m.err = msg
		return m, tea.Quit

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)

	// Is it a key press?
	case tea.KeyMsg:
		switch msg.String() {
		// return to previous view with backspace
		case tea.KeyBackspace.String():
			return m, func() tea.Msg {
				return models.BackMsg(true)
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
	switch m.Focused {
	case true:
		return focusedStyle.Render(m.List.View())
	default:
		return docStyle.Render(m.List.View())
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) CreatePlaylist() error {
	user, err := auth.Client.CurrentUser(context.Background())
	if err != nil {
		return err
	}
	newPlaylist, err := auth.Client.CreatePlaylistForUser(context.Background(), user.ID, "my-test-playlist", "created with mixify", true, false)
	if err != nil {
		return err
	}
	var trackIDs []spotify.ID
	tracks := m.List.Items()
	for _, t := range tracks {
		x := t.(track.Item)
		trackIDs = append(trackIDs, x.TrackID)
	}
	snapShotID, err := auth.Client.AddTracksToPlaylist(context.Background(), newPlaylist.ID, trackIDs...)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully created new playlist %v", snapShotID)
	return nil
}
