package observability

import (
	"context"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

func NewOtelLogger() error {
	l, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	log := otelzap.New(l, otelzap.WithMinLevel(zap.InfoLevel))

	_ = otelzap.ReplaceGlobals(log)

	return nil
}

func Logger(ctx context.Context) otelzap.LoggerWithCtx {
	return otelzap.Ctx(ctx)
}
