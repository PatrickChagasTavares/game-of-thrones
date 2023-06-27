package test

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
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

type TestContainer struct {
	testcontainers.Container
	URI string
}

var posgresMux = new(sync.Mutex)

func StartTestContainer(ctx context.Context) (*TestContainer, error) {
	posgresMux.Lock()
	defer posgresMux.Unlock()
	var env = map[string]string{
		"POSTGRES_PASSWORD": "postgres",
		"POSTGRES_USER":     "postgres",
		"POSTGRES_DB":       "game-of-thrones",
	}
	var port = "5432/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres",
			ExposedPorts: []string{port},
			Env:          env,
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %v", err)
	}

	ip, err := container.Host(ctx)
	if err != nil {
		err = fmt.Errorf("failed to get mongo container ip: %w", err)
		return nil, err
	}

	p, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, fmt.Errorf("failed to get container external port: %v", err)
	}

	log.Println("postgres container ready and running at port: ", p.Port())

	return &TestContainer{
		Container: container,
		URI: fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			env["POSTGRES_USER"],
			env["POSTGRES_PASSWORD"],
			ip,
			p.Port(),
			env["POSTGRES_DB"],
		),
	}, nil

}
