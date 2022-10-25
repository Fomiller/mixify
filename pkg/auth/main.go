package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const redirectURL = "http://localhost:42069/callback/"

var (
	Auth   *spotifyauth.Authenticator
	Ch     = make(chan *spotify.Client)
	State  = "abc123"
	Client *spotify.Client
)

func init() {
	// load env vars
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
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
}

func CompleteAuth(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(r)
	fmt.Println("complete auth")
	tok, err := Auth.Token(r.Context(), State, r)
	fmt.Println(tok)
	tokByte, _ := json.Marshal(tok)
	err = os.WriteFile("./test", tokByte, 0777)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != State {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, State)
	}

	// use the token to get an authenticated client
	client := spotify.New(Auth.Client(r.Context(), tok))
	// fmt.Fprintf(w, "Login Completed!")
	Ch <- client
}
