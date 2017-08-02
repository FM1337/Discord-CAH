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
	// Check to see if the game is running.
	if !Running {
		s.ChannelMessageSend(m.ChannelID, "No game is running!")
		return
	}

	// Check to see if a player is actually playing
	if !UserInGame(m.Author.ID) {
		s.ChannelMessageSend(m.ChannelID, "Only players may pick cards!")
		return
	}

	// If Judging don't allow card changes
	if Judging {
		return
	}
	// Split up the message into arguments
	args := strings.Split(m.Content, " ")
	// Make sure we have the correct number of cards being played.
	if len(args)-1 == BlackCards[RoundCardID].Cards {
		// Loop through our message arguments and see if they are correct.
		for _, arg := range args[1:] {
			// tmpCard is a temporary variable used to hold a matching
			// white card is a match is found.
			var tmpCard WhiteCard

			// match is a temporary bool that tells us if a match is found.
			match := false

			cardNum, err := strconv.Atoi(arg)

			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Invalid input!")
				return
			}

			// Now let's loop through our cardlist and see if a match
			// can be found.
			for _, card := range Players[m.Author.ID].Cards {
				// if a match is found.
				if card.Index == cardNum {
					match = true
					tmpCard = card
					break
				}
			}
			if !match {
				s.ChannelMessageSend(m.ChannelID, "Invalid input!")
				return
			}
			// tmpPlayer is a temporary copy of a player to update it's
			// entry in the players map
			tmpPlayer := Players[m.Author.ID]
			tmpPlayer.PlayedCards = append(tmpPlayer.PlayedCards, tmpCard)

			// Update the player entry
			Players[m.Author.ID] = tmpPlayer

		}
		// Now let's loop so that we can generate a message to send to
		// the player.
		// tmpString is a temporary string
		tmpString := RoundText
		for _, card := range Players[m.Author.ID].PlayedCards {
			tmpString = strings.Replace(tmpString, "_", card.Text, 1)
			// If no underscores, then add the card to the end.
			if !strings.Contains(tmpString, "_") {
				tmpString = fmt.Sprintf("%s %s", tmpString, card.Text)
			}
		}

		// Once the loop is done send a message to the player with the result.
		channel, _ := s.UserChannelCreate(m.Author.ID)
		s.ChannelMessageSend(channel.ID, fmt.Sprintf("You've played: %s", tmpString))
	}
}

// ChooseWinner will allow the Card Zar to choose a winner.
func ChooseWinner(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check to see if the game is running.
	if !Running {
		s.ChannelMessageSend(m.ChannelID, "No game is running!")
		return
	}

	// Check to see if a player is actually playing
	if !UserInGame(m.Author.ID) {
		s.ChannelMessageSend(m.ChannelID, "Only players may use this command!")
		return
	}

	// If not judging we don't want winners picked yet
	if !Judging {
		return
	}
	// Split up the message into arguments
	args := strings.Split(m.Content, " ")

	// Make sure we're only choosing one winner
	if len(args) == 2 {

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
