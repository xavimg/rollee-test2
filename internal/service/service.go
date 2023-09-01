package service

import (
	"context"
	"time"
	"words/internal/http"
	"words/internal/http/handlers"
	"words/internal/service/word"
	"words/internal/storage/in_memory"
)

// gci defines the frequency at which the cleanGarbageCollector() operates.
// This ensures efficient memory management due to the potentially large volume of storage.
// In the context of high-throughput scenarios, this interval might be set even shorter
// because we dont want that our in-memory storage grows a lot.
const gci time.Duration = 20 * time.Second

type App struct {
	ctx       context.Context
	ctxCancel context.CancelFunc

	port string
}

func NewApp(ctx context.Context, add string) *App {
	ctx, ctxCancel := context.WithCancel(ctx)
	return &App{
		ctx:       ctx,
		ctxCancel: ctxCancel,
		port:      add,
	}
}

func (s *App) Run() {
	repository := in_memory.NewInMemoryStorage(3 * time.Second)
	service := word.NewWordService(repository)
	handlers := handlers.NewHandler(service)

	// addr i gci entorn variables
	h := http.NewHTTP(s.ctx, gci, repository, handlers)

	go func() {
		if err := h.Start(); err != nil {
			panic(err)
		}
	}()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			go repository.CleanGarbageCollector()
		}
	}

}
