package types

const (
	PromHTTPHandlerPath   = "/metrics"
	MetaAuthorization     = "authorization"
	HTPPHeadAuthorization = "Authorization"
)

const (
	LogServerID      = "service_id"
	LogServerName    = "service_name"
	LogServerVersion = "service_version"
	LogTraceID       = "trace_id"
	LogSpanID        = "span_id"
	LogEnv           = "env"
	LogPodName       = "pod_name"
	LogGrpcMethod    = "grpc_method"
	LogRequestURL    = "request_url"
	LogStackTrace    = "stacktrace"
)

var SkipLog = map[string]struct{}{
	LogServerName: {},
	LogTraceID:    {},
	LogSpanID:     {},
	LogPodName:    {},
	LogEnv:        {},
	LogStackTrace: {},
}

type MetaDataAppIDCtx struct{}
