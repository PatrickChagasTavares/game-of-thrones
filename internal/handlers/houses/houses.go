package houses

import (
	"github.com/PatrickChagastavares/game-of-thrones/internal/controllers"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/httpRouter"
)

func New(router httpRouter.Router, Ctrl *controllers.Container) {

	router.Post("/houses", Ctrl.House.Create)
}
