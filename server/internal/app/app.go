package app

import (
	"github.com/kjuulh/scel/server/internal/persistence"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type App struct {
	Tracer         trace.Tracer
	DownloadsStore persistence.DownloadsStore
}

func NewApp() *App {
	var (
		tracer         = otel.Tracer("scel_server")
		downloadsStore = persistence.NewInMemoryDownloadsStore()
	)

	return &App{
		Tracer:         tracer,
		DownloadsStore: downloadsStore,
	}
}
