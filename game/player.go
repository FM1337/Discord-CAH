package game

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/FM1337/Discord-CAH/utils"
	"github.com/bwmarrin/discordgo"
)

// PickCard will allow players to pick a card to play from their hand.
func PickCard(s *discordgo.Session, m *discordgo.MessageCreate) {
	// If Judging don't allow card changes
	if Judging {
		return
	}

}

// ChooseWinner will allow the Card Zar to choose a winner.
func ChooseWinner(s *discordgo.Session, m *discordgo.MessageCreate) {
	// If not judging we don't want winners picked yet
	if !Judging {
		return
	}
}

// MessageHand will message a player their cards.
func MessageHand(PlayerID string, s *discordgo.Session) {
	// If the player is the zar, then we don't want to send them their hand.
	if Players[PlayerID].Zar {
		return
	}

	// MessageError is a bool that is set to true in the case of an error.
	MessageError := false
	// fields is a MessageEmbedField slice that holds fields for a discord message embed.
	fields := []*discordgo.MessageEmbedField{}

	// Open up a message channel to message the player their hand.
	channel, err := s.UserChannelCreate(PlayerID)

	if err != nil {
		// Set Paused to true.
		MessageError = true
	}

	// Append to the fields slice.
	for _, card := range Players[PlayerID].Cards {
		iString := strconv.Itoa(card.Index)
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   iString,
			Value:  card.Text,
			Inline: true,
		})
	}
	// The embed to send.
	embed := &discordgo.MessageEmbed{
		Type:        "rich",
		Color:       52,
		Title:       fmt.Sprintf("Your cards (%d)", len(Players[PlayerID].Cards)),
		Description: strings.Replace(RoundText, "_", "______", -1),
		Fields:      fields,
	}

	_, err = s.ChannelMessageSendEmbed(channel.ID, embed)

	if err != nil {
		// Set Paused to true.
		MessageError = true
	}
	// If there was an error, message the channel letting the user know the problem
	// then wait 30 seconds and try again.
	if MessageError {
		for MessageError {
			s.ChannelMessageSend(utils.Config.CAHChannelID, fmt.Sprintf("Hey %s, it looks like there was a problem sending your cards.\nA probable cause could be that you have the setting that allows server members to message you turned off, can you please try turning that on so I can send you your cards? I'll wait 30 seconds before trying again.\nError details: %s", Players[PlayerID].PlayerName, err))
			Wait30Seconds()
			if !Running {
				EndGame(s)
			}
			channel, err = s.UserChannelCreate(PlayerID)
			_, err = s.ChannelMessageSendEmbed(channel.ID, embed)
			if err != nil {
				continue
			}
			MessageError = false
		}
	}
	return
}
