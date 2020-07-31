package internal

import (
	util "github.com/Floor-Gang/utilpkg/botutil"
	"github.com/bwmarrin/discordgo"
	"strings"
)

// This listens for all new messages that the bot can see
func (b *Bot) onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages that aren't in a guild and don't start with the prefix
	if len(m.GuildID) == 0 || !strings.HasPrefix(m.Content, b.config.Prefix) {
		return
	}
	// args = [<prefix>, command]
	args := strings.Split(m.Content, " ")

	// Make sure they provided a command (ie .faq add)
	if len(args) < 2 {
		return
	}

	response, _ := b.auth.Auth(m.Author.ID)
	isAdmin := response.IsAdmin

	if isAdmin {
		switch args[1] {
		case "add":
			b.Add(m, args)
			break
		case "get":
			b.Get(s, m, args)
			break
		case "list":
			b.List(s, m)
			break
		case "remove":
			b.Remove(s, m, args)
			break
		case "set":
			b.Set(s, m, args)
			break
		case "sync":
			b.Sync(s, m)
			break
		}
	} else {
		_, _ = util.Reply(s, m.Message, "You must be admin to run this command.")
	}
}
