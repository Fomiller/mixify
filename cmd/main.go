package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Fomiller/mixify/internal/auth"
	"github.com/Fomiller/mixify/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/browser"
	"github.com/zmb3/spotify/v2"
)

const redirectURL = "http://localhost:42069/callback/"

var (
	AuthUrl string = auth.Auth.AuthURL(auth.State)
	User    *spotify.PrivateUser
)

type config struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Expiration   string `json:"expiry"`
}

// this isnt really doing anything, until I find a way to store the credentials
func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	// check if config file exists
	_, err = os.Stat(fmt.Sprintf("%s/.config/mixify/credentials.json", homeDir))

	// create config dir and file if it doesnt exist
	if err != nil {
		createConfig(homeDir)
	}
}

func main() {
	// http server setup
	http.HandleFunc("/callback/", auth.CompleteAuth)
	// http.HandleFunc("/refresh/", auth.RefreshAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":42069", nil)

	// prompt user if they want to log into spotify account
	// this should only happen once, or if refresh token is expired
	c := login("Login to spotify?")
	//log in using Oauth
	if c != false {
		err := browser.OpenURL(AuthUrl)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// dont log in and exit
		fmt.Println("unable to login to spotify")
		os.Exit(1)
	}
	// load the token object now that we have a token

	auth.Client = <-auth.Ch

	// // // use the client to make calls that require authorization
	// User, err := auth.Client.CurrentUser(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("You are logged in as:", user.ID)

	// // // _, playlist, err := client.FeaturedPlaylists()
	// playlist, err := auth.Client.CurrentUsersPlaylists(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// p := playlist.Playlists[0]
	// fmt.Println("Playlist: ", p)
	// fmt.Println("ID: ", p.ID)
	// // fmt.Println("Tracks: ", p.Tracks)
	// // fmt.Println("Tracks Endpoint: ", p.Tracks)

	// tracklist, err := auth.Client.GetPlaylistTracks(context.Background(), p.ID)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("trackList: %v\n", tracklist)
	// fmt.Printf("trackList.Tracks: %v\n", tracklist.Tracks)
	// fmt.Printf("trackList.Tracks[0]: %v\n", tracklist.Tracks[0].Track.Name)
	// fmt.Printf("trackList.Tracks[0]: %v\n", tracklist.Tracks[0].Track.Album.Name)
	// fmt.Println("--------------------------")
	// for _, t := range tracklist.Tracks {
	// 	fmt.Println(t)
	// 	fmt.Println(t.Track)
	// 	fmt.Println("--------------------------")
	// }

	//refresh token stuff
	// token := make(map[string]string)
	// token["refresh_token"] = "AQAtTr7ysrfqsfadY7gHj5wjSKyTN3W7v5CyjAm1w5SQKj2YHQypAzPzlkB2zgPUdM85SDL3_zehxyeIf-nQ9SDEz9olGk89bCoUdBBxUPF5-V4KQsK2HVaz53Vov0GF_Us"
	// postData, err := json.Marshal(token)
	// if err != nil {
	// 	panic(err)
	// }

	// data := bytes.NewBuffer(postData)
	// data := url.Values{}
	// data.Add("refresh_token", "AQAtTr7ysrfqsfadY7gHj5wjSKyTN3W7v5CyjAm1w5SQKj2YHQypAzPzlkB2zgPUdM85SDL3_zehxyeIf")
	// resp, err := http.PostForm("http://localhost:42069/refresh/", data)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(resp.Body)
	// /////////////////////////////

	// tui setup
	// rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	if err := tea.NewProgram(ui.New(), tea.WithAltScreen()).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func login(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [Y/n]: ", s)
		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "" {
			response = "y"

		}

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
