package playlistSelect

import "github.com/charmbracelet/lipgloss"

var docStyle = lipgloss.NewStyle().
	Margin(1, 2).
	Border(lipgloss.RoundedBorder())

var focusedStyle = lipgloss.NewStyle().
	Margin(1, 2).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("62"))

var selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("62"))
