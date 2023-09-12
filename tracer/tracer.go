package tracer

import (
	"context"

	"github.com/AmeerIbrahimm/goagent/config"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	// semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	"google.golang.org/grpc"
)

func NewTraceProvider(res *resource.Resource, bsp sdktrace.SpanProcessor) *sdktrace.TracerProvider {
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	return tracerProvider
}

func NewExporter(ctx context.Context, conn *grpc.ClientConn) (*otlptrace.Exporter, error) {
	return otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
}
func NewResource(ctx context.Context) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			attribute.String("service.name", config.ServiceName),
			attribute.String("library.language", "go"),
		),
	)
}
