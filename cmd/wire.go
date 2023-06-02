//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
  "github.com/go-kratos/kratos/v2"
  "github.com/go-kratos/kratos/v2/log"
  "github.com/google/wire"
  "go-leaf/internal/conf"
  "go-leaf/internal/logic"
  "go-leaf/internal/pkg"
  "go-leaf/internal/server"
  "go-leaf/internal/service"
)

// wireApp init kratos application.
func WireApp(conf.Server, conf.Conf, log.Logger) (*kratos.App, func(), error) {
  panic(wire.Build(
    server.ProviderSet,
    logic.ProviderSet,
    service.ProviderSet,
    pkg.NewPprof,
    pkg.NewHelper,
    NewApp))
}
