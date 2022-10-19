package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Fomiller/mixify/pkg/auth"
	"gopkg.in/yaml.v2"
)

const redirectURL = "http://localhost:42069/callback/"

var (
	// client *spotify.Client
	// Auth = spotifyauth.New(
	// 	spotifyauth.WithRedirectURL(redirectURL),
	// 	spotifyauth.WithScopes(
	// 		spotifyauth.ScopeUserReadCurrentlyPlaying,
	// 		spotifyauth.ScopeUserReadPlaybackState,
	// 		spotifyauth.ScopeUserModifyPlaybackState,
	// 		spotifyauth.ScopeUserLibraryModify,
	// 		spotifyauth.ScopeUserLibraryRead,
	// 		spotifyauth.ScopePlaylistModifyPublic,
	// 		spotifyauth.ScopePlaylistReadPrivate,
	// 	),
	// )
	Config config
)

type config struct {
	Token        string `yaml:"token"`
	RefreshToken string `yaml:"refreshToken"`
}

// func init() {
// 	fmt.Println(os.Getenv("SPOTIFY_ID"))
// 	homeDir, err := os.UserHomeDir()
// 	if err != nil {
// 		panic(err)
// 	}

// 	// check if config file exists
// 	_, err = os.Stat(fmt.Sprintf("%s/.config/mixify/config.yaml", homeDir))
// 	// create config dir and file if it doesnt exist
// 	if err != nil {
// 		createConfig(homeDir)
// 	}
// }

func main() {
	// err := readConfig()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(Config)

	// http server setup
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

func createConfig(homeDir string) {
	fmt.Println("creating config")
	// create config dir and parent folders if it doesnt exist
	err := os.MkdirAll(fmt.Sprintf("%s/.config/mixify", homeDir), 0777)
	if err != nil {
		panic(err)
	}

	// create config file if it doesnt exist
	_, err = os.Create(fmt.Sprintf("%s/.config/mixify/config.yaml", homeDir))
	if err != nil {
		panic(err)
	}
}

func readConfig() error {
	// create cfg path constant
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// read config file
	buf, err := os.ReadFile(fmt.Sprintf("%s/.config/mixify/config.yaml", homeDir))
	err = yaml.Unmarshal(buf, &Config)
	if err != nil {
		return err
	}

	return nil
}
