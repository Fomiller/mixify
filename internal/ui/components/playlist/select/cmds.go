package playlistSelect

import (
	"context"
	"log"

	"github.com/Fomiller/mixify/pkg/auth"
	"github.com/Fomiller/mixify/pkg/ui/models"
	tea "github.com/charmbracelet/bubbletea"
)

func GetUserCmd() tea.Msg {
	//get playlists from spotify
	user, err := auth.Client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return models.UserMsg(user)
}

func GetUserPlaylistsCmd() tea.Msg {
	playlist, err := auth.Client.CurrentUsersPlaylists(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return models.PlaylistMsg(playlist)
}
