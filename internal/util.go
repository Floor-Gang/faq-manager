package internal

import (
	util "github.com/Floor-Gang/utilpkg/botutil"
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

// This will remove a question from the FaQ channel (it will seek up to 100 messages)
func (b *Bot) RemoveFromChannel(s *discordgo.Session, m *discordgo.MessageCreate, question string) {
	messages, err := s.ChannelMessages(b.config.FaQChannel, 100, "", "", "")

	if err != nil {
		_, _ = util.Reply(s, m.Message,
			"Failed to get messages from FaQ channel to remove this question.")
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
					_, _ = util.Reply(s, m.Message,
						"Failed to delete the that question from the FaQ channel.")
					log.Println("Failed to delete question from FaQ channel.", err)
				}
				break
			}
		} else {
			log.Println("no")
		}
	}
}
