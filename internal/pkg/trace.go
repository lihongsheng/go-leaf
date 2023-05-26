package pkg

import (
	"context"
	"github.com/Shopify/sarama"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	AllTrace = "AllTrace"
)

type spanConfig struct {
	TraceID string `json:"trace_id"`
	SpanID  string `json:"span_id"`
}

type AMQPHeadersCarrier map[string]interface{}

func (a AMQPHeadersCarrier) Get(key string) string {
	v, ok := a[key]
	if !ok {
		return ""
	}
	return v.(string)
}

func (a AMQPHeadersCarrier) Set(key string, value string) {
	a[key] = value
}

func (a AMQPHeadersCarrier) Keys() []string {
	i := 0
	r := make([]string, len(a))

	for k := range a {
		r[i] = k
		i++
	}

	return r
}

// InjectAMQPHeaders injects the tracing from the context into the header map
func InjectAMQPHeaders(ctx context.Context) map[string]interface{} {
	h := make(AMQPHeadersCarrier)
	//otel.GetTextMapPropagator().Inject(ctx, h)
	if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
		textMap := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
		textMapCarr := propagation.MapCarrier{}
		textMap.Inject(ctx, textMapCarr)
		for _, k := range textMapCarr.Keys() {
			h.Set(k, textMapCarr.Get(k))
		}
	}
	return h
}

// ExtractAMQPHeaders extracts the tracing from the header and puts it into the context
func ExtractAMQPHeaders(ctx context.Context, headers map[string]interface{}) context.Context {
	if len(headers) > 0 {
		textMapCarr := propagation.MapCarrier{}
		for k, v := range headers {
			if str, ok := v.(string); ok {
				textMapCarr.Set(k, str)
			}
		}
		textMap := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
		ctx = textMap.Extract(ctx, textMapCarr)
	}
	return ctx
}

func InjectKafkaHeaders(ctx context.Context) []sarama.RecordHeader {
	var headers = make([]sarama.RecordHeader, 0, 1)
	if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
		textMap := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
		textMapCarr := propagation.MapCarrier{}
		textMap.Inject(ctx, textMapCarr)
		for _, k := range textMapCarr.Keys() {
			headers = append(headers, sarama.RecordHeader{
				Key:   []byte(k),
				Value: []byte(textMapCarr.Get(k)),
			})
		}
	}
	return headers
}

// ExtractKafkaHeaders extracts the tracing from the header and puts it into the context
func ExtractKafkaHeaders(ctx context.Context, headers []*sarama.RecordHeader) context.Context {
	if len(headers) > 0 {
		textMap := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
		textMapCarr := propagation.MapCarrier{}
		for _, header := range headers {
			textMapCarr.Set(string(header.Key), string(header.Value))
		}
		ctx = textMap.Extract(ctx, textMapCarr)
	}

	return ctx
}
