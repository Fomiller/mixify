package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Fomiller/mixify/pkg/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zmb3/spotify"
)

const redirectURL = "http://localhost:42069/callback/"

var (
	ch     = make(chan *spotify.Client)
	client *spotify.Client
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	if err := tea.NewProgram(ui.NewModel()).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	// load env vars
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	// // spotify setup
	// auth.SetAuthInfo(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))

	// // http server setup
	// http.HandleFunc("/callback/", completeAuth)
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Println("Got request for:", r.URL.String())
	// })
	// go http.ListenAndServe(":42069", nil)

	// url := auth.AuthURL(state)
	// go fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
	// client := <-ch

	// // use the client to make calls that require authorization
	// user, err := client.CurrentUser()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("You are logged in as:", user.ID)

	// _, playlist, err := client.FeaturedPlaylists()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Playlist: ", playlist)
	// // tui setup
	// p := tea.NewProgram(initialModel())
	// if err := p.Start(); err != nil {
	// 	fmt.Printf("Alas, there's been an error: %v", err)
	// 	os.Exit(1)
	// }
}
