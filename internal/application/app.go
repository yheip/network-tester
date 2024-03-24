package application

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/yheip/network-tester/internal/http"
	"golang.org/x/sync/errgroup"
)

type HttpServer interface {
	Start() error
	Shutdown() error
}
type Application struct {
	server HttpServer
}

func New(ctx context.Context) *Application {
	app := &Application{
		server: http.NewServer(),
	}

	return app
}

func (a *Application) Start(ctx context.Context) error {
	log := newLogger()

	ctx = log.WithContext(ctx)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return a.server.Start()
	})

	g.Go(func() error {
		<-ctx.Done()

		log.Info().Msg("Gracefully shutting down...")

		return a.server.Shutdown()
	})

	return g.Wait()
}

func newLogger() *zerolog.Logger {
	var output io.Writer
	if os.Getenv("ENV") == "" {
		output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	} else {
		output = os.Stdout
	}
	log := zerolog.New(output).With().Timestamp().Logger()

	zerolog.DefaultContextLogger = &log

	return &log
}
