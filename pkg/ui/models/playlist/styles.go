package playlist

import "github.com/charmbracelet/lipgloss"

var horizontalSize, verticalSize = lipgloss.NewStyle().GetFrameSize()

var focusedStyle = lipgloss.NewStyle().
	Margin(1, 2).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("62"))

var docStyle = lipgloss.NewStyle().
	Margin(1, 2).
	Border(lipgloss.RoundedBorder())
