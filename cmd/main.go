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
	"github.com/PatrickChagastavares/game-of-thrones/pkg/tracer"
	tracerjaeger "github.com/PatrickChagastavares/game-of-thrones/pkg/tracer/tracer_jaeger"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

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

	trace := tracer.New(tracerjaeger.NewExporter(configs.Tracer))
	defer trace.Close()

	var (
		router       = httpRouter.NewGinRouter()
		repositories = repositories.New(repositories.Options{
			WriterSqlx: otelsqlx.MustConnect(
				"postgres",
				configs.Database.Writer,
				otelsql.WithAttributes(
					semconv.DBSystemPostgreSQL,
					semconv.DBName("game-of-thrones"),
				)),
			ReaderSqlx: otelsqlx.MustConnect(
				"postgres",
				configs.Database.Reader,
				otelsql.WithAttributes(
					semconv.DBSystemPostgreSQL,
					semconv.DBName("game-of-thrones"),
				)),
			Log: log,
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
