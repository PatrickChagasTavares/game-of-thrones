package houses

import (
	"github.com/PatrickChagastavares/game-of-thrones/internal/controllers"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/httpRouter"
)

func New(router httpRouter.Router, Ctrl *controllers.Container) {

	router.Post("/houses", Ctrl.House.Create)
	router.Get("/houses", Ctrl.House.Find)
	router.Get("/houses/:id", Ctrl.House.FindByID)
	router.Put("/houses/:id", Ctrl.House.Update)
	router.Delete("/houses/:id", Ctrl.House.Delete)

}
