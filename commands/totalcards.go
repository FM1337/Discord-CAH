package commands

import (
	"fmt"

	"github.com/FM1337/Discord-CAH/cards"
	"github.com/bwmarrin/discordgo"
)

func TotalCards(s *discordgo.Session, m *discordgo.MessageCreate) {
	totalBlackCards := len(cards.CardList.BlackCards)
	totalDefaultBlackCards := len(cards.CardListMap["default"].BlackCards)
	totalCustomBlackCards := len(cards.CardListMap["custom"].BlackCards)
	totalWhiteCards := len(cards.CardList.WhiteCards)
	totalDefaultWhiteCards := len(cards.CardListMap["default"].WhiteCards)
	totalCustomWhiteCards := len(cards.CardListMap["custom"].WhiteCards)
	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("There are a total of %d black cards (%d default, %d custom) and %d white cards (%d default, %d custom)!", totalBlackCards, totalDefaultBlackCards, totalCustomBlackCards, totalWhiteCards, totalDefaultWhiteCards, totalCustomWhiteCards))
	if err != nil {
		fmt.Printf("%s", err)
	}
}
