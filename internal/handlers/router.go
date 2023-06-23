package handlers

import (
	"github.com/PatrickChagastavares/game-of-thrones/internal/controllers"
	"github.com/PatrickChagastavares/game-of-thrones/internal/handlers/swagger"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/httpRouter"
)

type (
	Options struct {
		Ctrl   *controllers.Container
		Router httpRouter.Router
	}
)

func NewRouter(opts Options) {
	swagger.New(opts.Router)
}
