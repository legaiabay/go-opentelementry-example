package config

import (
	"context"
	"log"
	"runtime"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const environment = "development"
const service = "payroll-service"
const url = "http://localhost:14268/api/traces"

var Tracer OpenTelemetry

type OpenTelemetry struct {
	Context        context.Context
	TracerProvider *tracesdk.TracerProvider
	Tracer         oteltrace.Tracer
	MainSpan       oteltrace.Span
	Cancel         context.CancelFunc
}

func Init(name string) OpenTelemetry {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		log.Fatal(err)
	}

	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),

		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", 1),
		)),
	)

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)

	// Create initial context
	ctx, cancel := context.WithCancel(context.Background())

	// Get current running func	tion name
	fpcs := make([]uintptr, 1)
	_ = runtime.Callers(2, fpcs)
	f := runtime.FuncForPC(fpcs[0] - 1)

	tr := otel.Tracer(f.Name())

	ctx, span := tr.Start(ctx, name)

	return OpenTelemetry{
		Context:        ctx,
		TracerProvider: tp,
		Tracer:         tr,
		MainSpan:       span,
		Cancel:         cancel,
	}

}

func (s *OpenTelemetry) Trace(name string, params string) oteltrace.Span {
	// Get current running function name
	fpcs := make([]uintptr, 1)
	_ = runtime.Callers(2, fpcs)
	f := runtime.FuncForPC(fpcs[0] - 1)

	tr := otel.Tracer(f.Name())

	_, span := tr.Start(s.Context, name)
	span.SetAttributes(attribute.Key("value").String(params))

	return span
}

func (s *OpenTelemetry) End() {
	s.MainSpan.End()

	s.Context, s.Cancel = context.WithTimeout(s.Context, time.Second*5)
	defer s.Cancel()
	if err := s.TracerProvider.Shutdown(s.Context); err != nil {
		log.Fatal(err)
	}
}
