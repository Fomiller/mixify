package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/Fomiller/mixify/pkg/auth"
	"github.com/Fomiller/mixify/pkg/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zmb3/spotify"
)

var (
	client *spotify.Client
)

func init() {
}

func main() {
	// // http server setup
	http.HandleFunc("/callback/", auth.CompleteAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})

	go http.ListenAndServe(":42069", nil)

	go fmt.Printf("Please log in to Spotify by visiting the following page in your browser: %s\n", auth.Auth.AuthURL(auth.State))

	client := <-auth.Ch

	// // use the client to make calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("You are logged in as:", user.ID)

	_, playlist, err := client.FeaturedPlaylists()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Playlist: ", playlist)

	// // tui setup
	rand.Seed(time.Now().UTC().UnixNano())

	if err := tea.NewProgram(ui.NewModel()).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
