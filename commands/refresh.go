package commands

import (
	"github.com/FM1337/Discord-CAH/cards"
	"github.com/FM1337/Discord-CAH/game"
	"github.com/bwmarrin/discordgo"
)

// Refresh will refresh the bot's cards, adding any new cards that might've
// been added.
func Refresh(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Do not refresh cards if the game is running
	if game.Running {
		s.ChannelMessageSend(m.ChannelID, "Sorry, you may not refresh the cards during a game!")
		return
	}
	// Set Refreshing to true.
	game.Refreshing = true

	// Blank out the current card slices.
	cards.CardList.BlackCards = nil
	cards.CardList.WhiteCards = nil

	// Now let's reload the cards.
	cards.LoadDefaultCards()

	// Now set Refreshing to false
	game.Refreshing = false
	s.ChannelMessageSend(m.ChannelID, "Cards refreshed successfully!")
	return
}
