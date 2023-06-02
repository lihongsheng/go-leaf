package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	adminv1 "go-leaf/api/admin/v1"
	g1 "go-leaf/api/general/v1"
	"go-leaf/internal/pkg"
	"go-leaf/internal/service"

	"go-leaf/internal/conf"
	"time"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(
	c conf.Server,
	logger log.Logger,
	general *service.General,
	admin *service.Admin,
) *http.Server {
	var timeOut time.Duration
	if c.TimeOut > 0 {
		timeOut = time.Duration(c.TimeOut)
	}
	srv := pkg.NewHTTPServer(logger, c.HTTPPort, timeOut)
	g1.RegisterGeneralHTTPServer(srv, general)
	adminv1.RegisterAdminHTTPServer(srv, admin)
	//// 增加auth
	//auth.RegisterHTTP(srv)
	return srv
}
