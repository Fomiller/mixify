package ui

import (
	"fmt"
	"net/http"
)

func (m Model) View() string {
	switch m.view {

	case "menu":
		return menuView(m)

	case "playlist":
		return playlistView(m)

	case "track":
		return trackView(m)

	case "test":
		return menuView(m)

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

	for i, choice := range m.choices {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if choice.selected {
			checked = "x" // selected!
		}

		choice, ok := choice.detail.(track)
		if ok {
			// for _, track := range choice {
			// Is the cursor pointing at this choice?
			// Render the row
			output += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.name)
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
		if choice.selected {
			checked = "x" // selected!
		}

		// Render the row
		choice, ok := choice.detail.(playlist)
		if ok {
			output += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.Name())
		}
	}

	// The footer
	output += "\nPress q to quit.\n"
	return output
}

func serverView(m Model) string {
	var output string
	// Send the UI for rendering
	output = fmt.Sprintf("Checking %s ... %s", url, m.view)

	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: Q%v\n\n", m.err)
	}

	if m.status > 0 {
		output += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	}
	return output
}

func menuView(m Model) string {
	var output string
	state := m.views[MENU]
	// Send the UI for rendering
	newState := state.(menuModel)
	output = fmt.Sprintf("Checking MENU %s ... %v", url, newState.tag)

	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: Q%v\n\n", m.err)
	}

	if m.status > 0 {
		output += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	}
	return output
}

func testView(m Model) string {
	var output string
	// Send the UI for rendering
	output = fmt.Sprintf("Checking %s ... %s", url, m.view)

	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: Q%v\n\n", m.err)
	}

	if m.status > 0 {
		output += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	}
	return output
}
