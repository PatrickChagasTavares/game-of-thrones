# Game of Thrones API
[![Test and coverage](https://github.com/PatrickChagastavares/game-of-thrones/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/PatrickChagastavares/game-of-thrones/actions/workflows/tests.yml)
[![codecov](https://codecov.io/gh/PatrickChagastavares/game-of-thrones/branch/main/graph/badge.svg?token=QB8MFFU8XL)](https://codecov.io/gh/PatrickChagastavares/game-of-thrones)

Esse projeto é uma API para armazenar as informações sobre as familias do jogo Game of thrones.

## 🚀 Começando

Essas instruções permitirão que você obtenha uma cópia do projeto em operação na sua máquina local para fins de desenvolvimento e teste.

### 📋 Pré-requisitos

Ferramentas:

- [Docker](https://docs.docker.com/engine/install/)
- [Golang](https://golang.org/doc/install)

## 🛠 Configurando ambiente

Antes de rodar qualquer coisa no projeto recomendo executar esse comando baixo:

- `make setup`: instala todas as dependencias necessarias para o projeto rodar.

## 📦 Desenvolvimento

Alguns comandos importantes para rodar o projeto e validar:

- `make up-local`: Inicia o docker compose (db and admin).
- `make run`: Wrapper para o `cd cmd && env=local go run main.go`.
- `make down-local`: encerra o docker-compose (db and admin).
- `make docker-up`: Inicia o projeto via docker.
- `make docker-down`: Encerra todos os componentes do docker-compose.
- `make docs`: Cria/atualiza os arquivos do swagger.
- `make mocks`: Cria/atualiza os arquivos do mock do projeto.
- `make test`: Rota dos testes do projeto e mostra o cover
- `make test-cover`: O mesmo do `make test`, porém abre o brawser para mais detalhes.
- `make itest`: Rota dos testes de e2e|integração do projeto


## 🗂 Arquitetura

### Descrição dos diretórios e arquivos mais importantes:

- `./cmd/main.go`: O codígo que inicia a aplicação.
- `./config`: Esse diretório possui todos os arquivos para ler as variaveis do projeto.
- `./docs`: Arquivos gerados pelo swagger, referente a documentação.
- `./internal`: O codígo relacionado a aplicação.
- `./internal/handles/**`: Esse diretório possui o registro todas as rotas existentes.
- `./internal/controllers/**`: Esse diretório possui toda as logica volta a camada de handles.
- `./internal/entities/**`: Este diretório possui todos os arquivos de modelos globais do projeto
- `./internal/services/**`: Esse diretório contem toda a regra de negocio da aplicação.
- `./internal/repositories/**`: Esse diretório possui todos os arquivos relacionado a banco ou cache.
- `./migrations`: Esse diretório possui todas migrations o projeto necessita para funcionar.
- `./pkg`: Esse diretório contem todos o pacotes externos que usamos (Gin, Log, Migration e etc...)
- `./test`: Sub-modulos necessários para manutenção do projeto em geral.


## 🛠️ Construído com

- [Gin](https://gin-gonic.com) - Framework Web
- [Go mod](https://blog.golang.org/using-go-modules) - Dependência
- [Docker](https://docs.docker.com) - Container
- [Viper](https://github.com/spf13/viper) - Configuração
- [Migrate](https://github.com/golang-migrate/migrate) - Database migrations
- [Logrus](github.com/sirupsen/logrus) - Log
- [Validator](github.com/go-playground/validator/v10) - Validador de structs
- [Test Container](https://golang.testcontainers.org/) - Run docker by code to integration/e2e test
- [Mockgen](https://github.com/golang/mock) - Gerador de mock de interface
