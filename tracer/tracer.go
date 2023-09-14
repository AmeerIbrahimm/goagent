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

var (
	exporter *otlptrace.Exporter
	bsp      *sdktrace.SpanProcessor
	tracr    *sdktrace.TracerProvider
)

func NewTraceProvider(res *resource.Resource, bsp sdktrace.SpanProcessor) *sdktrace.TracerProvider {
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	setTracerProvider(tracerProvider)
	return tracerProvider
}

func NewExporter(ctx context.Context, conn *grpc.ClientConn) (*otlptrace.Exporter, error) {
	exp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	setExporter(exp)

	return exp, err
}
func NewResource(ctx context.Context) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			attribute.String("service.name", config.ServiceName),
			attribute.String("service.identifier", config.ServiceIdentifier),
			attribute.String("app.name", config.AppName),
			attribute.String("telemetry.sdk.language", "go"),
		),
	)
}

func setExporter(exportr *otlptrace.Exporter) {
	exportr = exportr
}

func GetExporter() *otlptrace.Exporter {
	return exporter
}

func setSpanProcessor(processor *sdktrace.SpanProcessor) {
	bsp = processor
}

func GetProcessor() *sdktrace.SpanProcessor {
	return bsp
}

func setTracerProvider(tcr *sdktrace.TracerProvider) {
	tracr = tcr
}

func GetTracerProvider() *sdktrace.TracerProvider {
	return tracr
}
