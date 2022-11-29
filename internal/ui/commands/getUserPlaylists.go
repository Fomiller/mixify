package commands

import (
	"context"
	"log"

	"github.com/Fomiller/mixify/internal/auth"
	"github.com/Fomiller/mixify/internal/ui/messages"
	tea "github.com/charmbracelet/bubbletea"
)

func GetUserPlaylistsCmd() tea.Msg {
	playlist, err := auth.Client.CurrentUsersPlaylists(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return messages.PlaylistMsg(playlist)
}
