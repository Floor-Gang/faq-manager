package faq

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

// This will create a new FaQ embed to be sent into the FaQ channel
func buildEmbed(question string, answer string) discordgo.MessageEmbed {
	embed := discordgo.MessageEmbed{}
	embed.Color = 0x1385ef
	embed.Title = question
	embed.Description = answer

	return embed
}

// This checks whether or not a user has an bot admin role. if they don't then
// they will be responded with "you must be an admin to use this command."
func (b *Bot) adminCheck(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if auth, err := b.auth.Auth(m.Author.ID); err != nil {
		reply(s, m, "Something went wrong while authenticating you.")
		return false
	} else {
		if !auth.IsAdmin {
			log.Println(auth)
			reply(s, m, "You must be admin to use this command.")
		}
		return auth.IsAdmin
	}
}

// This is a reply function which just mentions the with the given context
// (ie) @dylan Hello!
func reply(s *discordgo.Session, m *discordgo.MessageCreate, context string) {
	_, err := s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf("<@%s> %s", m.Author.ID, context),
	)

	if err != nil {
		log.Printf(
			"Failed to respond to %s#%s in %s about %s\n%s\n",
			m.Author.Username, m.Author.Discriminator,
			m.ChannelID,
			context,
			err,
		)
	}
}
