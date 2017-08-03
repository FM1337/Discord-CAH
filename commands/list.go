package commands

import (
	"strings"

	"github.com/FM1337/Discord-CAH/game"
	"github.com/FM1337/Discord-CAH/utils"
	"github.com/bwmarrin/discordgo"
)

// Commands is a struct containing a slice of Command.
type Commands struct {
	CommandList []Command
}

// Command is a struct containing fields that hold command information.
type Command struct {
	Name    string                                                 // The name of the command.
	Command func(s *discordgo.Session, m *discordgo.MessageCreate) // The command function
	Admin   bool                                                   // Is the command admin only?
}

// CL is a Commands interface.
var CL Commands

func RegisterCommands() {
	CL.CommandList = append(CL.CommandList, Command{
		Name: "Start", Command: game.Start, Admin: false,
	})
	CL.CommandList = append(CL.CommandList, Command{
		Name: "Pause", Command: game.Pause, Admin: false,
	})
	CL.CommandList = append(CL.CommandList, Command{
		Name: "Stop", Command: game.Stop, Admin: false,
	})
	CL.CommandList = append(CL.CommandList, Command{
		Name: "Join", Command: game.Join, Admin: false,
	})
	CL.CommandList = append(CL.CommandList, Command{
		Name: "Leave", Command: game.Leave, Admin: false,
	})
	CL.CommandList = append(CL.CommandList, Command{
		Name: "Cards", Command: TotalCards, Admin: false,
	})
	CL.CommandList = append(CL.CommandList, Command{
		Name: "Pick", Command: game.PickCard, Admin: false,
	})
	CL.CommandList = append(CL.CommandList, Command{
		Name: "Choose", Command: game.ChooseWinner, Admin: false,
	})
	CL.CommandList = append(CL.CommandList, Command{
		Name: "Refresh", Command: Refresh, Admin: true,
	})
}

// RunCommand runs a specified command.
func RunCommand(name string, s *discordgo.Session, m *discordgo.MessageCreate) {
	for _, command := range CL.CommandList {
		if strings.ToLower(command.Name) == strings.ToLower(name) {
			// If the command is admin only.
			if command.Admin {
				// run is a temporary bool that will check if it can run
				// or not.
				run := false
				// Loop through the admin list
				for _, admin := range utils.Config.AdminIDs {
					if m.Author.ID == admin {
						run = true
						break
					}
				}
				// if not an admin, then don't continue
				if !run {
					return
				}

			}
			command.Command(s, m)
			return
		}
	}
}
