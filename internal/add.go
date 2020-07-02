package internal

import (
	"fmt"
	util "github.com/Floor-Gang/utilpkg"
	dg "github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

// This is the add command
// args = [prefix, add, question ... \n .. answer]
func (b *Bot) Add(m *dg.MessageCreate, args []string) {
	// first let's make sure they provided enough info
	if len(args) < 3 {
		_, _ = util.Reply(b.client, m.Message,
			fmt.Sprintf("Command usage: `%s add <question> new-line <answer>`", b.config.Prefix),
		)
		return
	}

	split := strings.SplitN(
		strings.Join(args[2:], " "),
		"\n",
		2,
	)

	// Make sure they split their question and answer with one new-line
	if len(split) < 2 {
		_, _ = util.Reply(b.client, m.Message,
			fmt.Sprintf("Command usage: `%s add <question> new-line <answer>`", b.config.Prefix),
		)
		return
	}

	question := strings.TrimSpace(split[0])
	answer := split[1]
	guild := m.GuildID

	// Add the question + answer to the database
	err := b.db.Add(guild, question, answer)

	if err != nil {
		_, _ = util.Reply(b.client, m.Message,
			"Something went wrong while adding this FaQ to the database.")
		return
	}
	// Build the FaQ embed
	embed := buildEmbed(question, answer)
	// Send the embed to the FaQ channel
	_, err = b.client.ChannelMessageSendEmbed(b.config.FaQChannel, &embed)

	if err != nil {
		_, _ = util.Reply(b.client, m.Message, "Failed to send question to the FaQ channel, is one set?")
		log.Println("Failed to send question to the FaQ channel", question, answer, err)
	}
}
