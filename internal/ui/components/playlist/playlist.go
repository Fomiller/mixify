package playlist

import (
	"github.com/Fomiller/mixify/internal/ui/components/base"
	"github.com/Fomiller/mixify/internal/ui/styles"
	"github.com/zmb3/spotify/v2"
)

type Playlist struct {
	Base          base.Item
	Playlist      spotify.SimplePlaylist
	ID            spotify.ID
	PlaylistTitle string
	Desc          string
}

func (p *Playlist) ToggleSelected() {
	p.Base.Selected = !p.Base.Selected
}

func (p Playlist) Title() string {
	if p.Base.Selected == true {
		return styles.SelectedItemStyle.Render(p.PlaylistTitle)
	} else {
		return p.PlaylistTitle
	}
}
func (p Playlist) Description() string { return p.Desc }
func (p Playlist) FilterValue() string { return p.PlaylistTitle }
