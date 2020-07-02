package internal

import (
	auth "github.com/Floor-Gang/authclient"
	"github.com/bwmarrin/discordgo"
	"log"
)

type Bot struct {
	config Config
	db     Controller
	auth   auth.AuthClient
	client *discordgo.Session
}

func Start(configPath string) {
	config := getConfig(configPath)
	db := GetController(config.DBLocation)
	authClient, err := auth.GetClient(config.Auth)

	if err != nil {
		log.Fatalln("Failed to connect to authentication server")
	}

	register, err := authClient.Register(
		auth.Feature{
			Name:          "FaQ Manager",
			Description:   "For managing the FaQ channel",
			CommandPrefix: config.Prefix,
			Commands: []auth.SubCommand{
				{
					Name:        "get",
					Description: "Get a question based on provided context",
					Example:     []string{"get", "question context"},
				},
				{
					Name:        "set",
					Description: "Set a question's answer",
					Example:     []string{"set", "question context NEWLINE", "...new answer..."},
				},
				{
					Name:        "add",
					Description: "Add a new question",
					Example:     []string{"add", "question NEWLINE", "...new answer..."},
				},
				{
					Name:        "list",
					Description: "List all the currently stored FaQ embeds",
					Example:     []string{"list"},
				},
				{
					Name:        "remove",
					Description: "Remove a question",
					Example:     []string{"remove", "question context"},
				},
				{
					Name:        "sync",
					Description: "Sync the FaQ channel with the stored FaQ's",
					Example:     []string{"sync"},
				},
			},
		},
	)

	if err != nil {
		panic(err)
	}

	client, _ := discordgo.New(register.Token)

	bot := Bot{
		config: config,
		db:     db,
		auth:   authClient,
		client: client,
	}

	client.AddHandler(bot.onMessage)

	if err = client.Open(); err != nil {
		panic(err)
	}
}
