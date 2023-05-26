package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	_ "github.com/go-sql-driver/mysql"
	"go-leaf/internal/conf"
	"go-leaf/internal/pkg"
	"go-leaf/internal/tools"
	"go-leaf/internal/types"
	_ "go.uber.org/automaxprocs"
	"google.golang.org/protobuf/encoding/protojson"
	"os"
	"strings"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// flagconf is the config flag.
	flagconf string
	id, _    = os.Hostname()
	env      string
	podName  string
)

func init() {
	flag.StringVar(&flagconf, "conf", "/app/configs", "config path, eg: -conf config.yaml;default dockerfile path")
	flag.StringVar(&env, "env", "", "env config tag, eg: -env local|test|product")
}

func main() {
	flag.Parse()
	json.MarshalOptions = protojson.MarshalOptions{
		EmitUnpopulated: true, //默认值不忽略
		UseProtoNames:   true, //使用proto name返回http字段
	}
	if env == "" {
		panic(errors.New("env param is empty"))
	}
	confFile := ""
	flagconf = strings.TrimRight(flagconf, "/") + "/"
	switch env {
	case conf.EnvLocal:
		confFile = flagconf + "config.yaml"
		break
	case conf.EnvTest:
		confFile = flagconf + "config_rpc_test.yaml"
		break
	case conf.EnvProduct:
		confFile = flagconf + "config_rpc_product.yaml"
		break
	}
	c := config.New(
		config.WithSource(
			file.NewSource(confFile),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Conf
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("%+v", bc))
	podName = os.Getenv("HOSTNAME")
	if podName == "" {
		podName = bc.Server.Name + "-" + tools.GenerateID()
	}
	bc.Env = env
	bc.Log.Path += podName + ".log"
	// 使用zap log
	logger := log.With(pkg.NewZapFileLog(bc.Log),
		types.LogServerID, id,
		types.LogServerName, bc.Server.Name,
		types.LogServerVersion, bc.Server.Version,
		types.LogTraceID, tracing.TraceID(),
		types.LogSpanID, tracing.SpanID(),
		types.LogGrpcMethod, pkg.GrpcPath(),
		types.LogRequestURL, pkg.HTTPPath(),
		types.LogEnv, bc.Env,
	)
	// jaeger初始化，用于链路追踪，及后续打日志里的traceID记录
	_ = tools.InitJaeger(bc.JaegerUrl, bc.Server.Name)
	// 新建一个 WithCancel ,告知 Prometheus 结束退出协程
	pkg.ProviderHTTPPrometheus(context.Background(), "go-leaf", "leaf")
	app, cleanup, err := WireApp(bc.Server, bc.Data, bc.Secret, bc, logger)
	defer func() {
		if cleanup != nil {
			cleanup()
		}
	}()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func NewApp(logger log.Logger, gs *grpc.Server, hs *http.Server, c conf.Server, pprof *pkg.Pprof) *kratos.App {
	return kratos.New(
		kratos.ID(""),
		kratos.Name(c.Name),
		kratos.Version(c.Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
			pprof,
		),
	)
}
