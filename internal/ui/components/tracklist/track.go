package tracklist

import (
	"context"
	"fmt"

	"github.com/Fomiller/mixify/internal/auth"
	"github.com/Fomiller/mixify/internal/ui/components/base"
	"github.com/Fomiller/mixify/internal/ui/components/track"
	"github.com/charmbracelet/bubbles/list"
	"github.com/zmb3/spotify/v2"
)

func (m Model) InsertTracks(playlist spotify.SimplePlaylist) Model {
	offset := 0
	for {
		result, err := auth.Client.GetPlaylistTracks(context.Background(), playlist.ID, spotify.Offset(offset))
		if err != nil {
			panic(err)
		}

		for _, t := range result.Tracks {
			m.List.InsertItem(len(m.List.Items())+1,
				track.Track{
					TrackTitle: t.Track.Name,
					Desc:       fmt.Sprintf("%v:%v", playlist.Name, playlist.ID),
					TrackID:    t.Track.ID,
					PlaylistID: playlist.ID,
					BaseComponent: base.Item{
						Selected: true,
					},
				})
		}

		if result.Next != "" {
			offset = offset + len(result.Tracks)
			continue
		} else {
			break
		}

	}

	return m
}
func (m Model) RemoveTracks(playlistID spotify.ID) Model {
	newList := []list.Item{}
	for _, t := range m.List.Items() {
		track, ok := t.(track.Track)
		if !ok {
			panic("could not assert list.Item to type Item")
		}
		if track.PlaylistID != playlistID {
			newList = append(newList, track)
		}
	}
	m.List.SetItems(newList)
	return m
}

func (m Model) GetSelectedTracks() []list.Item {
	selected := []list.Item{}
	for _, t := range m.List.Items() {
		track, ok := t.(track.Track)
		if !ok {
			panic("could not assert list.Item to type Item")
		}
		if track.BaseComponent.Selected != false {
			selected = append(selected, track)
		}
	}
	m.List.SetItems(selected)
	return m.List.Items()
}
