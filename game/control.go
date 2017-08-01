package game

import (
	"github.com/FM1337/Discord-CAH/utils"
	"github.com/bwmarrin/discordgo"
)

// Start will start the game.
func Start(s *discordgo.Session, m *discordgo.MessageCreate) {
	if Running || Starting {
		s.ChannelMessageSend(utils.Config.CAHChannelID, "A game is already running!")
		return
	}
	// Initialize the data before adding players.
	InitializeData()
}

// Pause will pause the game.
func Pause(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !Running {
		s.ChannelMessageSend(utils.Config.CAHChannelID, "No game is running!")
		return
	}
}

// Stop will stop the game
func Stop(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !Running {
		s.ChannelMessageSend(utils.Config.CAHChannelID, "No game is running!")
		return
	}
}

// Join will join you into the game
func Join(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !Running {
		s.ChannelMessageSend(utils.Config.CAHChannelID, "No game is running!")
		return
	}
}

// Leave will remove you from the game
func Leave(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !Running {
		s.ChannelMessageSend(utils.Config.CAHChannelID, "No game is running!")
		return
	}
}
