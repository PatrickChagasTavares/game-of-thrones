DATABASE_CONNECT="postgres://postgres:postgres@127.0.0.1:5432/game-of-thrones?sslmode=disable"
MIGRATION_SOURCE="file://migrations"

.PHONY: setup
setup:
	@echo "installing swaggo..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "installing golang-migrate..."
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "installing mockgen..."
	@go install github.com/golang/mock/mockgen@latest
	@echo "downloading project dependencies..."
	@go mod tidy

.PHONY: run
run:
	@cd cmd && env=local go run main.go

.PHONY: docs
docs:
	@swag init --parseDependency -g cmd/main.go

mocks:
	@go generate ./... 

.PHONY: up-local
up-local:
	@docker compose up -d --build

down-local:
	@docker compose down

.PHONY: test
test: ## runing unit tests with covarage
	GOARCH=amd64 go test ./internal/... -cover -failfast -coverprofile=coverage.out
	@go tool cover -func coverage.out | awk 'END{print sprintf("coverage: %s", $$3)}'

test-cover: test ## runing unit tests with covarage and opening cover profile on browser
	go tool cover -html coverage.out

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