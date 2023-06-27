package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/PatrickChagastavares/game-of-thrones/internal/controllers"
	"github.com/PatrickChagastavares/game-of-thrones/internal/entities"
	"github.com/PatrickChagastavares/game-of-thrones/internal/handlers"
	"github.com/PatrickChagastavares/game-of-thrones/internal/repositories"
	"github.com/PatrickChagastavares/game-of-thrones/internal/services"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/httpRouter"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/logger"
	migration "github.com/PatrickChagastavares/game-of-thrones/pkg/migrations"
	"github.com/PatrickChagastavares/game-of-thrones/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var secoundSleep = time.Second * 1

func Test_housesE2E(t *testing.T) {
	ctx := context.Background()

	// ============ Start Database by docker ============
	container, err := test.StartTestContainer(ctx)
	if err != nil {
		t.Error(err)
	}
	defer container.Terminate(ctx)

	// ============ RUN Migrations ============
	time.Sleep(secoundSleep)
	os.Setenv("env", "local")
	migration.RunMigrations(container.URI)

	// ============ Start Internal Modules ============
	var (
		log          = logger.NewLogrusLogger()
		router       = httpRouter.NewGinRouter()
		repositories = repositories.New(repositories.Options{
			WriterSqlx: sqlx.MustConnect("postgres", container.URI),
			ReaderSqlx: sqlx.MustConnect("postgres", container.URI),
			Log:        log,
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
		Ctrl:   controllers,
		Router: router,
	})

	// ============ RUN Server http ============
	go router.Server(":3003")
	time.Sleep(secoundSleep)

	// ============ Create Variable to reuse ============
	var houseID, lordID string

	// ============ Start Cases of Tests ============
	t.Run("Should to Create Lord", func(t *testing.T) {
		payloadCharacter := entities.CharacterRequest{
			Name:     "Patrick Chagas",
			TVSeries: []string{"session 1", "session 2"},
		}
		resp := map[string]any{}
		err := request(ctx, http.MethodPost, "/characters", payloadCharacter, &resp)

		assert.Nil(t, err)

		lordID = resp["id"].(string)
	})

	t.Run("Should to Create house", func(t *testing.T) {
		payloadHouse := entities.HouseRequest{
			Name:           "House Patrick",
			Region:         "são paulo",
			FoundationYear: "2023",
			CurrentLord:    lordID,
		}
		resp := map[string]any{}
		err := request(ctx, http.MethodPost, "/houses", payloadHouse, &resp)

		assert.Nil(t, err)

		houseID = resp["id"].(string)
	})

	t.Run("Should to Create secound house", func(t *testing.T) {
		payloadHouse := entities.HouseRequest{
			Name:           "House patrick secound",
			Region:         "floripa",
			FoundationYear: "2023",
			CurrentLord:    lordID,
		}
		resp := map[string]any{}
		err := request(ctx, http.MethodPost, "/houses", payloadHouse, &resp)

		assert.Nil(t, err)
	})

	t.Run("Should find house created", func(t *testing.T) {
		resp := entities.House{}

		err := request(ctx, http.MethodGet, "/houses/"+houseID, nil, &resp)

		assert.Nil(t, err)
		assert.Equal(t, houseID, resp.ID)
	})

	t.Run("Should find lord created", func(t *testing.T) {
		resp := entities.Character{}

		err := request(ctx, http.MethodGet, "/characters/"+lordID, nil, &resp)

		assert.Nil(t, err)
		assert.Equal(t, lordID, resp.ID)
	})

	t.Run("Should Update house created", func(t *testing.T) {
		req := entities.HouseRequest{
			Name:           "House Chagas",
			Region:         "são paulo",
			FoundationYear: "2023",
			CurrentLord:    lordID,
		}
		resp := entities.House{}

		err := request(ctx, http.MethodPut, "/houses/"+houseID, req, &resp)

		assert.Nil(t, err)
		assert.Equal(t, houseID, resp.ID)
		assert.Equal(t, req.Name, resp.Name)
	})

	t.Run("Should Update lord created", func(t *testing.T) {
		payloadCharacter := entities.CharacterRequest{
			Name:     "Patrick Chagas Tavares",
			TVSeries: []string{"session 1"},
		}
		resp := entities.Character{}

		err := request(ctx, http.MethodPut, "/characters/"+lordID, payloadCharacter, &resp)

		assert.Nil(t, err)
		assert.Equal(t, lordID, resp.ID)
		assert.Equal(t, payloadCharacter.Name, resp.Name)
		assert.Equal(t, payloadCharacter.TVSeries, resp.TVSeries)
	})

	t.Run("Should Find all Houses", func(t *testing.T) {
		resp := []entities.House{}

		err := request(ctx, http.MethodGet, "/houses", nil, &resp)

		assert.Nil(t, err)
		assert.Equal(t, lordID, resp[0].CurrentLord)
		assert.Equal(t, lordID, resp[1].CurrentLord)
	})

	t.Run("Should DeleteLord and remove currentlord of houses", func(t *testing.T) {
		errLord := request(ctx, http.MethodDelete, "/characters/"+lordID, nil, nil)

		assert.Nil(t, errLord)

		resp := []entities.House{}
		errHouses := request(ctx, http.MethodGet, "/houses", nil, &resp)

		assert.Nil(t, errHouses)
		assert.Equal(t, "", resp[0].CurrentLord)
		assert.Equal(t, "", resp[1].CurrentLord)
	})

	t.Run("Should Delete house created", func(t *testing.T) {
		err := request(ctx, http.MethodDelete, "/houses/"+houseID, nil, nil)

		assert.Nil(t, err)
	})
}

const baseURL = "http://localhost:3003"

func request(ctx context.Context, method, endpoint string, body interface{}, data interface{}) error {
	requestBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	url := baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	if res.StatusCode < 200 || res.StatusCode > 299 {
		detail := make(map[string]interface{})
		if err := decoder.Decode(&detail); err != nil {
			return err
		}

		return nil
	}

	if data != nil {
		err = decoder.Decode(data)
		if err != nil {
			return err
		}
	}

	return nil
}
