package repositories

import (
	"github.com/PatrickChagastavares/game-of-thrones/internal/repositories/database/houses"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type (
	// Container model to export instance repositories
	Container struct {
		Database SqlContainer
	}

	SqlContainer struct {
		House houses.IRepository
	}

	// Options struct of options to create a new repositories
	Options struct {
		WriterSqlx *sqlx.DB
		ReaderSqlx *sqlx.DB
		Log        logger.Logger
	}
)

// New Create a new instance of repositories
func New(opts Options) *Container {
	return &Container{
		Database: SqlContainer{
			House: houses.NewSqlx(opts.Log, opts.WriterSqlx, opts.ReaderSqlx),
		},
	}
}
