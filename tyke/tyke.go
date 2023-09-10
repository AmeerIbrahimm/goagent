package tyke

import (
	"context"
	"log"
	"time"

	"github.com/AmeerIbrahimm/goagent/tracer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitTracer() func() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	res, err := tracer.NewResource(ctx)
	reportErr(err, "failed to create res")

	conn, err := grpc.DialContext(ctx, "localhost:4317", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	reportErr(err, "failed to create gRPC connection to collector")

	// Set up a trace exporter
	traceExporter, err := tracer.NewExporter(ctx, conn)
	reportErr(err, "failed to create trace exporter")

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	batchSpanProcessor := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := tracer.NewTraceProvider(res, batchSpanProcessor)
	otel.SetTracerProvider(tracerProvider)

	return func() {
		// Shutdown will flush any remaining spans and shut down the exporter.
		reportErr(tracerProvider.Shutdown(ctx), "failed to shutdown TracerProvider")
		cancel()
	}
}

func reportErr(err error, message string) {
	if err != nil {
		log.Printf("%s: %v", message, err)
	}
}

func GetTraceProvider() trace.TracerProvider {
	return otel.GetTracerProvider()
}

func GetMeterProvider() metric.MeterProvider {
	return otel.GetMeterProvider()
}

func GetTextMapPropagator() propagation.TextMapPropagator {
	return otel.GetTextMapPropagator()
}
