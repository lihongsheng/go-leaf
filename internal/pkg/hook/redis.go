package hook

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net"
)

type RedisHook struct {
}

func NewRedisHook() redis.Hook {
	return &RedisHook{}
}

func (h *RedisHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

func (h *RedisHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		var span trace.Span
		if ctx != nil {
			if pSpan := trace.SpanContextFromContext(ctx); pSpan.HasTraceID() {
				tracer := otel.Tracer("redis")
				_, span = tracer.Start(ctx, "redis-cmd", trace.WithAttributes(), trace.WithSpanKind(trace.SpanKindServer))
			}
		}
		err := next(ctx, cmd)
		if span != nil {
			span.SetAttributes(attribute.String("query", cmd.String()))
			span.End()
		}
		return err
	}
}

func (h *RedisHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		var span trace.Span
		if ctx != nil {
			if pSpan := trace.SpanContextFromContext(ctx); pSpan.HasTraceID() {
				tracer := otel.Tracer("redis")
				_, span = tracer.Start(ctx, "redis-pipe", trace.WithAttributes(), trace.WithSpanKind(trace.SpanKindServer))
			}
		}
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		return next(ctx, cmds)
	}
}
