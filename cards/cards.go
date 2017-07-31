package cards

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type cards struct {
	BlackCards []blackCard
	WhiteCards []whiteCard
}

type blackCard struct {
	CardText   string
	Cards2Play int
}

type whiteCard struct {
	CardText string
}

var CardList cards

// LoadCards loads the cards into memory on startup
func LoadCards() {
	// Load default black cards.
	BCJson, err := ioutil.ReadFile("./cards/default/blackCards.json")

	if err != nil {
		fmt.Printf("Error loading black cards: %s", err)
		os.Exit(1)
	}

	err = json.Unmarshal(BCJson, &CardList.BlackCards)

	if err != nil {
		fmt.Printf("Error unmarshling black cards: %s", err)
		os.Exit(1)
	}

	// Load default white cards.
	WCJson, err := ioutil.ReadFile("./cards/default/whiteCards.json")

	if err != nil {
		fmt.Printf("Error loading white cards: %s", err)
		os.Exit(1)
	}

	err = json.Unmarshal(WCJson, &CardList.WhiteCards)

	if err != nil {
		fmt.Printf("Error unmarshling black cards: %s", err)
		os.Exit(1)
	}
}
