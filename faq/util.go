package faq

import (
	"github.com/bwmarrin/discordgo"
)

func buildEmbed(question string, answer string) discordgo.MessageEmbed {
	embed := discordgo.MessageEmbed{}
	embed.Color = 0x1385ef
	embed.Title = question
	embed.Description = answer

	return embed
}
