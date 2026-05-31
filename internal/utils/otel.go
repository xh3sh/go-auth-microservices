package utils

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func InitTracer() *sdktrace.TracerProvider {
	tp := sdktrace.NewTracerProvider()

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp
}

func GetTraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	traceID := span.SpanContext().TraceID().String()
	if traceID == "00000000000000000000000000000000" {
		return ""
	}
	return traceID
}
