package http

import (
	"context"
	"net/http"
	"os"
	"time"

	"words/internal/config"
	"words/internal/http/handlers"
	"words/internal/storage"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

type ServerHTTP struct {
	ctx       context.Context
	ctxCancel context.CancelFunc

	addr   string
	router *chi.Mux
	server *http.Server

	handlers   *handlers.WordHandler
	repository storage.WordRepository

	stopChan chan os.Signal
	errChan  chan error
}

// NewHTTP accepting interface for storage, it offers the flexibility
// to easily switch between different implementations. For instance, we can effortlessly swap our
// InMemoryStore with a PostgreSQLStore or even replace the WordService with another service implementation,
// ensuring a high level of decoupling and adaptability.
func NewHTTP(ctx context.Context, gci time.Duration, repository storage.WordRepository, handler *handlers.WordHandler) *ServerHTTP {
	addr := config.Settings.Api.Port
	ctx, ctxCancel := context.WithCancel(ctx)
	router := chi.NewRouter()
	return &ServerHTTP{
		ctx:        ctx,
		ctxCancel:  ctxCancel,
		addr:       addr,
		router:     router,
		handlers:   handler,
		repository: repository,
	}
}

// Start initializes and begins running the server, listening for incoming requests.
// It gracefully handles system interrupts or errors to shut down the server.
// The function returns nil if the server shuts down gracefully, or an error if there's an issue during its operation.
func (s *ServerHTTP) Start() error {
	s.routes()

	s.stopChan = make(chan os.Signal, 1)
	s.errChan = make(chan error, 1)

	s.server = &http.Server{
		Addr:    s.addr,
		Handler: s.router,
	}
	go func() {
		log.Info().Msgf("starting server on port %s", s.addr) // Log right after successful binding
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("Failed to start the server: %v", err)
			return
		}
	}()

	select {
	case <-s.ctx.Done():
		if err := s.server.Shutdown(s.ctx); err != nil {
			log.Error().Msgf("server shutdown error: %v", err)
		}
		log.Info().Msg("server gracefully closed")
		return nil
	}
}

func (s *ServerHTTP) routes() {
	s.router.Group(func(r chi.Router) {
		r.Post("/api/v0.1/words/{word}", s.handlers.AddWord)
		r.Get("/api/v0.1/words/{prefix}", s.handlers.FrequentWordByPrefix)
	})
}
