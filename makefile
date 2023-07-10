.PHONY: setup up-local down-local run docker-up docker-down docs mocks test itest test-cover migration-create migration-up migration-down

DATABASE_CONNECT="postgres://postgres:postgres@127.0.0.1:5432/game-of-thrones?sslmode=disable"
MIGRATION_SOURCE="file://migrations"
NAME_IMAGE = "game-of-thrones"

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

build: 
	@echo $(NAME_IMAGE): Compilando o micro-serviço
	@env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOFLAGS=-buildvcs=false go build -o dist/$(NAME_IMAGE)/main cmd/main.go 
	@echo $(NAME_IMAGE): Construindo a imagem
	@docker build -t $(NAME_IMAGE) .

docker-up:
	@docker compose -f "docker/docker-compose.yml" up -d --build

docker-down:
	@docker compose -f "docker/docker-compose.yml" down

up-local:
	@docker compose -f "docker/db/docker-compose.yml" up -d --build
	@docker compose -f "docker/observability/docker-compose.yml" up -d --build

down-local:
	@docker compose -f "docker/db/docker-compose.yml" down
	@docker compose -f "docker/observability/docker-compose.yml" down

run: up-local
	@cd cmd && env=local go run main.go

docs:
	@swag init --parseDependency -g cmd/main.go

mocks:
	@go generate ./... 

test: ## runing unit tests with covarage
	GOARCH=amd64 go test ./internal/... -cover -failfast -coverprofile=coverage.out
	@go tool cover -func coverage.out | awk 'END{print sprintf("coverage: %s", $$3)}'

itest: ## run unit and integration tests with covarage
	go test ./test/e2e/... -coverprofile=coverage.out -cover --tags=integration 

test-cover: test ## runing unit tests with covarage and opening cover profile on browser
	go tool cover -html coverage.out

# command to generate migration
migration-create:
	migrate create -ext sql -dir migrations -seq $(NAME)

migration-up:
	migrate -source $(MIGRATION_SOURCE) -database $(DATABASE_CONNECT) --verbose up

migration-down:
	migrate -source $(MIGRATION_SOURCE) -database $(DATABASE_CONNECT) --verbose down 1