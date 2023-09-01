package main

import (
	"words/internal/config"
	"words/internal/service"
	"words/pkg/context"

	"github.com/rs/zerolog/log"
)

// cargar variables entorno
func init() {
	if err := config.LoadSettings(); err != nil {
		log.Fatal().Err(err)
	}
}

func main() {
	ctx := context.ApplicationContext()

	// struct

	//pasarle el init
	app := service.NewApp(ctx, config.Settings.Api.Port)
	app.Run()

	// repository := in_memory.NewInMemoryStorage(gci)
	// service := service.NewWordService(repository)
	// handlers := handlers.NewHandler(service)

	// h := http.NewHTTP(ctx, addr, gci, repository, handlers)

	// if err := h.Start(); err != nil {
	// 	panic(err)
	// }
}
