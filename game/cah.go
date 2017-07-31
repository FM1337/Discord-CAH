package game

import (
	"fmt"
	"time"

	"github.com/FM1337/Discord-CAH/utils"
	"github.com/bwmarrin/discordgo"
)

type Player struct {
	PlayerID string
	Cards    []Card
	Points   int
}

type Card struct {
	Text  string
	Blank bool
}

var Players []Player

var running bool
var paused bool

var creatorID string
var pauserID string

func StartGame(s *discordgo.Session, m *discordgo.MessageCreate) {
	// TODO: Game code

	if running {
		s.ChannelMessageSend(utils.Config.CAHChannelID, "A game is already running")
		return
	}

	running = true
	creatorID = m.Author.ID
	s.ChannelMessageSend(utils.Config.CAHChannelID, fmt.Sprintf("%s started a game!", m.Author.Username))
	Players = append(Players, Player{PlayerID: m.Author.ID, Cards: nil, Points: 0})

	s.ChannelMessageSend(utils.Config.CAHChannelID, fmt.Sprintf("The game will start in 30 seconds, type %sjoin to join in!", utils.Config.Prefix))

	time.Sleep(30 * time.Second)
	if len(Players) < 3 {
		s.ChannelMessageSend(utils.Config.CAHChannelID, "Not enough players to start the game!")
		running = false
		creatorID = ""
		return
	}

	return
}

func PauseGame(s *discordgo.Session, m *discordgo.MessageCreate) {
	// TODO: Game code
	if !running {
		s.ChannelMessageSend(utils.Config.CAHChannelID, "No game is running")
		return
	}

	if !paused {
		paused = true
		s.ChannelMessageSend(utils.Config.CAHChannelID, "The game has been paused!")
		pauserID = m.Author.ID
		return
	}

	if m.Author.ID == pauserID || m.Author.ID == creatorID {
		paused = false
		s.ChannelMessageSend(utils.Config.CAHChannelID, "The game has been unpaused!")
		pauserID = ""
		return
	}

	pauser, _ := s.User(pauserID)
	pauserName := pauser.Username
	creator, _ := s.User(creatorID)
	creatorName := creator.Username

	s.ChannelMessageSend(utils.Config.CAHChannelID, fmt.Sprintf("Sorry only %s or %s can unpause the game!", pauserName, creatorName))
	return

}

// StopGame stops the game.
func StopGame(s *discordgo.Session, m *discordgo.MessageCreate) {
	// TODO: Game code
	if !running {
		s.ChannelMessageSend(utils.Config.CAHChannelID, "No game is running")
		return
	}

	for _, player := range Players {
		if m.Author.ID == player.PlayerID {
			s.ChannelMessageSend(utils.Config.CAHChannelID, fmt.Sprintf("%s stopped the game!", m.Author.Username))
			running = false
			Players = nil
			return
		}
	}

	s.ChannelMessageSend(utils.Config.CAHChannelID, "You must be a player to stop the game!")
	return
}
