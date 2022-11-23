package editview

import (
	"github.com/Fomiller/mixify/internal/ui/context"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
}

func New(msg context.ProgramContext) Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}
func (m Model) View() string {
	return ""
}
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}
