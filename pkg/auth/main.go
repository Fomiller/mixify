package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify"
)

const redirectURL = "http://localhost:42069/callback/"

var (
	auth = spotify.NewAuthenticator(redirectURL,
		spotify.ScopeUserReadCurrentlyPlaying,
		spotify.ScopeUserReadPlaybackState,
		spotify.ScopeUserModifyPlaybackState,
		spotify.ScopeUserLibraryModify,
		spotify.ScopeUserLibraryRead,
		spotify.ScopePlaylistModifyPublic,
	)

	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func init() {
	// load env vars
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	auth.SetAuthInfo(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	ch <- &client
}
