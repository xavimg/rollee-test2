package service

import (
	"context"
	"fmt"
	"time"
	config "words/internal/config"
	"words/internal/http"
	"words/internal/http/handlers"
	"words/internal/service/word"
	"words/internal/storage/in_memory"
)

type App struct {
	ctx       context.Context
	ctxCancel context.CancelFunc

	config *config.SettingsRoot
}

func NewApp(ctx context.Context) *App {
	ctx, ctxCancel := context.WithCancel(ctx)
	cf := config.Settings
	return &App{
		ctx:       ctx,
		ctxCancel: ctxCancel,
		config:    cf,
	}
}

func (s *App) Run() {
	repository := in_memory.NewInMemoryStorage()
	service := word.NewWordService(repository)

	handlers := handlers.NewHandler(service)
	h := http.NewHTTP(s.ctx, repository, handlers)

	go func() {
		if err := h.Start(); err != nil {
			panic(err)
		}
	}()

	timer, err := time.ParseDuration(s.config.Api.Gci)
	if err != nil {
		fmt.Println(err)
	}
	ticker := time.NewTicker(timer)
	defer ticker.Stop()

	for {
		<-ticker.C
		go repository.CleanGarbageCollector()
	}
}
