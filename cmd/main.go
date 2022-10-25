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
	Config config
)

type config struct {
	Token        string `yaml:"token"`
	RefreshToken string `yaml:"refreshToken"`
}

func init() {
	fmt.Println(os.Getenv("SPOTIFY_ID"))
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	// check if config file exists
	_, err = os.Stat(fmt.Sprintf("%s/.config/mixify/config.yaml", homeDir))
	// create config dir and file if it doesnt exist
	if err != nil {
		createConfig(homeDir)
	}
}

func main() {
	// err := readConfig()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(Config)

	// http server setup
	http.HandleFunc("/callback/", auth.CompleteAuth)
	// http.HandleFunc("/refresh/", auth.RefreshAuth)
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Println("Got request for:", r.URL.String())
	// })
	go http.ListenAndServe(":42069", nil)

	authUrl := auth.Auth.AuthURL(auth.State)
	fmt.Println(authUrl)
	fmt.Println("----------")

	fmt.Printf("Please log in to Spotify by visiting the following page in your browser: %s\n", authUrl)
	client := &http.Client{}
	res, err := client.Get(authUrl)
	if err != nil {
		panic("error logging in to spotify")
	}
	// fmt.Println()
	fmt.Println("----------")
	fmt.Println(res)
	fmt.Println("----------")
	fmt.Println(res.Request.URL.String())
	fmt.Println("----------")
	fmt.Println(res.Cookies())

	res, err = client.Get(res.Request.URL.String())
	if err != nil {
		panic("error logging in to spotify")
	}
	fmt.Println("----------")
	fmt.Println("----------")
	fmt.Println(res)
	fmt.Println("----------")
	fmt.Println(res.Request.URL.String())
	fmt.Println("----------")
	fmt.Println(res.Cookies())
	// fmt.Println(res.Location())
	x, err := client.Get(res.Request.URL.String())
	if err != nil {
		panic("error logging in to spotify")
	}
	fmt.Println("----------")
	fmt.Println("----------")
	fmt.Println(x)
	fmt.Println("----------")
	fmt.Println(x.Request.URL.String())
	// fmt.Println("----------")
	// fmt.Println(x.Cookies())

	// callbackUrl := res.Cookies()
	// fmt.Println(callbackUrl)
	// fmt.Println(callbackUrl.Path)
	res, err = http.Get("localhost:42069/callback/")
	if err != nil {
		panic("error logging calling callback")
	}
	fmt.Println(res)

	auth.Client = <-auth.Ch

	// // // use the client to make calls that require authorization
	user, err := auth.Client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("You are logged in as:", user.ID)

	// // _, playlist, err := client.FeaturedPlaylists()
	// playlist, err := auth.Client.CurrentUsersPlaylists(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// p := playlist.Playlists[0]
	// fmt.Println("Playlist: ", p)
	// fmt.Println("ID: ", p.ID)

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

	// if err := tea.NewProgram(ui.New(), tea.WithAltScreen()).Start(); err != nil {
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
