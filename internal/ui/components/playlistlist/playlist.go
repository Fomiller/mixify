package playlistlist

import (
	"context"
	"log"

	"github.com/Fomiller/mixify/internal/auth"
	"github.com/Fomiller/mixify/internal/ui/components/base"
	"github.com/Fomiller/mixify/internal/ui/components/playlist"
	"github.com/Fomiller/mixify/internal/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	emoji "github.com/tmdvs/Go-Emoji-Utils"
)

// TODO make this a command?
func GetUserPlaylists() list.Model {
	var items []list.Item

	spotifyUserPlaylists, err := auth.Client.CurrentUsersPlaylists(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range spotifyUserPlaylists.Playlists {
		item := playlist.Playlist{
			PlaylistTitle: emoji.RemoveAll(p.Name),
			Desc:          p.Description,
			Base: base.Item{
				Selected: false,
			},
			Playlist: p,
		}
		items = append(items, item)
	}

	// TODO make this height and width dynamic for now it works
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle.Foreground(styles.Selected)
	delegate.Styles.NormalTitle.Foreground(styles.Unselected)

	newList := list.New(items, delegate, 60, 50)
	newList.KeyMap.NextPage = key.NewBinding(key.WithKeys("pgdown", "J"))
	newList.KeyMap.PrevPage = key.NewBinding(key.WithKeys("pgup", "K"))
	return newList
}
