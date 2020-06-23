package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/floor-gang/auth-client"
)

type Bot struct {
	config Config
	db     Controller
	auth   auth_client.AuthClient
}

func Start(configLocation string) {
	config := getConfig(configLocation)
	db := getController(config.DBLocation)
	client, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		panic(err)
	}

	authClient, err := auth_client.GetClient(config.Auth)

	bot := Bot{
		config: config,
		db:     db,
		auth:   authClient,
	}

	if err != nil {
		panic(err)
	}

	client.AddHandler(bot.onMessage)

	if err = client.Open(); err != nil {
		panic(err)
	}
}
