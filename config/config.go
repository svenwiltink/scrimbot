package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	DiscordToken string `json:"discordToken"`
	BoltDatabase string `json:"boltDatabase"`
}

func LoadConfig(path string) (*Config, error) {
	contents, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	config := &Config{}

	if err = json.Unmarshal(contents, config); err != nil {
		return nil, err
	}

	return config, err
}
