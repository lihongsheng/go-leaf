package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"go-leaf/internal/conf"
	"go-leaf/internal/pkg"
	"time"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c conf.Server,
	logger log.Logger,
) *grpc.Server {
	var timeOut time.Duration
	if c.TimeOut > 0 {
		timeOut = time.Duration(c.TimeOut)
	}
	srv := pkg.NewGrpcServer(logger, c.GrpcPort, timeOut)
	//// 注册grpc 路由
	//messagev1.RegisterEmailServer(srv, email)
	//messagev1.RegisterSmsServer(srv, sms)
	//messagev1.RegisterPushServer(srv, push)
	//templatev1.RegisterTemplateServer(srv, template)
	//appv1.RegisterAppServer(srv, app)
	//channelv1.RegisterChannelServer(srv, channel)
	//riskv1.RegisterBlacklistServer(srv, blacklist)
	//logv1.RegisterLogServer(srv, log)
	//
	//// 注册auth
	//auth.RegisterGrpc(srv)
	return srv
}
