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

func Init() {
	// load env vars
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// // spotify setup
	Auth.SetAuthInfo(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))
}

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
