package faq

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
)

type Config struct {
	Prefix     string
	Token      string
	Auth       string // Authentication RPC server address
	FaQChannel string // FaQ channel ID
	DBLocation string // Database absolute location
}

// This will get the current configuration file. If it doesn't exist then a
// new one will be made.
func getConfig(location string) Config {
	// This will attempt to read from the config file. If it doesn't exist then
	// a new configuration file will be generated otherwise we continue to parse
	// the current config and return it
	if data, err := ioutil.ReadFile(location); err != nil {
		return genConfig(location)
	} else {
		configMap := make(map[string]interface{})
		err := yaml.Unmarshal(data, configMap)
		config := Config{
			Prefix:     "",
			Token:      "",
			Auth:       "",
			FaQChannel: "",
			DBLocation: "",
		}

		if err != nil {
			panic(err)
		}

		config.Prefix = configMap["prefix"].(string)
		config.Token = configMap["token"].(string)
		config.Auth = configMap["auth"].(string)
		config.FaQChannel = configMap["faq_channel"].(string)
		config.DBLocation = configMap["db_location"].(string)

		return config
	}
}

// This will create a new configuration file
func genConfig(location string) Config {
	newConfig := map[string]string{
		"prefix":      ".faq",
		"token":       "",
		"auth":        "127.0.0.1:6969",
		"faq_channel": "",
		"db_location": "./faq.db",
	}

	serialized, err := yaml.Marshal(newConfig)

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(location, serialized, 0660)

	if err != nil {
		panic(err)
	}

	return getConfig(location)
}
