package faq

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

func (b *Bot) onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.GuildID) == 0 || !strings.HasPrefix(m.Content, b.config.Prefix) {
		return
	}
	// args = [<prefix>, command]
	args := strings.Split(m.Content, " ")

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
		if b.adminCheck(s, m) {
			b.get(s, m, args)
		}
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

// args = [prefix, add, question ... answer]
func (b *Bot) add(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
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

	if len(split) < 2 {
		reply(s, m,
			fmt.Sprintf("Command usage: `%s add <question> new-line <answer>`", b.config.Prefix),
		)
		return
	}

	question := strings.TrimSpace(split[0])
	answer := split[1]
	guild := m.GuildID

	err := b.db.Add(guild, question, answer)

	if err != nil {
		reply(s, m, "Something went wrong while adding this FaQ to the database.")
		return
	}
	embed := buildEmbed(question, answer)
	_, err = s.ChannelMessageSendEmbed(b.config.FaQChannel, &embed)

	if err != nil {
		reply(s, m, "Failed to send question to the FaQ channel, is one set?")
		log.Println("Failed to send question to the FaQ channel", question, answer, err)
	}
}

// args = [prefix, get, question context ... ]
func (b *Bot) get(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 3 {
		reply(s, m,
			fmt.Sprintf("Command usage: `%s get <question context>`", b.config.Prefix),
		)
		return
	}

	context := strings.ToLower(strings.Join(args[2:], " "))
	faqs, err := b.db.GetAll(m.GuildID)

	if err != nil {
		reply(s, m, "Something went wrong while getting the FaQ for this server.")
		log.Printf("Something went wrong while getting FaQ for %s\n%s\n", m.GuildID, err)
	}

	for _, faq := range faqs {
		if strings.Contains(strings.ToLower(faq.Question), context) {
			embed := buildEmbed(faq.Question, faq.Answer)
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

func (b *Bot) list(s *discordgo.Session, m *discordgo.MessageCreate) {
	faqs, err := b.db.GetAll(m.GuildID)

	if err != nil {
		reply(s, m, "Something went wrong while getting the FaQ for this server.")
		log.Printf("Something went wrong while getting FaQ for %s\n%s\n", m.GuildID, err)
	}

	for _, faq := range faqs {
		embed := buildEmbed(faq.Question, faq.Answer)
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		if err != nil {
			log.Printf("Failed to send embed to %s\n%s\n", m.GuildID, err)
		}
	}
}

func (b *Bot) remove(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 3 {
		reply(s, m,
			fmt.Sprintf("Command usage: `%s remove <question context>`", b.config.Prefix),
		)
		return
	}

	context := strings.ToLower(strings.Join(args[2:], " "))
	faqs, err := b.db.GetAll(m.GuildID)

	if err != nil {
		reply(s, m, "Something went wrong while getting the FaQ for this server.")
		log.Printf("Something went wrong while getting FaQ for %s\n%s\n", m.GuildID, err)
	}

	for _, faq := range faqs {
		if strings.Contains(strings.ToLower(faq.Question), context) {
			err = b.db.Remove(m.GuildID, faq.Question)

			if err != nil {
				reply(s, m, "Something went wrong while removing the question")
				log.Printf(`Something went wrong while removing the following "%s" for "%s"`, faq.Question, m.GuildID)
			} else {
				reply(s, m, "Removed "+faq.Question)
				b.removeFromChannel(s, m, faq.Question)
			}
			return
		}
	}
	reply(s, m, "Couldn't find \""+context+"\"")
}

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

func (b *Bot) sync(s *discordgo.Session, m *discordgo.MessageCreate) {
	messages, err := s.ChannelMessages(b.config.FaQChannel, 100, "", "", "")

	if err != nil {
		reply(s, m, "Failed to get messages from the FaQ channel")
		return
	}

	for _, message := range messages {
		if message.Author.ID == s.State.User.ID {
			_ = s.ChannelMessageDelete(b.config.FaQChannel, message.ID)
		}
	}

	faqs, err := b.db.GetAll(m.GuildID)

	if err != nil {
		reply(s, m, "Failed to get FaQ from the database for this server.")
	}

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
