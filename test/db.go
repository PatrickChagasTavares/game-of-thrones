package test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

// GetDB retorna uma instancia do db mocado
func GetDB() (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New(sqlmock.MonitorPingsOption(true))
	return sqlx.NewDb(db, "postgres"), mock
}

// NewRows retorna um modelo para adicionar rows
func NewRows(columns ...string) *sqlmock.Rows {
	return sqlmock.NewRows(columns)
}
