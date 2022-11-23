package previewlist

import (
	"context"
	"log"

	"github.com/Fomiller/mixify/internal/auth"
	"github.com/Fomiller/mixify/internal/ui/components/track"
	"github.com/zmb3/spotify/v2"
)

func (m Model) CreatePlaylist(name string, desc string) error {
	var trackIDs []spotify.ID
	log.Println("description: ", desc)
	if desc == "" {
		desc = "Created with mixify"
	}
	if name == "" {
		log.Fatal("No Name in createplaylist")
	}
	log.Println("description: ", desc)

	user, err := auth.Client.CurrentUser(context.Background())
	if err != nil {
		return err
	}

	newPlaylist, err := auth.Client.CreatePlaylistForUser(context.Background(), user.ID, name, desc, true, false)
	if err != nil {
		return err
	}

	tracks := m.List.Items()
	for _, t := range tracks {
		x := t.(track.Track)
		trackIDs = append(trackIDs, x.TrackID)
	}

	// make multiple calls to add tracks if needed, spotify only supports 100 at a time
	if len(trackIDs) > 100 {
		tracks := chunkIDs(trackIDs, 100)
		for _, t := range tracks {
			_, err := auth.Client.AddTracksToPlaylist(context.Background(), newPlaylist.ID, t...)
			if err != nil {
				return err
			}
		}
	} else {
		_, err := auth.Client.AddTracksToPlaylist(context.Background(), newPlaylist.ID, trackIDs...)
		if err != nil {
			return err
		}
	}
	return nil
}

func chunkIDs(slice []spotify.ID, chunkSize int) [][]spotify.ID {
	var chunks [][]spotify.ID
	for {
		if len(slice) == 0 {
			break
		}
		// necessary check to avoid slicing beyond
		// slice capacity
		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}

		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}
	return chunks
}
