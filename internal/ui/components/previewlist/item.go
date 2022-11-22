package previewlist

import "github.com/zmb3/spotify/v2"

type Item struct {
	title      string
	desc       string
	TrackID    spotify.ID
	PlaylistID spotify.ID
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.title }
