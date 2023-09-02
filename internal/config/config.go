package config

import (
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var Settings *SettingsRoot

type SettingsRoot struct {
	Api Api `yml:"api"`
}

type Api struct {
	Port string `yml:"port"`
	Gci  string `yml:"gci"`
}

func LoadSettings() error {
	f, err := os.Open(os.Getenv("CONFIG_FILE"))
	if err != nil {
		log.Err(err)
		return err
	}
	defer f.Close()

	return yaml.NewDecoder(f).Decode(&Settings)
}
