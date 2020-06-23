package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
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

// This will remove a question from the FaQ channel (it will seek up to 100 messages)
func (b *Bot) removeFromChannel(s *discordgo.Session, m *discordgo.MessageCreate, question string) {
	messages, err := s.ChannelMessages(b.config.FaQChannel, 100, "", "", "")

	if err != nil {
		reply(s, m, "Failed to get messages from FaQ channel to remove this question.")
		log.Println("Failed getting messages from faq channel", err)
		return
	}

	for _, message := range messages {
		if len(message.Embeds) > 0 {
			embed := message.Embeds[0]
			title := strings.ToLower(embed.Title)

			if title == strings.ToLower(question) {
				err := s.ChannelMessageDelete(b.config.FaQChannel, message.ID)

				if err != nil {
					reply(s, m, "Failed to delete the that question from the FaQ channel.")
					log.Println("Failed to delete question from FaQ channel.", err)
				}
				break
			}
		} else {
			log.Println("no")
		}
	}
}
