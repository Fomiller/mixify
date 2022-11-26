package commands

import (
	"context"
	"log"

	"github.com/Fomiller/mixify/internal/auth"
	"github.com/Fomiller/mixify/internal/ui/components/track"
	"github.com/Fomiller/mixify/internal/ui/messages"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zmb3/spotify/v2"
)

const (
	defaultDescription = "Created with Mixify"
)

func CreatePlaylistCmd(name string, desc string, tracks []list.Item) tea.Cmd {
	return func() tea.Msg {
		if desc == "" {
			desc = defaultDescription
		}
		if name == "" {
			log.Fatal("No Name provided to createplaylist")
		}

		user, err := auth.Client.CurrentUser(context.Background())
		if err != nil {
			return messages.CreatePlaylistErrorMsg{Err: err}
		}

		newPlaylist, err := auth.Client.CreatePlaylistForUser(context.Background(), user.ID, name, desc, true, false)
		if err != nil {
			return messages.CreatePlaylistErrorMsg{Err: err}
		}

		err = addTracksToPlaylist(newPlaylist, tracks)
		if err != nil {
			return messages.CreatePlaylistErrorMsg{Err: err}
		}
		return messages.CreatePlaylistSuccessMsg(true)
	}
}

func addTracksToPlaylist(playlist *spotify.FullPlaylist, tracks []list.Item) error {
	var trackIDs []spotify.ID

	for _, t := range tracks {
		x := t.(track.Track)
		trackIDs = append(trackIDs, x.TrackID)
	}

	// make multiple calls to add tracks if needed, spotify only supports 100 at a time
	if len(trackIDs) > 100 {
		tracks := chunkIDs(trackIDs, 100)
		for _, t := range tracks {
			_, err := auth.Client.AddTracksToPlaylist(context.Background(), playlist.ID, t...)
			if err != nil {
				return err
			}
		}
	} else {
		_, err := auth.Client.AddTracksToPlaylist(context.Background(), playlist.ID, trackIDs...)
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
