package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go-leaf/internal/pkg"

	"go-leaf/internal/conf"
	"time"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(
	c conf.Server, logger log.Logger,
) *http.Server {
	var timeOut time.Duration
	if c.TimeOut > 0 {
		timeOut = time.Duration(c.TimeOut)
	}
	srv := pkg.NewHTTPServer(logger, c.HTTPPort, timeOut)
	//messagev1.RegisterSmsHTTPServer(srv, sms)
	//messagev1.RegisterEmailHTTPServer(srv, email)
	//messagev1.RegisterPushHTTPServer(srv, push)
	//templatev1.RegisterTemplateHTTPServer(srv, template)
	//appv1.RegisterAppHTTPServer(srv, app)
	//reportv1.RegisterCallbackHTTPServer(srv, callback)
	//cahnnelv1.RegisterChannelHTTPServer(srv, channel)
	//logv1.RegisterLogHTTPServer(srv, log)
	//riskv1.RegisterBlacklistHTTPServer(srv, blacklist)
	//// 增加auth
	//auth.RegisterHTTP(srv)
	return srv
}
