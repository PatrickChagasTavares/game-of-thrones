package migration

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

func RunMigrations(dbURL string) {
	fmt.Println("start migration")
	if err := getMigration(dbURL).Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}

func getMigration(dbURL string) *migrate.Migrate {
	dir, _ := os.Getwd()
	if os.Getenv("env") == "local" {
		dir = strings.SplitAfter(dir, "game-of-thrones")[0]
	}
	m, err := migrate.New(
		fmt.Sprintf("file://%s/migrations", dir),
		dbURL,
	)
	if err != nil {
		log.Fatal(err)
	}
	return m
}
