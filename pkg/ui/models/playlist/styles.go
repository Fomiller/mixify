package playlist

import "github.com/charmbracelet/lipgloss"

var focusedStyle = lipgloss.NewStyle().
	Margin(1, 2).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("62"))
