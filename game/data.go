package game

// Player is a struct that holds a player's data.
type Player struct {
	PlayerName string
	PlayerID   string
	Zar        bool
	Cards      []WhiteCard
}

// WhiteCard is a struct that holds the data about a white card.
type WhiteCard struct {
	CardID int
	Index  int
	Text   string
	Blank  bool
	taken  bool
}

// BlackCard is a struct that holds the data about a black card.
type BlackCard struct {
	CardID int
	Text   string
	Cards  int
}
