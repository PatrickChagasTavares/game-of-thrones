package services

import (
	"github.com/PatrickChagastavares/game-of-thrones/internal/repositories"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/logger"
)

type (
	Container struct {
	}

	Options struct {
		Repo *repositories.Container
		Log  logger.Logger
	}
)

func New(opts Options) *Container {
	return &Container{}
}
