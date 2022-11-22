package track

import (
	"github.com/Fomiller/mixify/internal/ui/components/basecomponents"
	"github.com/Fomiller/mixify/internal/ui/styles"
	"github.com/zmb3/spotify/v2"
)

type Track struct {
	Item       basecomponents.Item
	TrackTitle string
	TrackID    spotify.ID
	PlaylistID spotify.ID
	Desc       string
}

func (t Track) Description() string { return t.Desc }
func (t Track) FilterValue() string { return t.TrackTitle }
func (t Track) Title() string {
	if t.Item.Selected == true {
		return styles.SelectedItemStyle.Render(t.TrackTitle)
	} else {
		return t.TrackTitle
	}
}

func (t *Track) ToggleSelected() {
	t.Item.Selected = !t.Item.Selected
}
