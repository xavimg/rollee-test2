package main

import (
	"time"

	"words/internal/handlers"
	"words/internal/service"
	"words/internal/storage"
	"words/server"

	"github.com/rs/zerolog/log"
)

const (
	addr string = ":3001" // Should be in .env variable, but because of test purposes im remaining as a constant.

	// gci defines the frequency at which the cleanGarbageCollector() operates.
	// This ensures efficient memory management and reduces the risk of memory leaks.
	// In the context of high-throughput scenarios, this interval might be set even shorter
	// because we dont want that our in-memory storage grows a lot.
	gci time.Duration = 3 * time.Second
)

func main() {
	storage := storage.NewInMemoryStore(gci)
	service := service.NewWordService(storage)
	handlers := handlers.NewHandler(service)

	// Due to NewHTTP accepting interfaces for storage and service, it offers the flexibility
	// to easily switch between different implementations. For instance, we can effortlessly swap our
	// InMemoryStore with a PostgreSQLStore or even replace the WordService with another service implementation,
	// ensuring a high level of decoupling and adaptability.
	h := server.NewHTTP(addr, gci, storage, service, handlers)

	if err := h.Start(); err != nil {
		log.Fatal().Msgf("Failed to start the server: %v", err)
	}
}
