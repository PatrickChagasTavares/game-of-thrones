package characters

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PatrickChagastavares/game-of-thrones/internal/entities"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/logger"
	"github.com/PatrickChagastavares/game-of-thrones/test"

	"github.com/stretchr/testify/assert"
)

func Test_Create(t *testing.T) {
	data := entities.CharacterRequest{
		ID:        "id_123",
		Name:      "teste Patrick",
		TVSeries:  []string{"session 1", "session 2"},
		CreatedAt: time.Now(),
	}

	cases := map[string]struct {
		input       entities.CharacterRequest
		expectedErr error

		prepareMock func(mock sqlmock.Sqlmock)
	}{
		"Should return success": {
			input: data,
			prepareMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`INSERT INTO characters 
				(id,name,tv_series,created_at)
				VALUES ($1, $2, $3, $4);`)
				mock.ExpectExec(query).
					WithArgs(data.ID, data.Name, data.TVSeries, data.CreatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		"Should return Error": {
			input:       data,
			expectedErr: errors.New("problem to create character"),
			prepareMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`INSERT INTO characters 
				(id,name,tv_series,created_at)
				VALUES ($1, $2, $3, $4);`)
				mock.ExpectExec(query).
					WithArgs(data.ID, data.Name, data.TVSeries, data.CreatedAt).
					WillReturnError(errors.New("Problem to execute query"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			db, mock := test.GetDB()

			cs.prepareMock(mock)

			repo := NewSqlx(logger.NewLogrusLogger(), db, db)

			err := repo.Create(context.Background(), cs.input)

			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_Find(t *testing.T) {
	resp := []entities.Character{
		{ID: "id_1", Name: "Patrick", TVSeries: []string{"session 1", "session 2"}},
		{ID: "id_2", Name: "Patrick", TVSeries: []string{"session 2", "session 2"}},
	}

	cases := map[string]struct {
		expectedData []entities.Character
		expectedErr  error

		prepareMock func(mock sqlmock.Sqlmock)
	}{
		"Should return success": {
			expectedData: resp,
			prepareMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
				SELECT id, name, tv_series, created_at, updated_at
				FROM characters
				WHERE deleted_at is null
				ORDER BY created_at DESC;
				`)
				rows := test.NewRows("id", "name", "tv_series", "created_at", "updated_at").
					AddRow(resp[0].ID, resp[0].Name, resp[0].TVSeries, resp[0].CreatedAt, nil).
					AddRow(resp[1].ID, resp[1].Name, resp[1].TVSeries, resp[1].CreatedAt, nil)
				mock.ExpectQuery(query).
					WillReturnRows(rows)
			},
		},
		"Should return success without rows": {
			expectedData: []entities.Character{},
			prepareMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
				SELECT id, name, tv_series, created_at, updated_at
				FROM characters
				WHERE deleted_at is null
				ORDER BY created_at DESC;
				`)
				mock.ExpectQuery(query).
					WillReturnError(sql.ErrNoRows)
			},
		},
		"Should return Error": {
			expectedErr: errors.New("problem to find characters"),
			prepareMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
				SELECT id, name, tv_series, created_at, updated_at
				FROM characters
				WHERE deleted_at is null
				ORDER BY created_at DESC;
				`)
				mock.ExpectExec(query).
					WillReturnError(errors.New("Problem to execute query"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			db, mock := test.GetDB()

			cs.prepareMock(mock)

			repo := NewSqlx(logger.NewLogrusLogger(), db, db)

			data, err := repo.Find(context.Background())

			assert.Equal(t, cs.expectedErr, err)
			assert.Equal(t, cs.expectedData, data)
		})
	}
}

func Test_FindByID(t *testing.T) {
	resp := entities.Character{
		ID: "id_1", Name: "Patrick", TVSeries: []string{"session 1", "session 2"},
	}

	cases := map[string]struct {
		input        string
		expectedData entities.Character
		expectedErr  error

		prepareMock func(mock sqlmock.Sqlmock)
	}{
		"Should return success": {
			input:        resp.ID,
			expectedData: resp,
			prepareMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
				SELECT id, name, tv_series, created_at, updated_at
				FROM characters
				WHERE id =$1 AND deleted_at is null;`)
				rows := test.NewRows("id", "name", "tv_series", "created_at", "updated_at").
					AddRow(resp.ID, resp.Name, resp.TVSeries, resp.CreatedAt, nil)
				mock.ExpectQuery(query).
					WithArgs(resp.ID).
					WillReturnRows(rows)
			},
		},
		"Should return Error": {
			input:       resp.ID,
			expectedErr: errors.New("character is not found or deleted"),
			prepareMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
				SELECT id, name, tv_series, created_at, updated_at
				FROM characters
				WHERE id =$1 AND deleted_at is null;`)
				mock.ExpectExec(query).
					WithArgs(resp.ID).
					WillReturnError(errors.New("Problem to execute query"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			db, mock := test.GetDB()

			cs.prepareMock(mock)

			repo := NewSqlx(logger.NewLogrusLogger(), db, db)

			data, err := repo.FindByID(context.Background(), cs.input)

			assert.Equal(t, cs.expectedErr, err)
			assert.Equal(t, cs.expectedData, data)
		})
	}
}

func Test_Update(t *testing.T) {
	now := time.Now()
	resp := entities.Character{
		ID:        "id_1",
		Name:      "Patrick",
		TVSeries:  []string{"session 1", "session 2"},
		CreatedAt: now,
		UpdatedAt: &now,
	}

	cases := map[string]struct {
		input       *entities.Character
		expectedErr error

		prepareMock func(mock sqlmock.Sqlmock)
	}{
		"Should return success": {
			input: &resp,
			prepareMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
				UPDATE characters
				SET name = $1, tv_series = $2, updated_at = $3
				WHERE id = $4;
				`)
				mock.ExpectExec(query).
					WithArgs(resp.Name, resp.TVSeries, resp.UpdatedAt, resp.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		"Should return Error": {
			input:       &resp,
			expectedErr: errors.New("failed to update character"),
			prepareMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
				UPDATE characters
				SET name = $1, tv_series = $2, updated_at = $3
				WHERE id = $4;
				`)
				mock.ExpectExec(query).
					WithArgs(resp.Name, resp.TVSeries, resp.UpdatedAt, resp.ID).
					WillReturnError(errors.New("Problem to execute query"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			db, mock := test.GetDB()

			cs.prepareMock(mock)

			repo := NewSqlx(logger.NewLogrusLogger(), db, db)

			err := repo.Update(context.Background(), cs.input)

			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_Delete(t *testing.T) {
	id := "id_123"
	now := time.Now()
	timeNow = func() time.Time {
		return now
	}

	cases := map[string]struct {
		input       string
		expectedErr error

		prepareMock func(mock sqlmock.Sqlmock)
	}{
		"Should return success": {
			input: id,
			prepareMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
				UPDATE characters
				SET deleted_at = $1
				WHERE id = $2;
				`)
				mock.ExpectExec(query).
					WithArgs(now, id).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		"Should return Error": {
			input:       id,
			expectedErr: errors.New("failed to delete character"),
			prepareMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
				UPDATE characters
				SET deleted_at = $1
				WHERE id = $2;
				`)
				mock.ExpectExec(query).
					WithArgs(now, id).
					WillReturnError(errors.New("Problem to execute query"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			db, mock := test.GetDB()

			cs.prepareMock(mock)

			repo := NewSqlx(logger.NewLogrusLogger(), db, db)

			err := repo.Delete(context.Background(), cs.input)

			assert.Equal(t, cs.expectedErr, err)
		})
	}
}
