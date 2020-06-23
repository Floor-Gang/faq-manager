package faq

import (
	"fmt"
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

	switch args[1] {
	case "add":
		if b.adminCheck(s, m) {
			b.add(s, m, args)
		}
		break
	case "get":
		b.get(s, m, args)
		break
	case "list":
		if b.adminCheck(s, m) {
			b.list(s, m)
		}
		break
	case "remove":
		if b.adminCheck(s, m) {
			b.remove(s, m, args)
		}
		break
	case "set":
		if b.adminCheck(s, m) {
			b.set(s, m, args)
		}
		break
	case "sync":
		if b.adminCheck(s, m) {
			b.sync(s, m)
		}
		break
	}
}

// This is the add command
// args = [prefix, add, question ... \n .. answer]
func (b *Bot) add(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	// first let's make sure they provided enough info
	if len(args) < 3 {
		reply(s, m,
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
		reply(s, m,
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
		reply(s, m, "Something went wrong while adding this FaQ to the database.")
		return
	}
	// Build the FaQ embed
	embed := buildEmbed(question, answer)
	// Send the embed to the FaQ channel
	_, err = s.ChannelMessageSendEmbed(b.config.FaQChannel, &embed)

	if err != nil {
		reply(s, m, "Failed to send question to the FaQ channel, is one set?")
		log.Println("Failed to send question to the FaQ channel", question, answer, err)
	}
}

// This is the get command
// args = [prefix, get, question context ... ]
func (b *Bot) get(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 3 {
		reply(s, m,
			fmt.Sprintf("Command usage: `%s get <question context>`", b.config.Prefix),
		)
		return
	}

	// Get the question they're talking about
	context := strings.ToLower(strings.Join(args[2:], " "))
	// Get all the FaQ for this Discord server
	faqs, err := b.db.GetAll(m.GuildID)

	if err != nil {
		reply(s, m, "Something went wrong while getting the FaQ for this server.")
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
				reply(s, m, "Something went wrong while sending this FaQ to the FaQ channel.")
				log.Printf("Failed to send embed to %s\n%s\n", m.GuildID, err)
			}
			return
		}
	}
	reply(s, m, fmt.Sprintf(`I couldn't find anything with "%s"`, context))
}

// This is the list command
func (b *Bot) list(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Get all the FaQ's for this Discord server from the database
	faqs, err := b.db.GetAll(m.GuildID)

	if err != nil {
		reply(s, m, "Something went wrong while getting the FaQ for this server.")
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

// This is the remove command
// args = [prefix, remove, question context]
func (b *Bot) remove(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 3 {
		reply(s, m,
			fmt.Sprintf("Command usage: `%s remove <question context>`", b.config.Prefix),
		)
		return
	}

	// Get the question they're talking about
	context := strings.ToLower(strings.Join(args[2:], " "))
	faqs, err := b.db.GetAll(m.GuildID)

	if err != nil {
		reply(s, m, "Something went wrong while getting the FaQ for this server.")
		log.Printf("Something went wrong while getting FaQ for %s\n%s\n", m.GuildID, err)
	}

	// Iterate through all the FaQ's for this Discord server
	for _, faq := range faqs {
		if strings.Contains(strings.ToLower(faq.Question), context) {
			// Remove it from the database
			err = b.db.Remove(m.GuildID, faq.Question)

			if err != nil {
				reply(s, m, "Something went wrong while removing the question")
				log.Printf(`Something went wrong while removing the following "%s" for "%s"`, faq.Question, m.GuildID)
			} else {
				reply(s, m, "Removed "+faq.Question)
				// Remove it from the FaQ channel (only up to 100 messages)
				b.removeFromChannel(s, m, faq.Question)
			}
			return
		}
	}
	reply(s, m, "Couldn't find \""+context+"\"")
}

// This is the set command
// args = [prefix, set, question ... new answer]
func (b *Bot) set(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 3 {
		reply(s, m,
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
		reply(s, m,
			fmt.Sprintf("Command usage: `%s set <question> new-line <answer>`", b.config.Prefix),
		)
		return
	}

	context := strings.ToLower(split[0])
	newAnswer := split[1]
	faqs, err := b.db.GetAll(m.GuildID)

	if err != nil {
		reply(s, m, "Something went wrong while getting the FaQ for this server.")
		log.Printf("Something went wrong while getting FaQ for %s\n%s\n", m.GuildID, err)
	}

	for _, faq := range faqs {
		question := strings.ToLower(faq.Question)

		if strings.Contains(question, context) {
			err := b.db.Set(m.GuildID, faq.Question, newAnswer)

			if err != nil {
				reply(s, m, "Something went wrong while updating this question in the database")
				log.Println(err)
				return
			}
			b.removeFromChannel(s, m, faq.Question)
			embed := buildEmbed(faq.Question, newAnswer)
			_, err = s.ChannelMessageSendEmbed(b.config.FaQChannel, &embed)
			if err != nil {
				reply(s, m, "Something went while sending the new FaQ to the FaQ channel.")
				log.Printf("Failed to send embed to %s\n%s\n", m.GuildID, err)
			}

			reply(s, m, "Updated.")
			return
		}
	}
	reply(s, m, fmt.Sprintf(`Couldn't find "%s"'`, context))
}

// This is the sync command
func (b *Bot) sync(s *discordgo.Session, m *discordgo.MessageCreate) {
	messages, err := s.ChannelMessages(b.config.FaQChannel, 100, "", "", "")

	if err != nil {
		reply(s, m, "Failed to get messages from the FaQ channel")
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
		reply(s, m, "Failed to get FaQ from the database for this server.")
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
