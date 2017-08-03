package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/FM1337/Discord-CAH/cards"
	"github.com/FM1337/Discord-CAH/utils"
	"github.com/bwmarrin/discordgo"
)

// Structs
// Player is a struct that holds a player's data.
type Player struct {
	PlayerName  string      // The player's Discord user name.
	PlayerID    string      // The player's Discord ID.
	Zar         bool        // Is the player the card zar?
	Cards       []WhiteCard // The player's hand.
	PlayedCards []WhiteCard // The cards the player played.
	Score       int         // The player's score
}

// WhiteCard is a struct that holds the data about a white card.
type WhiteCard struct {
	CardID int    // The card's ID.
	Index  int    // The card's index value
	Text   string // The card's text.
	Blank  bool   // Is the card a blank card?
	taken  bool   // Is the card taken by a player?
}

// BlackCard is a struct that holds the data about a black card.
type BlackCard struct {
	CardID int    // The card's ID.
	Text   string // The card's text
	Cards  int    // The amount of cards to play for this card.
}

// RoundResult is a struct that contains data about the chosen cards by players.
type RoundResult struct {
	PlayerName string // The Player's Discord name.
	PlayerID   string // The Player's Discord ID.
	PlayString string // The result of the chosen cards.
}

// Bools
// Starting is a bool to tell if the game is in the start up stage.
var Starting bool

// Running is a bool to tell if the game is running.
var Running bool

// Paused is a bool to tell if the game is paused.
var Paused bool

// Judging is a bool to tell if the round is in the judging stage.
var Judging bool

// Refreshing is a bool to tell if the cards are currently being refreshed.
var Refreshing bool

// Maps
// BlackCards is a map of the BlackCard struct.
var BlackCards map[int]BlackCard

// WhiteCards is a map of the WhiteCard struct.
var WhiteCards map[int]WhiteCard

// Players is a map of the Player struct.
var Players map[string]Player

// Ints
// PlayerCount is an int that shows how many player have joined.
var PlayerCount int

// Rounds is an int that holds the number of rounds to be played in a game.
var Rounds int

// Round is an int that holds the current round number.
var Round int

// Zar is an int that holds the current Zar's index number of Zars.
var Zar int

// RoundCardID contains the round's black card ID
var RoundCardID int

// HighScore is an int that holds the highest score of the game.
var HighScore int

// Strings
// CreatorID is a string containing the Discord ID of the person who
// started the game.
var CreatorID string

// PauserID is a string containing the Discord ID of the person who
// paused the game.
var PauserID string

//  RoundText is a string containing the round's black card text.
var RoundText string

// HighScoreID is a string containing the ID of the player with the highest score.
var HighScoreID string

// Slices
// Zars is a string slice that will contain the order of which the next
// round's card zar will be chosen.
var Zars []string

// RoundResults is a slice of RoundResult.
var RoundResults []RoundResult

// InitializeData will prepare the maps and slices for the game.
func InitializeData() {
	// First we do the black cards
	BlackCards = make(map[int]BlackCard)
	ImportBlackCards()
	// Print the amount of Black Cards loaded into memory.
	fmt.Printf("%d Black Cards loaded!\n", len(BlackCards))

	// Then we do the white cards.
	WhiteCards = make(map[int]WhiteCard)
	ImportWhiteCards()
	// Print the amount of White Cards loaded into memory.
	fmt.Printf("%d White Cards loaded!\n", len(WhiteCards))

	// Moving on, we now want to make the players map.
	Players = make(map[string]Player)

	// We want to clear the following variables..
	PlayerCount = 0
	HighScore = 0
	Zar = 0
	Round = 0
	Rounds = 10
	RoundText = ""
	RoundResults = nil
	CreatorID = ""
	Zars = nil
	Judging = false
}

