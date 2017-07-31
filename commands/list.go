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
		Name: "Start", Command: game.StartGame,
	})
	CL.CommandList = append(CL.CommandList, Command{
		Name: "Pause", Command: game.PauseGame,
	})
	CL.CommandList = append(CL.CommandList, Command{
		Name: "Stop", Command: game.StopGame,
	})
	CL.CommandList = append(CL.CommandList, Command{
		Name: "Join", Command: game.JoinGame,
	})
	CL.CommandList = append(CL.CommandList, Command{
		Name: "Leave", Command: game.LeaveGame,
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
