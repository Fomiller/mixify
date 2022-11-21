package styles

import "github.com/charmbracelet/lipgloss"

var (
	AppStyle = lipgloss.NewStyle().Padding(1, 2)

	MenuStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			BorderStyle(lipgloss.NormalBorder())

	DisplayStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Background(lipgloss.Color("FFF")).
			BorderStyle(lipgloss.NormalBorder())

	ViewStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Background(lipgloss.Color("FFF")).
			BorderStyle(lipgloss.NormalBorder())

	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	StatusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render

	DocStyle = lipgloss.NewStyle().Margin(1, 2).Border(lipgloss.RoundedBorder())
)
