package utils

import (
	"fmt"
	"os"
)

// Config is a config interface.
var Config config

type config struct {
	DiscordToken string
	CAHChannelID string
	Prefix       string
}

func (conf *config) LoadConfig() {

	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		fmt.Println("Environment file not found, cannot continue!")
		os.Exit(1)
	}
	Config = config{
		DiscordToken: os.Getenv("DiscordToken"),
		CAHChannelID: os.Getenv("CAHChannelID"),
		Prefix:       os.Getenv("Prefix"),
	}
}
