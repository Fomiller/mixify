package ui

import (
	"fmt"
	"net/http"
)

func (m Model) View() string {
	switch m.view {

	case "playlistView":
		return playlistView(m)

	case "trackView":
		return trackView(m)

	case "choiceView":
		return choiceView(m)

	default:
		return choiceView(m)
	}
}

func playlistView(m Model) string {
	var output string
	return output
}

func trackView(m Model) string {
	var output string

	for i, playlist := range m.selected {

		for _, track := range playlist.tracks {
			// Is the cursor pointing at this choice?
			cursor := " " // no cursor
			if m.cursor == i {
				cursor = ">" // cursor!
			}

			// Is this choice selected?
			checked := " " // not selected
			if _, ok := m.selected[i]; ok {
				checked = "x" // selected!
			}

			// Render the row
			output += fmt.Sprintf("%s [%s] %s\n", cursor, checked, track)
		}
	}
	return output
}

func choiceView(m Model) string {
	var output string
	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		output += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.name)
	}

	// The footer
	output += "\nPress q to quit.\n"
	return output
}

func serverView(m Model) string {
	var output string
	// Send the UI for rendering
	output = fmt.Sprintf("Checking %s ... ", url)

	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: Q%v\n\n", m.err)
	}

	if m.status > 0 {
		output += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	}
	return output
}
