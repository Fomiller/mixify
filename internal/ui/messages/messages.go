package messages

import (
	"github.com/zmb3/spotify/v2"
)

// need to implement these down the road
type InitMsg map[string]string

type StatusMsg int
type ErrMsg error
type ResetStateMsg int

type BackMsg bool
type ExitInputMsg bool

type CreatePlaylistMsg struct {
	PlaylistName string
	Description  string
}

// I think these are not needed
type NextMsg bool
type PrevMsg bool

type UserMsg *spotify.PrivateUser
type PlaylistMsg *spotify.SimplePlaylistPage

type CreatePlaylistSuccessMsg bool
type CreatePlaylistErrorMsg struct{ Err error }
