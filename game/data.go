package game

import (
	"fmt"

	"github.com/FM1337/Discord-CAH/cards"
)

// Player is a struct that holds a player's data.
type Player struct {
	PlayerName string      // The player's Discord user name.
	PlayerID   string      // The player's Discord ID.
	Zar        bool        // Is the player the card zar?
	Cards      []WhiteCard // The player's hand.
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

// RoundResults is a struct that contains data about the chosen cards by players.
type RoundResults struct {
	PlayerName  string      // The Player's Discord name.
	PlayerID    string      // The Player's Discord ID.
	PlayerCards []WhiteCard // Slice of WhiteCard struct.
}

// Bools
// Starting is a bool to tell if the game is in the start up stage.
var Starting bool

// Running is a bool to tell if the game is running.
var Running bool

// Paused is a bool to tell if the game is paused.
var Paused bool

// Maps
// BlackCards is a map of the BlackCard struct.
var BlackCards map[int]BlackCard

// WhiteCards is a map of the WhiteCard struct.
var WhiteCards map[int]WhiteCard

// InitializeData will prepare the maps and slices for the game
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
}

// ImportBlackCards will import the black cards.
func ImportBlackCards() {
	// Loop through the black cards imported at start up of the bot and
	// add them to our map of black cards.
	for i, card := range cards.CardList.BlackCards {
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
