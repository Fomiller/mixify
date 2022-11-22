package playlistlist

import (
	"github.com/Fomiller/mixify/internal/ui/styles"
	"github.com/zmb3/spotify/v2"
)

type Item struct {
	title    string
	desc     string
	ID       spotify.ID
	Playlist spotify.SimplePlaylist
	Selected bool
}

func (i *Item) ToggleSelected() {
	i.Selected = !i.Selected
}

func (i Item) Title() string {
	if i.Selected == true {
		return styles.SelectedItemStyle.Render(i.title)
	} else {
		return i.title
	}
}
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.title }
