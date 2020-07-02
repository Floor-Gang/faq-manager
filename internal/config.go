package internal

import (
	util "github.com/Floor-Gang/utilpkg"
	"log"
	"strings"
)

type Config struct {
	Prefix     string `yaml:"prefix"`
	Token      string `yaml:"token"`
	Auth       string `yaml:"auth_server"`
	FaQChannel string `yaml:"channel"`
	DBLocation string `yaml:"db_location"`
}

// This will get the current configuration file. If it doesn't exist then a
// new one will be made.
func getConfig(location string) Config {
	config := Config{
		Prefix:     ".faq",
		Token:      "",
		Auth:       "localhost:6969",
		FaQChannel: "",
		DBLocation: "",
	}
	err := util.GetConfig(location, &config)

	if err != nil {
		if strings.Contains(err.Error(), "default") {
			log.Fatalln("A default configuration has been made.")
		} else {
			log.Fatalln("Something went wrong while generating config. " + err.Error())
		}
	}

	return config
}
