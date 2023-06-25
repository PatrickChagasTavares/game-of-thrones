package characters

import (
	"github.com/PatrickChagastavares/game-of-thrones/internal/controllers"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/httpRouter"
)

func New(router httpRouter.Router, Ctrl *controllers.Container) {

	router.Post("/characters", Ctrl.Character.Create)
	router.Get("/characters", Ctrl.Character.Find)
	router.Get("/characters/:id", Ctrl.Character.FindByID)
	router.Put("/characters/:id", Ctrl.Character.Update)
	router.Delete("/characters/:id", Ctrl.Character.Delete)

}
