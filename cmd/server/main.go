package main

import (
	"words/internal/config"
	"words/internal/service"
	"words/pkg/context"

	"github.com/rs/zerolog/log"
)

func init() {
	if err := config.LoadSettings(); err != nil {
		log.Fatal().Err(err)
	}
}

func main() {
	ctx := context.ApplicationContext()

	app := service.NewApp(ctx)
	app.Run()
}
