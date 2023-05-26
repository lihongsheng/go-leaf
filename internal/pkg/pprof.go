package pkg

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"message-center/internal/conf"
	"message-center/internal/types"
	"net"
	"net/http"
	"net/http/pprof"
	"time"
)

type Pprof struct {
	http   *http.Server
	config conf.Bootstrap
	log    log.Logger
}

func NewPprof(c conf.Bootstrap, log log.Logger) *Pprof {
	return &Pprof{
		config: c,
		log:    log,
	}
}

// Start pprof
func (p *Pprof) Start(ctx context.Context) error {
	addr := ":0"
	if !p.config.IsLocal() {
		addr = p.config.Pprof
	}
	mux := http.NewServeMux()
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      mux, // 用原生的 serverHttp 接口
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
	}
	mux.Handle(types.PromHTTPHandlerPath, promhttp.Handler())
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("ok"))
	})
	p.http = httpServer
	if p.config.IsLocal() {
		con, err := net.Listen("tcp", addr)
		if err != nil {
			return err
		}
		if p.log != nil {
			_ = p.log.Log(log.LevelInfo, log.DefaultMessageKey, "pprof", "pprof listening at", con.Addr().String())
		}
		err = httpServer.Serve(con)
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
	if p.log != nil {
		_ = p.log.Log(log.LevelInfo, log.DefaultMessageKey, "pprof", "pprof listening at", addr)
	}
	err := httpServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (p *Pprof) Stop(ctx context.Context) error {
	if p.http != nil {
		_ = p.http.Shutdown(ctx)
	}
	return nil
}
