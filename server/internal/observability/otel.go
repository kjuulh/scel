package observability

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"google.golang.org/grpc"
)

func getOtlpEndpoint(ctx context.Context) (string, error) {
	otelAgentAddr, ok := os.LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if !ok {
		otelAgentAddr = "0.0.0.0:4317"
	}

	return otelAgentAddr, nil
}

func NewOtlpMetricsClient(ctx context.Context) (*controller.Controller, error) {
	otlpEndpoint, err := getOtlpEndpoint(ctx)
	if err != nil {
		return nil, err
	}

	metricsClient := otlpmetricgrpc.NewClient(
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(otlpEndpoint),
	)

	metricsExp, err := otlpmetric.New(ctx, metricsClient)

	pusher := controller.New(
		processor.NewFactory(simple.NewWithHistogramDistribution(), metricsExp),
		controller.WithExporter(metricsExp),
		controller.WithCollectPeriod(2*time.Second),
	)
	global.SetMeterProvider(pusher)

	err = pusher.Start(ctx)
	if err != nil {
		return nil, err
	}

	return pusher, nil
}

func NewOtlpTraceClient(ctx context.Context) (*otlptrace.Exporter, error) {
	otlpEndpoint, err := getOtlpEndpoint(ctx)
	if err != nil {
		return nil, err
	}

	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otlpEndpoint),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)

	traceExp, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(
		ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("scel-server"),
		),
	)
	if err != nil {
		return nil, err
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExp)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(
			sdktrace.AlwaysSample(),
		),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	)
	otel.SetTracerProvider(tracerProvider)

	return traceExp, nil
}

func NewOtlp(ctx context.Context) (func(), error) {
	err := NewOtelLogger()
	if err != nil {
		return nil, err
	}

	metricsClient, err := NewOtlpMetricsClient(ctx)
	if err != nil {
		return nil, err
	}
	tracerClient, err := NewOtlpTraceClient(ctx)
	if err != nil {
		return nil, err
	}

	return func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()

		if err := tracerClient.Shutdown(ctx); err != nil {
			otel.Handle(err)
		}

		if err := metricsClient.Stop(ctx); err != nil {
			otel.Handle(err)
		}
	}, nil
}
