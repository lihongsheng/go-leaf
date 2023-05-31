package pkg

import (
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"go-leaf/internal/conf"
	"go-leaf/internal/types"
	"testing"
	"time"
)

func TestZapLog_Log(t *testing.T) {
	fmt.Println(time.Now().Format("2006-01-02T15:04:05.000000Z"))
	fmt.Println(time.Local)
	l := log.With(NewZapFileLog(conf.Log{}),
		types.LogEnv, "local",
		types.LogServerName, "message-center",
		types.LogServerVersion, "1.0.0",
		types.LogTraceID, tracing.TraceID(),
		types.LogSpanID, tracing.SpanID(),
		types.LogGrpcMethod, GrpcPath(),
		types.LogRequestURL, HTTPPath(),
		types.LogPodName, "local")
	l.Log(log.LevelError, log.DefaultMessageKey, "tag", "err", "err")
}

func TestNewHelper_Info(t *testing.T) {
	l := log.With(NewZapFileLog(conf.Log{}))
	h := NewHelper(l)
	err := errors.New("is error")
	h.Info("tag", "err", "err")
	h.Error("tag", "err", err.Error())
	h.Warn("tag", "err", "error")
}
