package internal

import (
	"fmt"
	util "github.com/Floor-Gang/utilpkg"
	dg "github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

// This is the remove command
// args = [prefix, remove, question context]
func (b *Bot) Remove(s *dg.Session, m *dg.MessageCreate, args []string) {
	if len(args) < 3 {
		_, _ = util.Reply(s, m.Message,
			fmt.Sprintf("Command usage: `%s remove <question context>`", b.config.Prefix),
		)
		return
	}

	// Get the question they're talking about
	context := strings.ToLower(strings.Join(args[2:], " "))
	faqs, err := b.db.GetAll(m.GuildID)

	if err != nil {
		_, _ = util.Reply(s, m.Message,
			"Something went wrong while getting the FaQ for this server.")
		log.Printf("Something went wrong while getting FaQ for %s\n%s\n", m.GuildID, err)
	}

	// Iterate through all the FaQ's for this Discord server
	for _, faq := range faqs {
		if strings.Contains(strings.ToLower(faq.Question), context) {
			// Remove it from the database
			err = b.db.Remove(m.GuildID, faq.Question)

			if err != nil {
				_, _ = util.Reply(s, m.Message, "Something went wrong while removing the question")
				log.Printf(`Something went wrong while removing the following "%s" for "%s"`, faq.Question, m.GuildID)
			} else {
				_, _ = util.Reply(s, m.Message, "Removed "+faq.Question)
				// Remove it from the FaQ channel (only up to 100 messages)
				b.RemoveFromChannel(s, m, faq.Question)
			}
			return
		}
	}
	_, _ = util.Reply(s, m.Message, "Couldn't find \""+context+"\"")
}