// ImportBlackCards will import the black cards.
func ImportBlackCards() {
	// Loop through the black cards imported at start up of the bot and
	// add them to our map of black cards.
	for i, card := range cards.CardList.BlackCards {
		// If there are no underscores, we still have to play a card.
		if card.Cards2Play == 0 {
			card.Cards2Play = 1
		}

		BlackCards[i] = BlackCard{
			CardID: i,
			Text:   card.CardText,
			Cards:  card.Cards2Play,
		}
	}
}

// ImportWhiteCards will import the white cards.
func ImportWhiteCards() {
	// Loop through the white cards imported at start up of the bot and
	// add them to our map of white cards.
	for i, card := range cards.CardList.WhiteCards {
		WhiteCards[i] = WhiteCard{
			CardID: i,
			Text:   card.CardText,
			taken:  false,
		}
	}
}

// AddPlayer adds a player to the game.
func AddPlayer(User *discordgo.User) {
	Players[User.ID] = Player{
		PlayerName: User.Username,
		PlayerID:   User.ID,
		Zar:        false,
	}
	// Add the player to the Zars list only if the game hasn't begun
	if Round != 0 {
		Zars = append(Zars, User.ID)
		GenerateHand(User.ID)
	}
	// Add 1 to PlayerCount.
	PlayerCount++
}

// RemovePlayer removes a player from the game.
func RemovePlayer(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s has left the game!", m.Author.Username))
	// Let's check if the game is paused, this is important because
	// if the person leaving the game is the one who paused, we want to
	// unpause it.
	if Paused {
		if m.Author.ID == PauserID {
			// Instead of writing extra code, let's just reference the
			// pause function.
			Pause(s, m)
		}
	}

	// Next we want to check if the person is the CardZar
	if Players[m.Author.ID].Zar {
		s.ChannelMessageSend(m.ChannelID, "The CardZar has left, next round!")
		// Reference a function that will move players to the next round.
	}

	// Return the player's cards to the deck.
	for _, card := range Players[m.Author.ID].Cards {
		ReleaseCard(card.CardID)
	}

	// Remove the player from the Zars list
	// We'll generate a temporary ZarsList to update the global one with.
	TmpZars := []string{}
	for _, zar := range Zars {
		if zar != m.Author.ID {
			TmpZars = append(TmpZars, zar)
		}
	}
	Zars = TmpZars

	// Moving on let's delete the player from the Players map.
	delete(Players, m.Author.ID)
	// then we remove one from the PlayerCount
	PlayerCount--
	// Finally we check to see if there's still enough players to keep playing
	if PlayerCount < 3 {
		s.ChannelMessageSend(m.ChannelID, "Not enough players to keep playing!")
		Running = false
		EndGame(s)
		return
	}
}

// UserInGame checks if a user is in the game.
func UserInGame(PlayerID string) bool {
	for _, player := range Players {
		// If a match is found return true
		if player.PlayerID == PlayerID {
			return true
		}
	}
	// Otherwise return false.
	return false
}

// GetUserName returns the username of the Discord User that matches the
// given ID.
func GetUserName(ID string, s *discordgo.Session) string {
	// Create a user object based on the ID
	user, err := s.User(ID)

	// if no user is found just return the ID.
	if err != nil {
		return ID
	}
	// Otherwise return the user's username.
	return user.Username
}

// ReleaseCard will release the player's card back to the deck to allow it
// to be used by other players.
func ReleaseCard(CardID int) {
	// tmpcard is a temporary element which will allow us to update the
	// card's taken bool.
	tmpcard := WhiteCards[CardID]
	// Set taken to false.
	tmpcard.taken = false
	// Update the map entry with the new data.
	WhiteCards[CardID] = tmpcard
}

