package game

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/FM1337/Discord-CAH/cards"
	"github.com/FM1337/Discord-CAH/utils"
	"github.com/bwmarrin/discordgo"
)

type Player struct {
	PlayerID     string
	PlayerName   string
	PlayerNumber int
	Cards        []WhiteCard
	Points       int
	CardZar      bool
	PickedCards  []WhiteCard
}

type CardZar struct {
	PlayerNumber int
	PlayerID     string
}

type WhiteCard struct {
	CardID int
	Text   string
	Blank  bool
	taken  bool
}

type BlackCard struct {
	CardID int
	Text   string
	Cards  int
}

var Players map[string]Player = make(map[string]Player)
var Cardzars map[int]CardZar = make(map[int]CardZar)
var PlayerIDList []int

var RoundBlackCards map[int]BlackCard = make(map[int]BlackCard)
var RoundWhiteCards map[int]WhiteCard = make(map[int]WhiteCard)
var roundText string
var roundPlay int
var roundZar string

var running bool
var starting bool
var paused bool

var creatorID string
var pauserID string
var LastPlayerNumber int

func StartGame(s *discordgo.Session, m *discordgo.MessageCreate) {
	// TODO: Game code
	if running || starting {
		s.ChannelMessageSend(m.ChannelID, "A game is already running")
		return
	}

	starting = true
	creatorID = m.Author.ID
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s started a game!", m.Author.Username))
	Players[m.Author.ID] = Player{PlayerID: m.Author.ID, Cards: nil, Points: 0, PlayerNumber: 1, PlayerName: m.Author.Username}
	LastPlayerNumber = 1

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The game will start in 30 seconds, type %sjoin to join in!", utils.Config.Prefix))

	time.Sleep(30 * time.Second)
	fmt.Printf("%d players", len(Players))
	if len(Players) < 1 && starting {
		s.ChannelMessageSend(m.ChannelID, "Not enough players to start the game!")
		starting = false
		creatorID = ""
		Players = make(map[string]Player)
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Loading please wait...")
	starting = false
	running = true
	LoadTmpCards()
	GenerateHand()
	CardZarOrder()
	s.ChannelMessageSend(m.ChannelID, "The game has started!")
	Game(s, m)
	return
}

func PauseGame(s *discordgo.Session, m *discordgo.MessageCreate) {
	// TODO: Game code
	if !running {
		s.ChannelMessageSend(m.ChannelID, "No game is running")
		return
	}

	if !paused {
		paused = true
		s.ChannelMessageSend(m.ChannelID, "The game has been paused!")
		pauserID = m.Author.ID
		return
	}

	if m.Author.ID == pauserID || m.Author.ID == creatorID {
		paused = false
		s.ChannelMessageSend(m.ChannelID, "The game has been unpaused!")
		pauserID = ""
		return
	}

	pauser, _ := s.User(pauserID)
	pauserName := pauser.Username
	creator, _ := s.User(creatorID)
	creatorName := creator.Username

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Sorry only %s or %s can unpause the game!", pauserName, creatorName))
	return

}

// StopGame stops the game.
func StopGame(s *discordgo.Session, m *discordgo.MessageCreate) {
	// TODO: Game code
	if !running {
		if !starting {
			s.ChannelMessageSend(m.ChannelID, "No game is running")
			return
		}
	}

	for _, player := range Players {
		if m.Author.ID == player.PlayerID {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s stopped the game!", m.Author.Username))
			running = false
			starting = false
			Players = make(map[string]Player)
			return
		}
	}

	s.ChannelMessageSend(m.ChannelID, "You must be a player to stop the game!")
	return
}

func JoinGame(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !running {
		if !starting {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No game is running! Type %sstart to start a game", utils.Config.Prefix))
			return
		}

	}

	for _, player := range Players {
		if m.Author.ID == player.PlayerID {
			s.ChannelMessageSend(m.ChannelID, "You've already joined!")
			return
		}
	}
	Players[m.Author.ID] = Player{PlayerID: m.Author.ID, Cards: nil, Points: 0, PlayerNumber: LastPlayerNumber + 1, PlayerName: m.Author.Username}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s joined the game!", m.Author.Username))
	return
}

func LeaveGame(s *discordgo.Session, m *discordgo.MessageCreate) {

	if !running {
		if !starting {
			s.ChannelMessageSend(m.ChannelID, "No game is running")
			return
		}
	}

	delete(Players, m.Author.ID)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s has left the game!", m.Author.Username))

	if len(Players) < 3 && running {
		s.ChannelMessageSend(m.ChannelID, "Not enough players to continue!")
		running = false
		Players = make(map[string]Player)
		return
	}

	i := 1
	for _, player := range Players {
		player.PlayerNumber = i
		Players[player.PlayerID] = player
		LastPlayerNumber = i
		i = i + 1
	}
	CardZarOrder()
	return
}

func LoadTmpCards() {
	for i, card := range cards.CardList.BlackCards {
		RoundBlackCards[i] = BlackCard{
			CardID: i,
			Text:   card.CardText,
			Cards:  card.Cards2Play,
		}
	}
	for i, card := range cards.CardList.WhiteCards {
		RoundWhiteCards[i] = WhiteCard{
			CardID: i,
			Text:   card.CardText,
			taken:  false,
		}
	}
}

func Game(s *discordgo.Session, m *discordgo.MessageCreate) {
	round := 1
	cardzarNumber := 1
	for {
		if round > 10 || !running {
			break
		}
		if paused {
			time.Sleep(3 * time.Second)
			continue
		}
		roundZar = ChooseCardZar(cardzarNumber, true)
		seed := rand.NewSource(time.Now().UnixNano())
		random := rand.New(seed)
		roundCard := RoundBlackCards[random.Intn(len(RoundBlackCards))]
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Round %d, %s is the cardzar!", round, roundZar))
		roundPlay = roundCard.Cards
		if roundPlay == 0 {
			roundPlay = 1
		}
		roundText = roundCard.Text
		roundText = strings.Replace(roundText, "_", "______", -1)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s", roundText))
		time.Sleep(500 * time.Millisecond)
		MessageCards(s, m)
		s.ChannelMessageSend(m.ChannelID, "Players you have 30 seconds to choose a card/cards to play!")
		time.Sleep(30 * time.Second)
		delete(RoundBlackCards, roundCard.CardID)
		ChooseCardZar(cardzarNumber, false)
		round = round + 1
		cardzarNumber = cardzarNumber + 1
		if cardzarNumber > LastPlayerNumber {
			cardzarNumber = 1
		}

	}
	RoundBlackCards = make(map[int]BlackCard)
}

// GenerateHand will generate a player's hand when the game first starts.
func GenerateHand() {
	for _, player := range Players {
		var cardList []WhiteCard
		seed := rand.NewSource(time.Now().UnixNano())
		for i := 0; i != 10; i++ {
			var RandomCard WhiteCard
			for {
				random := rand.New(seed)
				RandomChoice := RoundWhiteCards[random.Intn(len(RoundWhiteCards))]
				if !RandomChoice.taken {
					RandomCard = RandomChoice
					RandomCard.taken = true
					RoundWhiteCards[RandomCard.CardID] = RandomCard
					break
				}
			}
			cardList = append(cardList, RandomCard)
		}
		player.Cards = cardList
		fmt.Printf("%v\n", player.Cards)
		Players[player.PlayerID] = player
	}
}

func CardsEmbed(PlayerCards []WhiteCard) *discordgo.MessageEmbed {

	fields := []*discordgo.MessageEmbedField{}
	for i, card := range PlayerCards {
		iString := strconv.Itoa(i + 1)
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   iString,
			Value:  card.Text,
			Inline: true,
		})
	}

	embed := &discordgo.MessageEmbed{
		Type:        "rich",
		Color:       52,
		Title:       fmt.Sprintf("Your cards (%d)", len(PlayerCards)),
		Description: roundText,
		Fields:      fields,
	}

	return embed
}

