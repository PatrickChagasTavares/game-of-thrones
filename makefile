DATABASE_CONNECT="postgres://postgres:postgres@127.0.0.1:5432/game-of-thrones?sslmode=disable"
MIGRATION_SOURCE="file://migrations"

.PHONY: setup
setup:
	@echo "installing swaggo..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "installing golang-migrate..."
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "downloading project dependencies..."
	@go mod tidy

.PHONY: run
run:
	@cd cmd && env=local go run main.go

.PHONY: docs
docs:
	@swag init -g cmd/main.go

# command to generate migration
.PHONY: migration-create
migration-create:
	migrate create -ext sql -dir migrations -seq $(NAME)

.PHONY: migration-up
migration-up:
	migrate -source $(MIGRATION_SOURCE) -database $(DATABASE_CONNECT) --verbose up

.PHONY: migration-down
migration-down:
	migrate -source $(MIGRATION_SOURCE) -database $(DATABASE_CONNECT) --verbose down 1