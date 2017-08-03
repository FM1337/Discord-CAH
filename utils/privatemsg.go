package utils

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// AllowCommandPrivmsg will check to see if a command should be allowed to
// be executed in a private message rather than on the channel.
func AllowCommandPrivmsg(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	// open up a private channel
	channel, err := s.UserChannelCreate(m.Author.ID)

	if err != nil {
		return false
	}

	// Check if the IDs are the same
	if m.ChannelID == channel.ID {
		// For now only let the pick command be used.
		args := strings.Split(m.Content, " ")
		if strings.ToLower(args[0]) == "$pick" {
			return true
		}
	}

	return false
}
