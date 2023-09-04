package http

import (
	"context"
	"net/http"

	"words/internal/config"
	"words/internal/http/handlers"
	"words/internal/storage"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

const (
	logGracefully     = "server gracefully closed"
	logServerFailed   = "Failed to start the server: %v"
	logServerShutdown = "server shutdown error: %v"
	logServerStarted  = "starting server on port %s"
)

type ServerHTTP struct {
	ctx       context.Context
	ctxCancel context.CancelFunc

	addr   string
	router *chi.Mux
	server *http.Server

	handlers   *handlers.WordHandler
	repository storage.WordRepository
}

// NewHTTP accepting interface for storage, it offers the flexibility
// to easily switch between different implementations. For instance, we can effortlessly swap our
// InMemoryStore with a PostgreSQLStore or even replace the WordService with another service implementation,
// ensuring a high level of decoupling and adaptability.
func NewHTTP(ctx context.Context, repository storage.WordRepository, handler *handlers.WordHandler) *ServerHTTP {
	ctx, ctxCancel := context.WithCancel(ctx)

	return &ServerHTTP{
		ctx:        ctx,
		ctxCancel:  ctxCancel,
		addr:       config.Settings.Api.Port,
		router:     chi.NewRouter(),
		handlers:   handler,
		repository: repository,
	}
}

// Start initializes and begins running the server, listening for incoming requests.
// It gracefully handles system interrupts or errors to shut down the server.
// The function returns nil if the server shuts down gracefully, or an error if there's an issue during its operation.
func (s *ServerHTTP) Start() error {
	s.routes()

	s.server = &http.Server{
		Addr:    s.addr,
		Handler: s.router,
	}

	go func() {
		log.Info().Msgf(logServerStarted, s.addr)

		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf(logServerFailed, err)
			return
		}
	}()

	<-s.ctx.Done()
	if err := s.server.Shutdown(s.ctx); err != nil {
		log.Error().Msgf(logServerShutdown, err)
	}
	log.Info().Msg(logGracefully)

	return nil
}

func (s *ServerHTTP) routes() {
	s.router.Group(func(r chi.Router) {
		r.Post("/api/v0.1/words/{word}", s.handlers.AddWord)
		r.Get("/api/v0.1/words/{prefix}", s.handlers.FrequentWordByPrefix)
	})
}
