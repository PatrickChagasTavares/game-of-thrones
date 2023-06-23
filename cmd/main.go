package main

import (
	"github.com/PatrickChagastavares/game-of-thrones/config"
	"github.com/PatrickChagastavares/game-of-thrones/internal/controllers"
	"github.com/PatrickChagastavares/game-of-thrones/internal/handlers"
	"github.com/PatrickChagastavares/game-of-thrones/internal/repositories"
	"github.com/PatrickChagastavares/game-of-thrones/internal/services"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/httpRouter"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/logger"
	migration "github.com/PatrickChagastavares/game-of-thrones/pkg/migrations"
	"github.com/jmoiron/sqlx"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	var log = logger.NewLogrusLogger()
	configs, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("failed to read config: ", err)
		return
	}

	migration.RunMigrations(configs.Database.Writer)

	var (
		router       = httpRouter.NewGinRouter()
		repositories = repositories.New(repositories.Options{
			WriterSqlx: sqlx.MustConnect("postgres", configs.Database.Writer),
			ReaderSqlx: sqlx.MustConnect("postgres", configs.Database.Reader),
			Log:        log,
		})
		services = services.New(services.Options{
			Repo: repositories,
			Log:  log,
		})
		controllers = controllers.New(controllers.Options{
			Srv: services,
			Log: log,
		})
	)

	handlers.NewRouter(handlers.Options{
		Router: router,
		Ctrl:   controllers,
	})

	log.Info("start serve in port:", configs.Port)
	if err := router.Server(configs.Port); err != nil {
		log.Fatal(err)
	}
}
