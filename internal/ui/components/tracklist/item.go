package tracklist

import (
	"github.com/Fomiller/mixify/internal/ui/styles"
	"github.com/zmb3/spotify/v2"
)

type Item struct {
	ItemTitle  string
	Desc       string
	TrackID    spotify.ID
	PlaylistID spotify.ID
	Selected   bool
}

func (i Item) Title() string {
	if i.Selected == true {
		return styles.SelectedItemStyle.Render(i.ItemTitle)
	} else {
		return i.ItemTitle
	}
}
func (i Item) Description() string { return i.Desc }
func (i Item) FilterValue() string { return i.ItemTitle }

func (i *Item) ToggleSelected() {
	i.Selected = !i.Selected
}
