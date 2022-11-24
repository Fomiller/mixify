package playlist

import (
	"github.com/Fomiller/mixify/internal/ui/components/base"
	"github.com/Fomiller/mixify/internal/ui/styles"
	"github.com/zmb3/spotify/v2"
)

type Playlist struct {
	BaseComponent base.Item
	Playlist      spotify.SimplePlaylist
	ID            spotify.ID
	PlaylistTitle string
	Desc          string
}

func (p *Playlist) ToggleSelected() {
	p.BaseComponent.Selected = !p.BaseComponent.Selected
}

func (p Playlist) Title() string {
	if p.BaseComponent.Selected == true {
		return styles.SelectedItemStyle.Render(p.PlaylistTitle)
	} else {
		return p.PlaylistTitle
	}
}
func (p Playlist) Description() string { return p.Desc }
func (p Playlist) FilterValue() string { return p.PlaylistTitle }
