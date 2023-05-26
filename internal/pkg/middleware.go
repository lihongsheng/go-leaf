package pkg

import (
  "context"
  "fmt"
  prom "github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
  "github.com/go-kratos/kratos/v2/errors"
  "github.com/go-kratos/kratos/v2/log"
  "github.com/go-kratos/kratos/v2/middleware"
  "github.com/go-kratos/kratos/v2/middleware/logging"
  mmd "github.com/go-kratos/kratos/v2/middleware/metadata"
  "github.com/go-kratos/kratos/v2/middleware/recovery"
  "github.com/go-kratos/kratos/v2/middleware/tracing"
  "github.com/go-kratos/kratos/v2/middleware/validate"
  "github.com/go-kratos/kratos/v2/transport"
  "github.com/go-kratos/kratos/v2/transport/grpc"
  "github.com/go-kratos/kratos/v2/transport/http"
  "github.com/gorilla/handlers"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
  "go-leaf/internal/pkg/metrics"
  "go-leaf/internal/types"
  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/trace"
  g "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"
  "google.golang.org/grpc/keepalive"
  http2 "net/http"
  "strings"
  "time"
)

var (
  _metricSeconds  *prometheus.HistogramVec
  _metricRequests *prometheus.CounterVec
)

// ProviderHTTPPrometheus 确保放在main 文件最前面
func ProviderHTTPPrometheus(ctx context.Context, namespace string, subSystem string) {
  // prometheus 指标只支持下划线格式
  namespace = strings.Replace(namespace, "-", "_", -1)
  subSystem = strings.Replace(subSystem, "-", "_", -1)
  _metricSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
    Namespace: namespace,
    Subsystem: subSystem,
    Name:      "duration_sec",
    Help:      "server requests duration(sec).",
    Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.250, 0.5, 1},
  }, []string{"kind", "operation"})

  _metricRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
    Namespace: namespace,
    Subsystem: subSystem,
    Name:      "code_total",
    Help:      "The total number of processed requests",
  }, []string{"kind", "operation", "code", "reason"})

  prometheus.MustRegister(_metricSeconds, _metricRequests)
}

func providerPrometheus() middleware.Middleware {
  return metrics.Server(
    metrics.WithSeconds(prom.NewHistogram(_metricSeconds)),
    metrics.WithRequests(prom.NewCounter(_metricRequests)),
  )
}

// GrpcClientConn
func GrpcClientConn(target string, logger log.Logger) (*g.ClientConn, error) {
  return grpc.DialInsecure(context.Background(), grpc.WithEndpoint(target),
    grpc.WithOptions(g.WithKeepaliveParams(keepalive.ClientParameters{
      Time:                5 * time.Second, // send pings every 10 seconds if there is no activity
      Timeout:             2 * time.Second, // wait 1 second for ping ack before considering the connection dead
      PermitWithoutStream: true,            // send pings even without active streams
    })),
    grpc.WithMiddleware(
      tracing.Client(),
      logging.Server(logger),
      mmd.Client(),
    ))
}

// NewGrpcClient 原生grpc
func NewGrpcClient(target string, logger log.Logger) (*g.ClientConn, error) {
  var opts []g.DialOption
  opts = append(opts, g.WithTransportCredentials(insecure.NewCredentials()))
  // 加上拦截器
  opts = append(opts, g.WithChainUnaryInterceptor(
    UnaryClientInterceptorTrace(),
    UnaryClientInterceptorLogger(logger)),
    g.WithKeepaliveParams(keepalive.ClientParameters{
      Time:                5 * time.Second, // send pings every 10 seconds if there is no activity
      Timeout:             2 * time.Second, // wait 1 second for ping ack before considering the connection dead
      PermitWithoutStream: true,            // send pings even without active streams
    }))
  return g.Dial(target, opts...)
}

// link https://github.com/grpc-ecosystem/go-grpc-middleware
func UnaryClientInterceptorLogger(logger log.Logger) g.UnaryClientInterceptor {
  return func(ctx context.Context, method string, req, reply interface{}, cc *g.ClientConn, invoker g.UnaryInvoker, opts ...g.CallOption) error {
    err := invoker(ctx, method, req, reply, cc, opts...)
    return err
  }
}

