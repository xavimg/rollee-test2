package config

import (
	"fmt"
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
}

func LoadSettings() error {
	f, err := os.Open(os.Getenv("CONFIG_FILE"))
	if err != nil {
		fmt.Println("SSSS")
		log.Err(err)
		return err
	}
	defer f.Close()

	return yaml.NewDecoder(f).Decode(&Settings)
}
