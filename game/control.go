package game

import (
	"fmt"
	"time"

	"github.com/FM1337/Discord-CAH/utils"
	"github.com/bwmarrin/discordgo"
)

// Start will start the game.
func Start(s *discordgo.Session, m *discordgo.MessageCreate) {
	if Running || Starting {
		s.ChannelMessageSend(m.ChannelID, "A game is already running!")
		return
	}
	// Make sure the cards aren't being refreshed
	if Refreshing {
		s.ChannelMessageSend(m.ChannelID, "Sorry the cards are in the progress of being refreshed, please try again in a few moments!")
		return
	}
	// Initialize the data before adding players.
	InitializeData()
	// Set the Creator ID
	CreatorID = m.Author.ID
	// Add the game starter to the player list.
	AddPlayer(m.Author)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s has started a game!", m.Author.Username))
	// Wait half a second before sending the next message
	time.Sleep(500 * time.Millisecond)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The game will start in 30 seconds\nType %sJoin to join!", utils.Config.Prefix))
	// Now that's start up is complete, we can change Running to true.
	Starting = false
	Running = true
	Wait30Seconds()
	// If the game has been stopped during the waiting period, don't continue.
	if !Running {
		return
	}

	// Check if PlayerCount is less than 3
	if PlayerCount < 3 {
		s.ChannelMessageSend(m.ChannelID, "Not enough players!")
		// Wait half a second before sending the next message.
		time.Sleep(500 * time.Millisecond)
		s.ChannelMessageSend(m.ChannelID, "Waiting 3 minutes for minimum amount needed!")
		ExtendedWait(s, m)
		// If the game has been stopped don't continue.
		if !Running {
			return
		}
	}
	s.ChannelMessageSend(m.ChannelID, "The game is starting!")
	// Let's take care of a few more things before we start the game.
	PrepareGame()
	RoundStart(s, m)
}

// Pause will pause the game.
func Pause(s *discordgo.Session, m *discordgo.MessageCreate) {
	// If a game is not running, obviously we don't pause it.
	if !Running {
		s.ChannelMessageSend(m.ChannelID, "No game is running!")
		return
	}

	// Make sure that the person trying to pause/unpause the game is
	// actually in the game.
	if !UserInGame(m.Author.ID) {
		s.ChannelMessageSend(m.ChannelID, "Only players may pause the game!")
		return
	}

	// If the game is paused, we want to unpause it
	if Paused {
		// First let's check if the person attempting to unpause it
		// is the one who paused it or started the game
		if m.Author.ID == PauserID || m.Author.ID == CreatorID {
			s.ChannelMessageSend(m.ChannelID, "The game has been unpaused!")
			Paused = false
			// Set the PauserID to blank now that the game is unpaused.
			PauserID = ""
			return
		}
		// If the IDs do not match do the following.
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Only %s or %s may unpause the game", GetUserName(PauserID, s), GetUserName(CreatorID, s)))
		return

	}

	// If the game is not paused and the user is in the game, let's pause it
	Paused = true
	// Set PauserID to the Pauser's Discord ID.
	PauserID = m.Author.ID
	s.ChannelMessageSend(m.ChannelID, "The game has been paused!")
}

// Stop will stop the game
func Stop(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check if the game is running
	if !Running {
		s.ChannelMessageSend(m.ChannelID, "No game is running!")
		return
	}

	// Check if the user is a player
	if !UserInGame(m.Author.ID) {
		s.ChannelMessageSend(m.ChannelID, "Only players may pause the game!")
		return
	}

	Running = false
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s stopped the game!", m.Author.Username))
}

// Join will join you into the game
func Join(s *discordgo.Session, m *discordgo.MessageCreate) {

	if !Running {
		s.ChannelMessageSend(m.ChannelID, "No game is running!")
		return
	}

	// Check if user already in game.
	if UserInGame(m.Author.ID) {
		s.ChannelMessageSend(m.ChannelID, "You've already joined!")
		return
	}

	// Otherwise add them to the game.
	AddPlayer(m.Author)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s has joined!", m.Author.Username))
}

// Leave will remove you from the game
func Leave(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Can't leave a game that ain't running.
	if !Running {
		s.ChannelMessageSend(m.ChannelID, "No game is running!")
		return
	}

	// First check if the user attempting to leave is actually in the game
	if !UserInGame(m.Author.ID) {
		s.ChannelMessageSend(m.ChannelID, "Only players may leave the game!")
		return
	}

	RemovePlayer(s, m)
}

// Wait30Seconds will wait 30 seconds, it also respects if the game is paused
// or stopped.
func Wait30Seconds() {
	// The 30 second for loop.
	for i := 0; i <= 30; i++ {
		// Check to see if has been stopped
		if !Running {
			return
		}
		// Check to see if the game is paused.
		if Paused {
			// If it is paused, take 1 away from i
			i = i - 1
		}
		// Sleep for 1 second
		time.Sleep(1 * time.Second)
	}

}

// Wait1Minute will wait 1 minute, it also respects if the game is paused
// or stopped.
func Wait1Minute() {
	// The 1 minute for loop.
	for i := 0; i <= 60; i++ {
		// Sleep for 1 second
		time.Sleep(1 * time.Second)
		// Check to see if has been stopped
		if !Running {
			return
		}
		// Check to see if the game is paused.
		if Paused {
			// If it is paused, take 1 away from i
			i = i - 1
			continue
		}
		// Check to see if judging is over.
		if !Judging {
			return
		}
	}

}

// ExtendedWait will wait 3 minutes for the minimum amount of players,
//it also respects if the game is paused or stopped.
func ExtendedWait(s *discordgo.Session, m *discordgo.MessageCreate) {
	// The 3 minute for loop.
	for i := 0; i <= 180; i++ {
		// Check to see if has been stopped
		if !Running {
			return
		}
		// Check to see if the game is paused.
		if Paused {
			// If it is paused, take 1 away from i
			i = i - 1
		}
		// Sleep for 1 second

		// Check to see if the minmum player limit has been hit
		if PlayerCount >= 3 {
			return
		}

		time.Sleep(1 * time.Second)
	}
	// If the time is up and we still haven't got the required amount,
	// then we stop the game.
	Running = false
	s.ChannelMessageSend(m.ChannelID, "Not enough players, game has been stopped!")
	return

}
