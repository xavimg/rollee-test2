package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"words/internal/handlers"
	"words/internal/service"
	"words/internal/storage"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

type ServerHTTP struct {
	addr   string
	router *chi.Mux
	server *http.Server

	handlers *handlers.WordHandler
	service  service.WordServicer
	storage  storage.WordStorager

	stopChan chan os.Signal
	errChan  chan error
}

func NewHTTP(addr string, gci time.Duration, storage storage.WordStorager, service service.WordServicer, handler *handlers.WordHandler) *ServerHTTP {
	router := chi.NewRouter()
	return &ServerHTTP{
		addr:     addr,
		router:   router,
		handlers: handler,
		service:  service,
		storage:  storage,
	}
}

// Start initializes and begins running the server, listening for incoming requests.
// It gracefully handles system interrupts or errors to shut down the server.
// The function returns nil if the server shuts down gracefully, or an error if there's an issue during its operation.
func (s *ServerHTTP) Start() error {
	s.routes()

	s.stopChan = make(chan os.Signal, 1)
	s.errChan = make(chan error, 1)
	signal.Notify(s.stopChan, os.Interrupt, syscall.SIGTERM)

	s.server = &http.Server{
		Addr:    s.addr,
		Handler: s.router,
	}
	go func() {
		log.Info().Msgf("starting server on port %s", s.addr) // Log right after successful binding
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.errChan <- err
			return
		}
	}()

	select {
	case <-s.stopChan:
		signal.Stop(s.stopChan)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.server.Shutdown(ctx); err != nil {
			log.Error().Msgf("server shutdown error: %v", err)
		}
		log.Info().Msg("server gracefully closed")
		return nil
	case err := <-s.errChan:
		return err
	}
}

func (s *ServerHTTP) routes() {
	s.router.Group(func(r chi.Router) {
		r.Post("/api/v0.1/words/{word}", s.handlers.AddWord)
		r.Get("/api/v0.1/words/{prefix}", s.handlers.FrequentWordByPrefix)
	})
}
