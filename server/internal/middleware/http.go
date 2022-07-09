package middleware

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/context"
)

func ObservabilityMiddleware(ctx context.Context, next http.Handler, method string) http.Handler {
	meter := global.Meter("scel-meter")
	serverAttribute := attribute.String("scel-attribute", "scel-server")
	commonLabels := []attribute.KeyValue{serverAttribute}
	requestCount, _ := meter.
		SyncInt64().
		Counter("scel_server/request_counts", instrument.WithDescription("The number of requests received"))

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestCount.Add(ctx, 1, commonLabels...)
		span := trace.SpanFromContext(ctx)
		bag := baggage.FromContext(ctx)

		baggageAttributes := []attribute.KeyValue{}
		baggageAttributes = append(baggageAttributes, serverAttribute)
		for _, member := range bag.Members() {
			baggageAttributes = append(baggageAttributes, attribute.String("baggage key:"+member.Key(), member.Value()))
		}
		span.SetAttributes(baggageAttributes...)

		next.ServeHTTP(w, r)
	})

	return otelhttp.NewHandler(h, method)
}
