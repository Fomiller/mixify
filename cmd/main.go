package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Fomiller/mixify/pkg/auth"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const redirectURL = "http://localhost:42069/callback/"

var (
	// client *spotify.Client
	Auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURL),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadCurrentlyPlaying,
			spotifyauth.ScopeUserReadPlaybackState,
			spotifyauth.ScopeUserModifyPlaybackState,
			spotifyauth.ScopeUserLibraryModify,
			spotifyauth.ScopeUserLibraryRead,
			spotifyauth.ScopePlaylistModifyPublic,
			spotifyauth.ScopePlaylistReadPrivate,
		),
	)
)

func main() {
	// // // http server setup
	http.HandleFunc("/callback/", auth.CompleteAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":42069", nil)

	url := auth.Auth.AuthURL(auth.State)
	fmt.Printf("Please log in to Spotify by visiting the following page in your browser: %s\n", url)

	client := <-auth.Ch

	// // use the client to make calls that require authorization
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("You are logged in as:", user.ID)

	// _, playlist, err := client.FeaturedPlaylists()
	playlist, err := client.CurrentUsersPlaylists(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	p := playlist.Playlists[0]
	fmt.Println("Playlist: ", p)
	fmt.Println("ID: ", p.ID)

	// tui setup
	// rand.Seed(time.Now().UTC().UnixNano())

	// if err := tea.NewProgram(ui.New()).Start(); err != nil {
	// 	fmt.Println("Error running program:", err)
	// 	os.Exit(1)
	// }
}