// MessageCards will message users their hand
func MessageCards(s *discordgo.Session, m *discordgo.MessageCreate) {
	for _, player := range Players {
		channel, err := s.UserChannelCreate(player.PlayerID)

		if err != nil {
			log.Fatal(err)
		}
		CEmbed := CardsEmbed(player.Cards)
		s.ChannelMessageSendEmbed(channel.ID, CEmbed)
		time.Sleep(1 * time.Second)

	}
}

func CardZarOrder() {
	Cardzars = make(map[int]CardZar)
	for _, player := range Players {
		Cardzars[player.PlayerNumber] = CardZar{
			PlayerNumber: player.PlayerNumber,
			PlayerID:     player.PlayerID,
		}
	}
}

func ChooseCardZar(CardzarNum int, current bool) string {
	PlayerID := Cardzars[CardzarNum].PlayerID
	tmpPlayer := Players[PlayerID]
	if current {
		tmpPlayer.CardZar = true
	} else {
		tmpPlayer.CardZar = false
	}
	Players[PlayerID] = tmpPlayer
	return tmpPlayer.PlayerName
}

func PickCard(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !running {
		s.ChannelMessageSend(m.ChannelID, "No game is running!")
		return
	}

	playing := false
	for _, player := range Players {
		if m.Author.ID == player.PlayerID {
			playing = true
			break
		}
	}

	if !playing {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You aren't playing! Type %sjoin to join", utils.Config.Prefix))
		return
	}

	if m.Author.Username == roundZar {
		s.ChannelMessageSend(m.ChannelID, "You're the CardZar, you don't play cards this round!")
		return
	}

	args := strings.Split(m.Content, " ")
	if len(args)-1 == roundPlay {
		tmpPlayer := Players[m.Author.ID]
		tmpPlayer.PickedCards = nil
		tmpRoundText := roundText
		for _, card := range args[1:] {

			cardNum, err := strconv.Atoi(card)

			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Invalid number!")
				fmt.Printf("%s\n", err)
				return
			}
			if strings.Contains(roundText, "______") {
				tmpPlayer.PickedCards = append(tmpPlayer.PickedCards, tmpPlayer.Cards[utils.IndexFixer(len(tmpPlayer.Cards), cardNum)])

				tmpRoundText = strings.Replace(tmpRoundText, "______", tmpPlayer.Cards[utils.IndexFixer(len(tmpPlayer.Cards), cardNum)].Text, 1)
			} else {
				tmpRoundText = fmt.Sprintf("%s %s", tmpRoundText, tmpPlayer.Cards[utils.IndexFixer(len(tmpPlayer.Cards), cardNum)].Text)
			}

		}
		Players[m.Author.ID] = tmpPlayer

		channel, err := s.UserChannelCreate(m.Author.ID)
		if err != nil {
			log.Fatal(err)
		}
		s.ChannelMessageSend(channel.ID, fmt.Sprintf("You played: %s", tmpRoundText))
		return

	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You've chosen an incorrect amount of cards to play!, The amount you played was %d when you should've played %d", len(args)-1, roundPlay))
}
