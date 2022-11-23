package track

import (
	"github.com/Fomiller/mixify/internal/ui/components/base"
	"github.com/Fomiller/mixify/internal/ui/styles"
	"github.com/zmb3/spotify/v2"
)

type Track struct {
	Base       base.Item
	TrackTitle string
	TrackID    spotify.ID
	PlaylistID spotify.ID
	Desc       string
}

func (t Track) Description() string { return t.Desc }
func (t Track) FilterValue() string { return t.TrackTitle }
func (t Track) Title() string {
	if t.Base.Selected == true {
		return styles.SelectedItemStyle.Render(t.TrackTitle)
	} else {
		return t.TrackTitle
	}
}

func (t *Track) ToggleSelected() {
	t.Base.Selected = !t.Base.Selected
}
