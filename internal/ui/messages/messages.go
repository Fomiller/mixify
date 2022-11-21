package messages

import "github.com/zmb3/spotify/v2"

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
