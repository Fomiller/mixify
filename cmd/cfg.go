package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Fomiller/mixify/pkg/auth"
)

func createConfig(homeDir string) {
	// create config dir and parent folders if it doesnt exist
	err := os.MkdirAll(fmt.Sprintf("%s/.config/mixify", homeDir), 0777)
	if err != nil {
		panic(err)
	}

	// create config file if it doesnt exist
	_, err = os.Create(fmt.Sprintf("%s/.config/mixify/credentials.json", homeDir))
	if err != nil {
		panic(err)
	}
}

func loadToken() error {
	// create cfg path constant
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// get info about config file
	f, err := os.Stat(fmt.Sprintf("%s/.config/mixify/credentials.json", homeDir))
	if err != nil {
		return err
	}

	// check file for contents and load if it does
	if f.Size() > 0 {
		buf, err := os.ReadFile(fmt.Sprintf("%s/.config/mixify/credentials.json", homeDir))
		err = json.Unmarshal(buf, &auth.Token)
		// return err from loading token
		if err != nil {
			return err
		}
	}

	// return nil if load was successful
	return nil
}
