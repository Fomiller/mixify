package previewlist

import (
	"log"

	"github.com/Fomiller/mixify/internal/ui/components/base"
	"github.com/Fomiller/mixify/internal/ui/components/textinput"
	"github.com/Fomiller/mixify/internal/ui/context"
	"github.com/Fomiller/mixify/internal/ui/messages"
	"github.com/Fomiller/mixify/internal/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type view string

type Model struct {
	base.Component
	Base    base.List
	Confirm bool
	state   view
	List    list.Model
	name    string
}

func New(msg context.ProgramContext) Model {
	items := []list.Item{}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle.Foreground(styles.Selected)
	delegate.Styles.NormalTitle.Foreground(styles.Unselected)

	list := list.New(items, delegate, 60, 50)
	list.KeyMap.NextPage = key.NewBinding(key.WithKeys("pgdown", "J"))
	list.KeyMap.PrevPage = key.NewBinding(key.WithKeys("pgup", "K"))

	return Model{
		Base: base.List{
			Focused: false,
			Width:   msg.ScreenWidth,
			Height:  msg.ScreenHeight,
		},
		List: list,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case messages.StatusMsg:
		m.Base.Status = int(msg)
		return m, nil

	case messages.ErrMsg:
		m.Base.Err = msg
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
	h, _ := styles.DocStyle.GetFrameSize()

	if m.Confirm == true {
		input := textinput.New()
		return styles.DocStyle.Render(input.View())
	}

	switch m.Base.Focused {
	case true:
		log.Println("COMBINED WIDTH: ", m.Base.Width)
		return styles.FocusedStyle.Width((m.Base.Width / 3) - h).Render(m.List.View())
	default:
		return styles.DocStyle.Width((m.Base.Width / 3) - h).Render(m.List.View())
	}
}

func (m *Model) SetWidth(width int) {
	m.Base.Width = width
}

func (m *Model) SetHeight(height int) {
	m.Base.Height = height
}
