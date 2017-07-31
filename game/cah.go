package game

import (
	"fmt"
	"time"

	"github.com/FM1337/Discord-CAH/utils"
	"github.com/bwmarrin/discordgo"
)

type Player struct {
	PlayerID string
	Cards    []WhiteCard
	Points   int
}

type WhiteCard struct {
	CardID int
	Text   string
	Blank  bool
}

type BlackCard struct {
	CardID int
	Text   string
	Cards  int
}

var Players map[string]Player = make(map[string]Player)

var running bool
var starting bool
var paused bool

var creatorID string
var pauserID string

func StartGame(s *discordgo.Session, m *discordgo.MessageCreate) {
	// TODO: Game code
	if running || starting {
		s.ChannelMessageSend(utils.Config.CAHChannelID, "A game is already running")
		return
	}

	starting = true
	creatorID = m.Author.ID
	s.ChannelMessageSend(utils.Config.CAHChannelID, fmt.Sprintf("%s started a game!", m.Author.Username))
	Players[m.Author.ID] = Player{PlayerID: m.Author.ID, Cards: nil, Points: 0}

	s.ChannelMessageSend(utils.Config.CAHChannelID, fmt.Sprintf("The game will start in 30 seconds, type %sjoin to join in!", utils.Config.Prefix))

	time.Sleep(30 * time.Second)
	fmt.Printf("%d players", len(Players))
	if len(Players) < 3 && starting {
		s.ChannelMessageSend(utils.Config.CAHChannelID, "Not enough players to start the game!")
		starting = false
		creatorID = ""
		Players = make(map[string]Player)
		return
	}

	starting = false
	running = true
	s.ChannelMessageSend(utils.Config.CAHChannelID, "The game has started!")
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
		if !starting {
			s.ChannelMessageSend(utils.Config.CAHChannelID, "No game is running")
			return
		}
	}

	for _, player := range Players {
		if m.Author.ID == player.PlayerID {
			s.ChannelMessageSend(utils.Config.CAHChannelID, fmt.Sprintf("%s stopped the game!", m.Author.Username))
			running = false
			starting = false
			Players = make(map[string]Player)
			return
		}
	}

	s.ChannelMessageSend(utils.Config.CAHChannelID, "You must be a player to stop the game!")
	return
}

func JoinGame(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !running {
		if !starting {
			s.ChannelMessageSend(utils.Config.CAHChannelID, fmt.Sprintf("No game is running! Type %sstart to start a game", utils.Config.Prefix))
			return
		}

	}

	for _, player := range Players {
		if m.Author.ID == player.PlayerID {
			s.ChannelMessageSend(utils.Config.CAHChannelID, "You've already joined!")
			return
		}
	}
	Players[m.Author.ID] = Player{PlayerID: m.Author.ID, Cards: nil, Points: 0}
	s.ChannelMessageSend(utils.Config.CAHChannelID, fmt.Sprintf("%s joined the game!", m.Author.Username))
	return
}

func LeaveGame(s *discordgo.Session, m *discordgo.MessageCreate) {

	if !running {
		if !starting {
			s.ChannelMessageSend(utils.Config.CAHChannelID, "No game is running")
			return
		}
	}

	delete(Players, m.Author.ID)
	s.ChannelMessageSend(utils.Config.CAHChannelID, fmt.Sprintf("%s has left the game!", m.Author.Username))

	if len(Players) < 3 && running {
		s.ChannelMessageSend(utils.Config.CAHChannelID, "Not enough players to continue!")
		running = false
		Players = make(map[string]Player)
		return
	}

	return
}
