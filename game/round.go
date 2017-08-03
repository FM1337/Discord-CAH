package game

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/FM1337/Discord-CAH/utils"
	"github.com/bwmarrin/discordgo"
)

// RoundStart holds the code that runs when a round starts.
func RoundStart(s *discordgo.Session, m *discordgo.MessageCreate) {
	// While the game is running
	for Running {
		// Choose a black card for the round
		RoundCardID = BlackCards[rand.Intn(len(BlackCards))].CardID
		RoundText = strings.Replace(BlackCards[RoundCardID].Text, "_", "______", -1)

		// If the RoundText is blank, then find another black card that has text.
		if RoundText == "" {
			for {
				RoundCardID = BlackCards[rand.Intn(len(BlackCards))].CardID
				RoundText = strings.Replace(BlackCards[RoundCardID].Text, "_", "______", -1)
				if RoundText != "" {
					break
				}
				// Wait half a second then continue.
				time.Sleep(500 * time.Millisecond)
			}
		}

		// Set the Zar
		TmpPlayer := Players[Zars[Zar]]
		TmpPlayer.Zar = true
		Players[Zars[Zar]] = TmpPlayer

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Round #%d", Round))
		// Wait half a second before sending each message
		time.Sleep(500 * time.Millisecond)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s", RoundText))
		time.Sleep(500 * time.Millisecond)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s is the Cardzar!", Players[Zars[Zar]].PlayerName))

		// Now send the players their cards.
		for _, player := range Players {
			MessageHand(player.PlayerID, s)
			// Wait half a second before sending the next
			time.Sleep(500 * time.Millisecond)
		}
		// Once that's done.
		s.ChannelMessageSend(m.ChannelID, "Players you have 30 seconds to choose what to play!")
		Wait30Seconds()
		// If the game has ended
		if !Running {
			EndGame(s)
			return
		}

		s.ChannelMessageSend(m.ChannelID, "Time is up!")
		Judging = true
		time.Sleep(500 * time.Millisecond)

		// Now let's generate the results.
		for _, player := range Players {
			// skip is a temporary bool that will tell if we should skip
			// appending.
			skip := false
			// let's make sure the player isn't the zar
			if player.Zar {
				continue
			}
			// if there are no played cards.
			if len(player.PlayedCards) == 0 {
				skip = true
			}
			// tmpString is a temporary string
			tmpString := BlackCards[RoundCardID].Text
			for _, card := range player.PlayedCards {
				// check to see if the card is blank
				if card.Text == "" {
					// If it is, then we break out of this loop and skip the player.
					skip = true
					break
				}
				if !strings.Contains(tmpString, "_") {
					tmpString = fmt.Sprintf("%s %s", tmpString, card.Text)
					continue
				}
				tmpString = strings.Replace(tmpString, "_", card.Text, 1)
				// If no underscores, then add the card to the end.
			}
			if !skip {
				RoundResults = append(RoundResults, RoundResult{
					PlayerName: player.PlayerName,
					PlayerID:   player.PlayerID,
					PlayString: tmpString,
				})
			}
			// Let's wait half a second before moving onto the next player
			time.Sleep(500 * time.Millisecond)
		}

		// if we don't have enough results
		if len(RoundResults) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Not enough results, next round!")
			NextRound(s)
			continue
		}

		s.ChannelMessageSend(m.ChannelID, "Here are the results:")
		// Now we loop through the results.
		for i, result := range RoundResults {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("#%d. %s", i+1, result.PlayString))
			// Let's wait half a second before sending the next
			time.Sleep(500 * time.Millisecond)
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Cardzar you have one minute to pick a winner!\nUse %sChoose to pick a winner!", utils.Config.Prefix))
		Wait1Minute()
		// If the CardZar did not pick a winner.
		if Judging {
			s.ChannelMessageSend(m.ChannelID, "Because the Cardzar did not pick a winner in time, a random one will be picked instead!")
			// Picking a random result
			randomResult := RoundResults[rand.Intn(len(RoundResults))]
			tmpWinPlayer := Players[randomResult.PlayerID]
			tmpWinPlayer.Score++
			Players[tmpWinPlayer.PlayerID] = tmpWinPlayer
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Congratulations %s, you win this round with %s", randomResult.PlayerName, randomResult.PlayString))
		}
		NextRound(s)
	}

}

// NextRound holds the code that runs just before the next round starts.
func NextRound(s *discordgo.Session) {
	// Add 1 to Round
	Round++
	// Check to see if we've played 10 rounds.
	if Round > Rounds {
		// If we've then end the game.
		Running = false
		EndGame(s)
		return
	}
	// Delete last round's card so that it doesn't get reused.
	delete(BlackCards, RoundCardID)
	// Choose the next Card Zar
	NextZar()
	// Replace used cards.
	SwapCard()
	// Set Judging to false
	Judging = false
	// Nil round results.
	RoundResults = nil

	// nil Player's chosen cards
	for _, player := range Players {
		player.PlayedCards = nil
		Players[player.PlayerID] = player
	}
}
