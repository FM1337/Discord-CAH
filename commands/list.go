package commands

import (
	"strings"

	"github.com/FM1337/Discord-CAH/game"
	"github.com/bwmarrin/discordgo"
)

// Commands is a struct containing a slice of Command.
type Commands struct {
	CommandList []Command
}

// Command is a struct containing fields that hold command information.
type Command struct {
	Name    string
	Command func(s *discordgo.Session, m *discordgo.MessageCreate)
}

// CL is a Commands interface.
var CL Commands

func RegisterCommands() {
	CL.CommandList = append(CL.CommandList, Command{
		Name: "Start", Command: game.Start,
	})
	CL.CommandList = append(CL.CommandList, Command{
		Name: "Pause", Command: game.Pause,
	})
	CL.CommandList = append(CL.CommandList, Command{
		Name: "Stop", Command: game.Stop,
	})
	CL.CommandList = append(CL.CommandList, Command{
		Name: "Join", Command: game.Join,
	})
	CL.CommandList = append(CL.CommandList, Command{
		Name: "Leave", Command: game.Leave,
	})
	CL.CommandList = append(CL.CommandList, Command{
		Name: "Cards", Command: TotalCards,
	})
	CL.CommandList = append(CL.CommandList, Command{
		Name: "Pick", Command: game.PickCard,
	})
	CL.CommandList = append(CL.CommandList, Command{
		Name: "Choose", Command: game.ChooseWinner,
	})
}

// RunCommand runs a specified command.
func RunCommand(name string, s *discordgo.Session, m *discordgo.MessageCreate) {
	for _, command := range CL.CommandList {
		if strings.ToLower(command.Name) == strings.ToLower(name) {
			command.Command(s, m)
			return
		}
	}
}