// GenerateHand will generate a player's hand at the beginning of the game.
func GenerateHand(PlayerID string) {
	// tmpCards is a temporary slice used to hold the cards to be added
	// to a player's hand.
	tmpCards := []WhiteCard{}
	// tmpPlayer is a temporary copy of a player in the Players map and
	// is used to update the player's information with the new cards.
	tmpPlayer := Players[PlayerID]
	for i := 0; i <= 9; i++ {
		for {
			// RandomCard is a random card chosen from the WhiteCards map.
			RandomCard := WhiteCards[rand.Intn(len(WhiteCards))]
			// Check if the RandomCard isn't already taken
			if !RandomCard.taken {
				RandomCard.taken = true
				RandomCard.Index = i + 1
				tmpCards = append(tmpCards, RandomCard)
				WhiteCards[RandomCard.CardID] = RandomCard
				break
			}
		}
	}
	// Set the player's cards.
	tmpPlayer.Cards = tmpCards
	// Update the player's data.
	Players[PlayerID] = tmpPlayer
}

// PrepareGame takes care of what's left to do before the game starts.
func PrepareGame() {
	// Let's generate the player's hands
	for _, player := range Players {
		// Generate the player's hand.
		GenerateHand(player.PlayerID)
		// Add the player to the Zars slice.
		Zars = append(Zars, player.PlayerID)
	}
	// Set round to 1
	Round = 1
}

// NextZar chooses the next zar in the Zars slice.
func NextZar() {
	// Remove current Zar
	tmpPlayer := Players[Zars[Zar]]
	tmpPlayer.Zar = false
	Players[tmpPlayer.PlayerID] = tmpPlayer

	// If length of the Zars list minus 1 is equal to the current Zar
	// then we set Zar to 0
	if len(Zars)-1 == Zar {
		Zar = 0
		return
	}
	// If not, then just add one to Zar
	Zar++

}

// EndGame is the function that is run at the end of the game.
func EndGame(s *discordgo.Session) {
	// Loop through all players
	for _, player := range Players {
		// if a player's score is higher than the current one, then
		// update the high score with the player's score.
		if player.Score > HighScore {
			HighScore = player.Score
			HighScoreID = player.PlayerID
		}
	}
	// Error catching here, if HighScoreID is blank, then no winners.
	if HighScoreID == "" {
		s.ChannelMessageSend(utils.Config.CAHChannelID, "The game has ended with no winners. Better luck next time!")
		return
	}

	s.ChannelMessageSend(utils.Config.CAHChannelID, fmt.Sprintf("Congratulations %s You won the game with %d points!", Players[HighScoreID].PlayerName, HighScore))
	// At this point, I would reset all the variables and stuff, but the
	// InitializeData function that runs when a new game is started
	// does it for me.
}

// SwapCard replaces the used cards with new ones after each round.
func SwapCard() {

	// oldCards is a slice that holds the used cards.
	oldCards := []WhiteCard{}

	// Loop through the player list.
	for _, player := range Players {

		// Loop through the player's played cards
		for _, card := range player.PlayedCards {
			oldCards = append(oldCards, card)
		}

		// Loop through the player's hand
		for i, card := range player.Cards {
			// Then do loop of the oldCards list
			for _, oldCard := range oldCards {
				// If the card is in the old cards list, get a new one.
				if card.CardID == oldCard.CardID {
					card = DrawCard(card.Index)
					player.Cards[i] = card
					break
				}
			}
		}
		// Update the player's data
		Players[player.PlayerID] = player
	}

	// Finally let's loop through the old cards list and release them
	// all back into the deck.
	for _, card := range oldCards {
		card.taken = false
		card.Index = 0
		WhiteCards[card.CardID] = card
	}

}

// DrawCard draws a random card and returns it
func DrawCard(index int) WhiteCard {
	for {

		randomCard := WhiteCards[rand.Intn(len(WhiteCards))]
		// If the random card isn't already taken
		if !randomCard.taken {
			// Set it as taken
			randomCard.taken = true
			randomCard.Index = index
			fmt.Printf("Card text: %s\n", randomCard.Text)
			// Update it in the map
			WhiteCards[randomCard.CardID] = randomCard
			return randomCard
		}
		// Wait half a second before looping again
		time.Sleep(500 * time.Millisecond)
	}
}
