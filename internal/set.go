package internal

import (
	"fmt"
	util "github.com/Floor-Gang/utilpkg"
	dg "github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

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
