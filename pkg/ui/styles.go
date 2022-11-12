package ui

import "github.com/charmbracelet/lipgloss"

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	menuStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			BorderStyle(lipgloss.NormalBorder())

	displayStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Background(lipgloss.Color("FFF")).
			BorderStyle(lipgloss.NormalBorder())

	viewStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Background(lipgloss.Color("FFF")).
			BorderStyle(lipgloss.NormalBorder())

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render

	docStyle = lipgloss.NewStyle().Margin(1, 2).Border(lipgloss.RoundedBorder())
)
