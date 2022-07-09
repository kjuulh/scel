package app

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("scel_server")

type App struct {
	Tracer trace.Tracer
}

func NewApp() *App {
	return &App{
		Tracer: tracer,
	}
}
