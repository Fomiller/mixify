package track

import (
	"github.com/Fomiller/mixify/pkg/ui/models"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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
	Name    string
}

type item struct {
	title string
	desc  string
	ID    spotify.ID
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func New() Model {
	items := []list.Item{}
	trackList := list.New(items, list.NewDefaultDelegate(), 60, 50)
	trackList.KeyMap.NextPage = key.NewBinding(
		key.WithKeys("pgdown", "J"),
	)
	trackList.KeyMap.PrevPage = key.NewBinding(
		key.WithKeys("pgup", "K"),
	)

	return Model{Focused: false, List: trackList}
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
		// h, v := docStyle.GetFrameSize()
		m.List.SetSize(msg.Width/3, msg.Height)

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
