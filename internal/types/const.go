package types

import "time"

const (
	NameSpace                       = "message_center"
	ConsumerChannelLabelName        = "channel"
	ConsumerChannelSuccessLabelName = "code"
	ConsumerSuccessName             = "success"
	ConsumerFailName                = "fail"
	ConsumerTemplateLabelName       = "template"
	PromHTTPHandlerPath             = "/message_center_679b94a_metrics"
	MetaAppID                       = "x-md-global-appid"
	MetaAppSign                     = "x-md-global-sign"
	MetaAppTime                     = "x-md-global-app-time"

	// 特殊的
	ChuangLanCallBack = "/callback/chuanglan"
	EmailReplyTo      = "reply_to"
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

const (
	TemplateCacheSimpleKey  = "template_cache_simple_key_%d"
	TemplateCacheSimpleTime = 24 * time.Hour
	TemplateCacheAppIDKey   = "template_cache_appID_key_%s_%d"
	TemplateCacheAPPIDTime  = 24 * time.Hour

	TemplateCacheIDKey  = "template_cache_key_%d"
	TemplateCacheIDTime = 24 * time.Hour
	// ChannelCacheKey %d:channelID, %s channelType
	ChannelCacheKey  = "Channel_Cache_Key_%d"
	ChannelCacheTime = 24 * time.Hour
	// 基于appID 缓存
	ChannelAppIDCacheKey  = "Channel_appID_Cache_Key_%s_%s"
	ChannelAppIDCacheTime = 24 * time.Hour
	// app table %s : app_id
	AppTableCacheKey  = "app_cache_key_%s"
	AppTableCacheTime = 24 * time.Hour
)

const (
	// SendKey , explame: Email Send %s:bacthID, %s:email
	SendKey       = "send_key_%s_%s"
	SendTime      = 2 * time.Hour
	SendLockKey   = "send_lock_key_%s_%s"
	SendLockTime  = 20 * time.Minute
	SendRetryTime = 20 * time.Minute
)

const (
	PageDefault     = 1
	PageSizeDefault = 20
	// ES
	ActionUpdate = "update"
	ActionCreate = "create"
)

const (
	EmailAwsRateLimit = 10 // aws邮件发送速率限制，每秒发20个
)

type MetaDataAppIDCtx struct{}
