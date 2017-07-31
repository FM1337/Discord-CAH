package commands

import (
	"fmt"

	"github.com/FM1337/Discord-CAH/cards"
	"github.com/bwmarrin/discordgo"
)

func TotalCards(s *discordgo.Session, m *discordgo.MessageCreate) {
	totalBlackCards := len(cards.CardList.BlackCards)
	totalWhiteCards := len(cards.CardList.WhiteCards)
	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("There are a total of %d black cards and %d white cards!", totalBlackCards, totalWhiteCards))
	if err != nil {
		fmt.Printf("%s", err)
	}
}
