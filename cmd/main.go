package main

import (
	"time"
	"words/internal/handlers"
	"words/internal/service"
	"words/internal/store"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

const (
	addr string = ":3001"

	// gci defines the frequency at which the cleanGarbageCollector() operates.
	// This ensures efficient memory management and reduces the risk of memory leaks.
	// In the context of high-throughput scenarios, this interval might be set even shorter
	// because we dont want that our in-memory storage grows a lot.
	gci time.Duration = 30 * time.Second
)

func main() {
	store := store.NewInMemoryStore(gci)
	service := service.NewWordService(store)
	handler := handlers.NewHandler(service)

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Post("/api/v0.1/words/{word}", handler.AddWord)
		r.Get("/api/v0.1/words/{prefix}", handler.FrequentWordByPrefix)
	})

	server := NewServer(addr, r)
	if err := server.Start(); err != nil {
		log.Fatal().Msgf("Failed to start the server: %v", err)
	}

}
