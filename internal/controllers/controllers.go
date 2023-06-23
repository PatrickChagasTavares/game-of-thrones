package controllers

import (
	"github.com/PatrickChagastavares/game-of-thrones/internal/services"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/logger"
)

type (
	Container struct {
	}

	Options struct {
		Srv *services.Container
		Log logger.Logger
	}
)

func New(opts Options) *Container {
	return &Container{}
}
