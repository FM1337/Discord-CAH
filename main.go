package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/FM1337/Discord-CAH/cards"
	"github.com/FM1337/Discord-CAH/commands"
	"github.com/FM1337/Discord-CAH/utils"
	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	utils.Config.LoadConfig()
	commands.RegisterCommands()
	cards.LoadDefaultCards()
	// Generate a random seed on startup.
	rand.Seed(time.Now().UnixNano())
	discord()
}

func discord() {
	discordToken := utils.Config.DiscordToken
	fmt.Println("Starting bot..")
	bot, err := discordgo.New("Bot " + discordToken)

	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		log.Fatal(err)
	}

	err = bot.Open()

	if err != nil {
		log.Fatal(err)
	}

	bot.AddHandler(messageCreate)
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	// Cleanly close down the Discord session.
	bot.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Print("[" + m.Author.Username + "] " + m.Content)

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.ChannelID == utils.Config.CAHChannelID || utils.AllowCommandPrivmsg(s, m) {
		if strings.HasPrefix(m.Content, utils.Config.Prefix) {
			command := strings.Split(m.Content[1:len(m.Content)], " ")
			name := strings.ToLower(command[0])
			commands.RunCommand(name, s, m)
			return
		}

		return
	}
}
