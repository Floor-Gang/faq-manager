package internal

import (
	"fmt"
	util "github.com/Floor-Gang/utilpkg"
	"github.com/bwmarrin/discordgo"
	"log"
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

// This is the get command
// args = [prefix, get, question context ... ]
func (b *Bot) Get(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 3 {
		_, _ = util.Reply(s, m.Message,
			fmt.Sprintf("Command usage: `%s get <question context>`", b.config.Prefix),
		)
		return
	}

	// Get the question they're talking about
	context := strings.ToLower(strings.Join(args[2:], " "))
	// Get all the FaQ for this Discord server
	faqs, err := b.db.GetAll(m.GuildID)

	if err != nil {
		_, _ = util.Reply(s, m.Message, "Something went wrong while getting the FaQ for this server.")
		log.Printf("Something went wrong while getting FaQ for %s\n%s\n", m.GuildID, err)
	}

	// Iterate through all the FaQ's for this server and see if any of them have
	// the context that the user gave us
	for _, faq := range faqs {
		// See if the question contains the context that the user gave us
		if strings.Contains(strings.ToLower(faq.Question), context) {
			// if it did then make an FaQ embed
			embed := buildEmbed(faq.Question, faq.Answer)
			// reply to the user with the FaQ embed
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, &embed)
			if err != nil {
				_, _ = util.Reply(s, m.Message, "Something went wrong while sending this FaQ to the FaQ channel.")
				log.Printf("Failed to send embed to %s\n%s\n", m.GuildID, err)
			}
			return
		}
	}
	_, _ = util.Reply(s, m.Message, fmt.Sprintf(`I couldn't find anything with "%s"`, context))
}

// This is the sync command
func (b *Bot) Sync(s *discordgo.Session, m *discordgo.MessageCreate) {
	messages, err := s.ChannelMessages(b.config.FaQChannel, 100, "", "", "")

	if err != nil {
		_, _ = util.Reply(s, m.Message, "Failed to get messages from the FaQ channel")
		return
	}

	// Go through all the messages in the FaQ channel
	// and delete them.
	for _, message := range messages {
		if message.Author.ID == s.State.User.ID {
			_ = s.ChannelMessageDelete(b.config.FaQChannel, message.ID)
		}
	}

	faqs, err := b.db.GetAll(m.GuildID)

	if err != nil {
		_, _ = util.Reply(s, m.Message, "Failed to get FaQ from the database for this server.")
	}

	// Get all the FaQ from the database for this Discord server
	// Post them into the FaQ channel.
	for _, faq := range faqs {
		embed := buildEmbed(faq.Question, faq.Answer)
		_, err = s.ChannelMessageSendEmbed(b.config.FaQChannel, &embed)
		if err != nil {
			log.Printf("Failed to send embed to %s\n%s\n", m.GuildID, err)
		}
	}
}
