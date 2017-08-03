package cards

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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

// CardList is a cards interface.
var CardList cards

// CardListMap is a map of the cards struct.
var CardListMap map[string]cards = make(map[string]cards)

// LoadDefaultCards loads the default cards into memory on startup
func LoadDefaultCards() {
	// defaultCards is used to import the default cards
	defaultCards := CardListMap["default"]

	// Load default black cards.
	BCJson, err := ioutil.ReadFile("./cards/default/blackCards.json")

	if err != nil {
		fmt.Printf("Error loading black cards: %s", err)
		os.Exit(1)
	}

	err = json.Unmarshal(BCJson, &defaultCards.BlackCards)

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

	err = json.Unmarshal(WCJson, &defaultCards.WhiteCards)

	if err != nil {
		fmt.Printf("Error unmarshling black cards: %s", err)
		os.Exit(1)
	}

	// Update the map
	CardListMap["default"] = defaultCards
	loadCustomCards()
	loadIntoSlice()
}

func loadCustomCards() {
	// blackCardFiles is a string slice that holds the file names of the
	// json files for black cards
	blackCardFiles := []string{}

	// whiteCardFiles is a string slice that holds the file names of the
	// json files for white cards
	whiteCardFiles := []string{}

	// customCard holds the custom cards in a map.
	customCards := CardListMap["custom"]
	// Look for .json files in the black cards directory first.
	files, err := ioutil.ReadDir("./cards/custom/BlackCards/")
	if err != nil {
		fmt.Println("looks like there was a problem reading ./cards/custom/BlackCards/, does it exist?")
		return
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			blackCardFiles = append(blackCardFiles, file.Name())
		}
	}

	// Now let's look for json files in the white cards directory.
	files, err = ioutil.ReadDir("./cards/custom/WhiteCards/")
	if err != nil {
		fmt.Println("looks like there was a problem reading ./cards/custom/WhiteCards/, does it exist?")
		return
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			whiteCardFiles = append(whiteCardFiles, file.Name())
		}
	}

	// now let's load the black cards first
	tmpBlackcards := []blackCard{}
	for _, file := range blackCardFiles {
		Json, err := ioutil.ReadFile("./cards/custom/BlackCards/" + file)

		if err != nil {
			fmt.Printf("Error loading black cards: %s", err)
			continue
		}

		err = json.Unmarshal(Json, &customCards.BlackCards)

		if err != nil {
			fmt.Printf("Error unmarshling black cards: %s", err)
			continue
		}
		// Loop through and append to the tmpBlackcards slice
		for _, bc := range customCards.BlackCards {
			tmpBlackcards = append(tmpBlackcards, bc)
		}
	}

	// moving on to the white cardss
	tmpWhitecards := []whiteCard{}
	for _, file := range whiteCardFiles {
		Json, err := ioutil.ReadFile("./cards/custom/WhiteCards/" + file)

		if err != nil {
			fmt.Printf("Error loading white cards: %s", err)
			continue
		}

		err = json.Unmarshal(Json, &customCards.WhiteCards)

		if err != nil {
			fmt.Printf("Error unmarshling white cards: %s", err)
			continue
		}
		// Loop through and append to the tmpBlackcards slice
		for _, wc := range customCards.WhiteCards {
			tmpWhitecards = append(tmpWhitecards, wc)
		}
	}

	// Update the customCards map
	customCards.BlackCards = tmpBlackcards
	customCards.WhiteCards = tmpWhitecards
	// Now let's update the map
	CardListMap["custom"] = customCards

}

func loadIntoSlice() {
	// Loads the cards from the map into the CardList
	for _, card := range CardListMap {
		// load the black cards into CardList first.
		for _, blackcard := range card.BlackCards {
			CardList.BlackCards = append(CardList.BlackCards, blackcard)
		}
		// then we load the white cards.
		for _, whitecard := range card.WhiteCards {
			CardList.WhiteCards = append(CardList.WhiteCards, whitecard)
		}
	}
}
