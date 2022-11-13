package models

import tea "github.com/charmbracelet/bubbletea"

type statusMsg int

type errMsg struct{ err error }

func CreatePlaylistCmd() tea.Msg {
	return CreatePlaylistMsg(true)
}

func (e errMsg) Error() string { return e.err.Error() }

// func checkServer() tea.Msg {
// 	c := &http.Client{Timeout: 10 * time.Second}
// 	res, err := c.Get(url)

// 	if err != nil {
// 		return errMsg{err}
// 	}
// 	return statusMsg(res.StatusCode)
// }

// func (m Model) CreatePlaylistCmd(name string, desc string) tea.Msg {
// 	c := playlist.Model.combined.(combined.Model)
// 	err := c.CreatePlaylist(name, desc)
// 	if err != nil {
// 		return models.ErrMsg(err)
// 	}
// 	return models.ResetStateMsg(1)
// }

// func (m Model) CreatePlaylistCmd(name string, desc string) error {
// 	user, err := auth.Client.CurrentUser(context.Background())
// 	if err != nil {
// 		return err
// 	}
// 	if desc == "" {
// 		desc = "Created with mixify"
// 	}
// 	newPlaylist, err := auth.Client.CreatePlaylistForUser(context.Background(), user.ID, name, desc, true, false)
// 	if err != nil {
// 		return err
// 	}
// 	var trackIDs []spotify.ID
// 	tracks := m.combined.List.Items()
// 	for _, t := range tracks {
// 		x := t.(track.Item)
// 		trackIDs = append(trackIDs, x.TrackID)
// 	}
// 	snapShotID, err := auth.Client.AddTracksToPlaylist(context.Background(), newPlaylist.ID, trackIDs...)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Printf("Successfully created new playlist %v", snapShotID)
// 	return models.ResetStateMsg
// }
