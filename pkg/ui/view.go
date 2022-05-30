package ui

import (
	"fmt"
	"net/http"
)

func (m Model) View() string {
	// The header
	s := "What should we buy at the market?\n\n"

	// Iterate over our choices
	if m.state == "choice" || m.state == "" {
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
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
		}

		// The footer
		s += "\nPress q to quit.\n"
	}

	if m.state == "server" {
		// Send the UI for rendering
		s = fmt.Sprintf("Checking %s ... ", url)

		if m.err != nil {
			return fmt.Sprintf("\nWe had some trouble: Q%v\n\n", m.err)
		}

		if m.status > 0 {
			s += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
		}
	}

	return s
}
