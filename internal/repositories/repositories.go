package repositories

import (
	"github.com/PatrickChagastavares/game-of-thrones/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type (
	// Container model to export instance repositories
	Container struct {
		Database SqlContainer
	}

	SqlContainer struct {
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
		Database: SqlContainer{},
	}
}
