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
	"github.com/joho/godotenv"
	"github.com/zmb3/spotify"
)

const redirectURL = "http://localhost:42069/callback/"

var (
	ch     = make(chan *spotify.Client)
	client *spotify.Client
)

func main() {
	// load env vars
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// // spotify setup
	auth.Auth.SetAuthInfo(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))

	// // http server setup
	http.HandleFunc("/callback/", auth.CompleteAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})

	go http.ListenAndServe(":42069", nil)

	url := auth.Auth.AuthURL(auth.State)
	go fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

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
