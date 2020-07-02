package internal

import (
	util "github.com/Floor-Gang/utilpkg"
	dg "github.com/bwmarrin/discordgo"
	"log"
)

// This is the list command
func (b *Bot) List(s *dg.Session, m *dg.MessageCreate) {
	// Get all the FaQ's for this Discord server from the database
	faqs, err := b.db.GetAll(m.GuildID)

	if err != nil {
		_, _ = util.Reply(s, m.Message,
			"Something went wrong while getting the FaQ for this server.")
		log.Printf("Something went wrong while getting FaQ for %s\n%s\n", m.GuildID, err)
	}

	// Iterate through all the FaQ, put them in embeds, and send it to the channel
	for _, faq := range faqs {
		embed := buildEmbed(faq.Question, faq.Answer)
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		if err != nil {
			log.Printf("Failed to send embed to %s\n%s\n", m.GuildID, err)
		}
	}
}
