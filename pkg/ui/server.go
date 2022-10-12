package ui

import (
	"fmt"
	"net/http"

	"github.com/Fomiller/mixify/pkg/ui/models/playlist/track"
)

func serverView(m mainModel) string {
	var output string
	state := m.views[TRACK].(track.Model)

	// Send the UI for rendering
	output = fmt.Sprintf("Checking %s ... %s", url, m.view)

	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: Q%v\n\n", m.err)
	}

	if m.status > 0 {
		output += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	}
	output += fmt.Sprintf("STATE: %s\n", state.Name)
	return output
}