func UnaryClientInterceptorTrace() g.UnaryClientInterceptor {
  return func(parentCtx context.Context, method string, req, reply interface{}, cc *g.ClientConn, invoker g.UnaryInvoker, opts ...g.CallOption) error {
    //span := trace.SpanFromContext(parentCtx)
    tracer := otel.Tracer("grpc-client")
    newCtx, span := tracer.Start(parentCtx, "client-"+method, trace.WithAttributes(), trace.WithSpanKind(trace.SpanKindClient))
    err := invoker(newCtx, method, req, reply, cc, opts...)
    span.End()
    return err
  }
}

func providerLog(logger log.Logger) middleware.Middleware {
  return func(handler middleware.Handler) middleware.Handler {
    return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
      var (
        code      int32
        reason    string
        kind      string
        operation string
      )
      startTime := time.Now()
      if info, ok := transport.FromServerContext(ctx); ok {
        kind = info.Kind().String()
        operation = info.Operation()
      }
      reason = operation
      reply, err = handler(ctx, req)
      if se := errors.FromError(err); se != nil {
        code = se.Code
        reason = se.Reason
      }
      level, stack := extractError(err)
      _ = log.WithContext(ctx, logger).Log(level,
        log.DefaultMessageKey, reason,
        "kind", "server",
        "component", kind,
        "operation", operation,
        "args", extractArgs(req),
        "code", code,
        "reason", reason,
        "stack", stack,
        "latency", time.Since(startTime).Seconds(),
      )
      return
    }
  }
}

func extractError(err error) (log.Level, string) {
  if err != nil {
    return log.LevelError, fmt.Sprintf("%+v", err)
  }
  return log.LevelInfo, ""
}

// extractArgs returns the string of the req
func extractArgs(req interface{}) string {
  if redacter, ok := req.(logging.Redacter); ok {
    return redacter.Redact()
  }
  if stringer, ok := req.(fmt.Stringer); ok {
    return stringer.String()
  }
  return fmt.Sprintf("%+v", req)
}

func NewGrpcServer(log log.Logger, addr string, t time.Duration) *grpc.Server {
  var opts = []grpc.ServerOption{
    grpc.Middleware(
      recovery.Recovery(),
      tracing.Server(),
      providerLog(log),
      validate.Validator(),
      mmd.Server(),
      providerPrometheus(),
    ),
    grpc.Address(addr),
  }
  if t > 0 {
    opts = append(opts, grpc.Timeout(t))
  }
  return grpc.NewServer(opts...)
}

func NewHTTPServer(log log.Logger, addr string, t time.Duration) *http.Server {
  var opts = []http.ServerOption{
    http.Middleware(
      recovery.Recovery(),
      tracing.Server(),
      providerLog(log),
      validate.Validator(),
      mmd.Server(),
      providerPrometheus(),
    ),
    // 允许跨域请求
    // 如果服务需要告知客户端支持跨域，还需要在自行加入拦截器，设置header允许跨域
    http.Filter(handlers.CORS(
      handlers.AllowedOrigins([]string{"*"}),
      handlers.AllowedMethods([]string{"GET", "POST", "PUT"}),
    )),
    http.Address(addr),
    // 这里不要随意改
    http.ResponseEncoder(func(w http2.ResponseWriter, r *http2.Request, v interface{}) error {
      // DefaultResponseEncoder encodes the object to the HTTP response.
      // 创蓝回调，需要写入ok
      if r.URL.Path == types.ChuangLanCallBack {
        _, _ = w.Write([]byte("ok"))
        return nil
      }
      return http.DefaultResponseEncoder(w, r, v)
    }),
  }
  if t > 0 {
    opts = append(opts, http.Timeout(t))
  }

  srv := http.NewServer(opts...)
  srv.Handle(types.PromHTTPHandlerPath, promhttp.Handler())
  srv.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
    _, _ = w.Write([]byte("ok"))
  })
  return srv
}
