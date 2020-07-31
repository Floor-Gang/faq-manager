package internal

import (
	"fmt"
	util "github.com/Floor-Gang/utilpkg/botutil"
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

// This is the get command
// args = [prefix, get, question context ... ]
func (b *Bot) Get(s *dg.Session, m *dg.MessageCreate, args []string) {
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

// This is the set command
// args = [prefix, set, question ... new answer]
func (b *Bot) Set(s *dg.Session, m *dg.MessageCreate, args []string) {
	if len(args) < 3 {
		_, _ = util.Reply(s, m.Message,
			fmt.Sprintf("Command usage: `%s set <question> new-line <answer>`", b.config.Prefix),
		)
		return
	}

	split := strings.SplitN(
		strings.Join(args[2:], " "),
		"\n",
		2,
	)

	if len(split) < 2 {
		_, _ = util.Reply(s, m.Message,
			fmt.Sprintf("Command usage: `%s set <question> new-line <answer>`", b.config.Prefix),
		)
		return
	}

	context := strings.ToLower(split[0])
	newAnswer := split[1]
	faqs, err := b.db.GetAll(m.GuildID)

	if err != nil {
		_, _ = util.Reply(s, m.Message,
			"Something went wrong while getting the FaQ for this server.")
		log.Printf("Something went wrong while getting FaQ for %s\n%s\n", m.GuildID, err)
	}

	for _, faq := range faqs {
		question := strings.ToLower(faq.Question)

		if strings.Contains(question, context) {
			err := b.db.Set(m.GuildID, faq.Question, newAnswer)

			if err != nil {
				_, _ = util.Reply(s, m.Message,
					"Something went wrong while updating this question in the database")
				log.Println(err)
				return
			}
			b.RemoveFromChannel(s, m, faq.Question)
			embed := buildEmbed(faq.Question, newAnswer)
			_, err = s.ChannelMessageSendEmbed(b.config.FaQChannel, &embed)
			if err != nil {
				_, _ = util.Reply(s, m.Message,
					"Something went while sending the new FaQ to the FaQ channel.")
				log.Printf("Failed to send embed to %s\n%s\n", m.GuildID, err)
			}

			_, _ = util.Reply(s, m.Message, "Updated.")
			return
		}
	}
	_, _ = util.Reply(s, m.Message, fmt.Sprintf(`Couldn't find "%s"'`, context))
}

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

// This is the sync command
func (b *Bot) Sync(s *dg.Session, m *dg.MessageCreate) {
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
