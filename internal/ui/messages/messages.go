package messages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zmb3/spotify/v2"
)

// need to implement these down the road
type StatusMsg int
type ErrMsg error
type ResetStateMsg int

type BackMsg bool
type ExitInputMsg bool

type CreatePlaylistMsg bool

// I think these are not needed
type NextMsg bool
type PrevMsg bool

type UserMsg *spotify.PrivateUser
type PlaylistMsg *spotify.SimplePlaylistPage

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func CreatePlaylistCmd() tea.Msg {
	return CreatePlaylistMsg(true)
}
