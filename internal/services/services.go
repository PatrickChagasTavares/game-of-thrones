package services

import (
	"github.com/PatrickChagastavares/game-of-thrones/internal/repositories"
	"github.com/PatrickChagastavares/game-of-thrones/internal/services/houses"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/logger"
)

type (
	Container struct {
		House houses.IService
	}

	Options struct {
		Repo *repositories.Container
		Log  logger.Logger
	}
)

func New(opts Options) *Container {
	return &Container{
		House: houses.New(opts.Repo, opts.Log),
	}
}
