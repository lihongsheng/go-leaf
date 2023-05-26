package pkg

import (
	"context"
	"errors"
	"github.com/Shopify/sarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"message-center/internal/tools"
	"testing"
)

func TestInjectAMQPHeaders2(t *testing.T) {
	// jaeger初始化，用于链路追踪，及后续打日志里的traceID记录
	//_ = tools.InitJaeger("http://127.0.0.1:14268/api/traces", "test3")
	//start := otel.Tracer("test")
	//ctx := context.Background()
	//spanCtx, span := start.Start(ctx, "test", trace.WithAttributes(), trace.WithSpanKind(trace.SpanKindProducer))
	//if span := trace.SpanContextFromContext(spanCtx); !span.HasTraceID() {
	//	t.Errorf("not trace %s", span.TraceID().String())
	//}
	//time.Sleep(1 * time.Second)
	//span.End()
	//textMap := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	//textMap2 := propagation.MapCarrier{}
	//textMap.Inject(spanCtx, textMap2)
	//fmt.Println(textMap2)
	//time.Sleep(10 * time.Second)
}
func TestInjectAMQPHeaders3(t *testing.T) {
	//_ = tools.InitJaeger("http://127.0.0.1:14268/api/traces", "test2")
	//textMap2 := propagation.MapCarrier{}
	//textMap2.Set("traceparent", "00-11a8a5337f83dd55e6dcc85ac01060af-65d60eb96ece96b8-01")
	//ctx2 := context.Background()
	//textMap3 := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	//ctx2 = textMap3.Extract(ctx2, textMap2)
	//if span := trace.SpanContextFromContext(ctx2); span.HasTraceID() {
	//	t.Logf("not trace %s", span.TraceID().String())
	//}
	//
	//start2 := otel.Tracer("test2")
	//_, span2 := start2.Start(ctx2, "test2", trace.WithAttributes(), trace.WithSpanKind(trace.SpanKindConsumer))
	//time.Sleep(1 * time.Second)
	//
	//span2.End()
	//
	//time.Sleep(10 * time.Second)
}

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
