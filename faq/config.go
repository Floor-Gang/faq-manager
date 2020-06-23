package faq

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
)

type Config struct {
	Prefix     string
	Token      string
	Auth       string
	FaQChannel string
	DBLocation string
}

func getConfig(location string) Config {
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
