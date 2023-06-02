package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	adminv1 "go-leaf/api/admin/v1"
	g1 "go-leaf/api/general/v1"
	"go-leaf/internal/conf"
	"go-leaf/internal/pkg"
	"go-leaf/internal/service"
	"time"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c conf.Server,
	logger log.Logger,
	general *service.General,
	admin *service.Admin,
) *grpc.Server {
	var timeOut time.Duration
	if c.TimeOut > 0 {
		timeOut = time.Duration(c.TimeOut)
	}
	srv := pkg.NewGrpcServer(logger, c.GrpcPort, timeOut)
	// 注册grpc 路由
	g1.RegisterGeneralServer(srv, general)
	adminv1.RegisterAdminServer(srv, admin)
	//
	//// 注册auth
	//auth.RegisterGrpc(srv)
	return srv
}
