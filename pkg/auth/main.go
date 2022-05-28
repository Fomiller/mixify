package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zmb3/spotify"
)

const redirectURL = "http://localhost:42069/callback/"

var (
	Auth = spotify.NewAuthenticator(redirectURL,
		spotify.ScopeUserReadCurrentlyPlaying,
		spotify.ScopeUserReadPlaybackState,
		spotify.ScopeUserModifyPlaybackState,
		spotify.ScopeUserLibraryModify,
		spotify.ScopeUserLibraryRead,
		spotify.ScopePlaylistModifyPublic,
	)

	Ch    = make(chan *spotify.Client)
	State = "abc123"
)

func CompleteAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := Auth.Token(State, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != State {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, State)
	}
	// use the token to get an authenticated client
	client := Auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	Ch <- &client
}
