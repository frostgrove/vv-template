package observability

import (
	"context"
	"errors"
	"log/slog"
	"time"

	vvotel "github.com/frostgrove/vv/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/fx"

	"github.com/frostgrove/vv-template/src/config"
)

// Telemetry is the application-owned OpenTelemetry state. The framework only
// borrows the providers (vvotel); the SDK, the exporters, sampling, flush and
// shutdown belong here. Nothing is installed as an OTel global.
type Telemetry struct {
	Enabled        bool
	TracerProvider trace.TracerProvider
	MeterProvider  metric.MeterProvider
	Propagator     propagation.TextMapPropagator
	Frost          *vvotel.Telemetry
}

func NewTelemetry(lifecycle fx.Lifecycle, configuration *config.Config, logger *slog.Logger) (*Telemetry, error) {
	settings := configuration.Observability.Otel
	if !settings.Enabled {
		return &Telemetry{
			TracerProvider: tracenoop.NewTracerProvider(),
			MeterProvider:  metricnoop.NewMeterProvider(),
			Propagator:     propagation.NewCompositeTextMapPropagator(),
			Frost:          vvotel.Must(vvotel.Config{Disabled: true}),
		}, nil
	}

	ctx := context.Background()
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.ServiceName(configuration.App.Name),
			semconv.DeploymentEnvironment(configuration.App.Env),
		),
	)
	if err != nil {
		return nil, err
	}

	var traceOptions []otlptracegrpc.Option
	var metricOptions []otlpmetricgrpc.Option
	if settings.Endpoint != "" {
		traceOptions = append(traceOptions, otlptracegrpc.WithEndpoint(settings.Endpoint))
		metricOptions = append(metricOptions, otlpmetricgrpc.WithEndpoint(settings.Endpoint))
	}
	if settings.Insecure {
		traceOptions = append(traceOptions, otlptracegrpc.WithInsecure())
		metricOptions = append(metricOptions, otlpmetricgrpc.WithInsecure())
	}
	traceExporter, err := otlptracegrpc.New(ctx, traceOptions...)
	if err != nil {
		return nil, err
	}
	metricExporter, err := otlpmetricgrpc.New(ctx, metricOptions...)
	if err != nil {
		_ = traceExporter.Shutdown(ctx)
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(settings.SampleRatio))),
		sdktrace.WithBatcher(traceExporter),
	)
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
	)
	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			flushCtx, cancelFlush := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
			defer cancelFlush()
			flushed := errors.Join(tracerProvider.ForceFlush(flushCtx), meterProvider.ForceFlush(flushCtx))

			shutdownCtx, cancelShutdown := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
			defer cancelShutdown()
			return errors.Join(flushed, tracerProvider.Shutdown(shutdownCtx), meterProvider.Shutdown(shutdownCtx))
		},
	})

	frost, err := vvotel.New(vvotel.Config{TracerProvider: tracerProvider, MeterProvider: meterProvider})
	if err != nil {
		return nil, err
	}
	logger.Info("opentelemetry export is on",
		slog.String("endpoint", settings.Endpoint), slog.Float64("sample_ratio", settings.SampleRatio))
	return &Telemetry{
		Enabled:        true,
		TracerProvider: tracerProvider,
		MeterProvider:  meterProvider,
		Propagator:     propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
		Frost:          frost,
	}, nil
}
