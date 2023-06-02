package logic

import (
	"go-leaf/internal/conf"
	"go-leaf/internal/pkg"
)

type General struct {
	config conf.Conf
	logger *pkg.Helper
}

func NewGeneral(c conf.Conf, logger *pkg.Helper) *General {
	return &General{
		config: c,
		logger: logger,
	}
}
