package tyke

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/AmeerIbrahimm/goagent/config"
	"github.com/AmeerIbrahimm/goagent/tracer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func InitTracer() func() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	var secureOption grpc.DialOption

	if strings.ToLower(config.Insecure) == "false" || config.Insecure == "0" || strings.ToLower(config.Insecure) == "f" {
		secureOption = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))
	} else {
		secureOption = grpc.WithTransportCredentials(insecure.NewCredentials())
	}
	conn, err := grpc.DialContext(ctx, config.CollectorEndpoint, secureOption, grpc.WithBlock())
	reportErr(err, "failed to create gRPC connection to collector")

	exporter, err := tracer.NewExporter(ctx, conn)
	if err != nil {
		log.Fatalf("Failed to create exporter: %v", err)
	}
	resources, err := tracer.NewResource(ctx)
	if err != nil {
		log.Fatalf("Could not set resources: %v", err)
	}
	batchSpanProcessor := sdktrace.NewBatchSpanProcessor(exporter)

	otel.SetTracerProvider(
		// sdktrace.NewTracerProvider(
		// 	sdktrace.WithSampler(sdktrace.AlwaysSample()),
		// 	sdktrace.WithBatcher(exporter),
		// 	sdktrace.WithResource(resources),
		// ),
		tracer.NewTraceProvider(resources, batchSpanProcessor),
	)
	return func() {
		// Shutdown will flush any remaining spans and shut down the exporter.
		reportErr(exporter.Shutdown(ctx), "failed to shutdown TracerProvider")
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
