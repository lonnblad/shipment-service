package trace

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/metric/global"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"

	"github.com/lonnblad/shipment-service-backend/config"
)

var pipeline *controller.Controller

// Start will setup an Open Telemetry tracer.
func Start() error {
	resource := resource.NewWithAttributes(
		semconv.ServiceNameKey.String(config.GetServiceName()),
		semconv.ServiceVersionKey.String(config.GetServiceVersion()),
		semconv.DeploymentEnvironmentKey.String(config.GetEnvironment().String()),
	)

	exporter, err := stdout.NewExporter(stdout.WithPrettyPrint())
	if err != nil {
		return err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(resource),
		sdktrace.WithBatcher(exporter),
	)

	pipeline = controller.New(
		processor.New(
			simple.NewWithInexpensiveDistribution(),
			exporter,
		),
		append(
			[]controller.Option{controller.WithResource(resource)},
			controller.WithExporter(exporter),
		)...,
	)

	if err = pipeline.Start(context.Background()); err != nil {
		return err
	}

	otel.SetTracerProvider(tp)
	global.SetMeterProvider(pipeline.MeterProvider())

	return nil
}

// Stop will stop the existing tracer.
func Stop(ctx context.Context) {
	if err := pipeline.Stop(ctx); err != nil {
		log.Println(err)
		return
	}
}

// Tracer will return a default tracer.
func Tracer() trace.Tracer {
	return otel.Tracer("")
}
