package swagger

import (
	"github.com/PatrickChagastavares/game-of-thrones/docs"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/httpRouter"
	httpSwagger "github.com/swaggo/http-swagger"
)

func New(router httpRouter.Router) {

	docs.SwaggerInfo.Title = "Swagger about router of solidAPI"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router.Get("swagger/*any", router.ParseHandler(
		httpSwagger.Handler(httpSwagger.URL("doc.json")),
	))
}
