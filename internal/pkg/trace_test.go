package pkg

import (
	"context"
	"errors"
	"github.com/Shopify/sarama"
	"go-leaf/internal/tools"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"testing"
)

func TestExtractAMQPHeaders(t *testing.T) {
	// jaeger初始化，用于链路追踪，及后续打日志里的traceID记录
	_ = tools.InitJaeger("http://127.0.0.1:14268/api/traces", "test")
	start := otel.Tracer("test")
	ctx := context.Background()
	spanCtx, span := start.Start(ctx, "test", trace.WithAttributes(), trace.WithSpanKind(trace.SpanKindConsumer))
	if span := trace.SpanContextFromContext(spanCtx); !span.HasTraceID() {
		t.Errorf("not trace %s", span.TraceID().String())
	}
	defer span.End()
	header := InjectAMQPHeaders(spanCtx)
	t.Log(header)
	newCtx := ExtractAMQPHeaders(ctx, header)
	if span := trace.SpanContextFromContext(newCtx); !span.HasTraceID() {
		t.Errorf("heder not trace %s", span.TraceID().String())
	}
}

func TestInjectAMQPHeaders(t *testing.T) {
	// jaeger初始化，用于链路追踪，及后续打日志里的traceID记录
	_ = tools.InitJaeger("http://127.0.0.1:14268/api/traces", "test")
	start := otel.Tracer("test")
	ctx := context.Background()
	spanCtx, span := start.Start(ctx, "test", trace.WithAttributes(), trace.WithSpanKind(trace.SpanKindConsumer))
	if span := trace.SpanContextFromContext(spanCtx); !span.HasTraceID() {
		t.Errorf("not trace %s", span.TraceID().String())
	}
	defer span.End()
	if span := trace.SpanContextFromContext(spanCtx); span.HasTraceID() {
		t.Logf(" trace %s", span.TraceID().String())
	}
	headers := InjectAMQPHeaders(spanCtx)
	if len(headers) <= 0 {
		t.Error(errors.New("fail ,not find headers"))
	}
	t.Log(headers)
}

func TestInjectKafkaHeaders(t *testing.T) {
	// jaeger初始化，用于链路追踪，及后续打日志里的traceID记录
	_ = tools.InitJaeger("http://127.0.0.1:14268/api/traces", "test")
	start := otel.Tracer("test")
	ctx := context.Background()
	spanCtx, span := start.Start(ctx, "test", trace.WithAttributes(), trace.WithSpanKind(trace.SpanKindConsumer))
	if span := trace.SpanContextFromContext(spanCtx); !span.HasTraceID() {
		t.Errorf("not trace %s", span.TraceID().String())
	}
	defer span.End()

	headers := InjectKafkaHeaders(spanCtx)
	if len(headers) <= 0 {
		t.Error(errors.New("fail ,not find headers"))
	}
	t.Log(headers)
}

func TestExtractKafkaHeaders(t *testing.T) {
	// jaeger初始化，用于链路追踪，及后续打日志里的traceID记录
	_ = tools.InitJaeger("http://127.0.0.1:14268/api/traces", "test")
	start := otel.Tracer("test")
	ctx := context.Background()
	spanCtx, span := start.Start(ctx, "test", trace.WithAttributes(), trace.WithSpanKind(trace.SpanKindConsumer))
	if span := trace.SpanContextFromContext(spanCtx); !span.HasTraceID() {
		t.Errorf("not trace %s", span.TraceID().String())
	}
	defer span.End()
	headers := InjectKafkaHeaders(spanCtx)
	headers2 := make([]*sarama.RecordHeader, 0, len(headers))
	for _, header := range headers {
		headers2 = append(headers2, &header)
	}
	newCtx := ExtractKafkaHeaders(context.Background(), headers2)
	t.Log(headers)
	newSpanInfo := trace.SpanContextFromContext(newCtx)
	oldSpanInfo := trace.SpanContextFromContext(spanCtx)
	if newSpanInfo.TraceID().String() == "" || oldSpanInfo.TraceID().String() != newSpanInfo.TraceID().String() {
		t.Error(errors.New("fail ExtractKafkaHeaders"))
	}
}
