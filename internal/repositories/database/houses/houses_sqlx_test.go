package houses

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
	data := entities.HouseRequest{
		ID:             "id_123",
		Name:           "house Patrick",
		Region:         "sao paulo",
		FoundationYear: "2023",
		CurrentLord:    "id_1",
		CreatedAt:      time.Now(),
	}

	cases := map[string]struct {
		input       entities.HouseRequest
		expectedErr error

		prepareMock func(mock sqlmock.Sqlmock)
	}{
		"Should return success": {
			input: data,
			prepareMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`INSERT INTO houses 
				(id,name,region,foundation_year,current_lord,created_at)
				VALUES ($1, $2, $3, $4, $5, $6);`)
				mock.ExpectExec(query).
					WithArgs(data.ID, data.Name, data.Region, data.FoundationYear, data.CurrentLord, data.CreatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		"Should return Error": {
			input:       data,
			expectedErr: errors.New("problem to create house"),
			prepareMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`INSERT INTO houses 
				(id,name,region,foundation_year,current_lord,created_at)
				VALUES ($1, $2, $3, $4, $5, $6);`)
				mock.ExpectExec(query).
					WithArgs(data.ID, data.Name, data.Region, data.FoundationYear, data.CurrentLord, data.CreatedAt).
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
	resp := []entities.House{
		{ID: "id_123", Name: "house Patrick", Region: "sao paulo", FoundationYear: "2023", CurrentLord: "id_1", CreatedAt: time.Now()},
		{ID: "id_234", Name: "house chagas ", Region: "sao paulo", FoundationYear: "2023", CurrentLord: "id_2", CreatedAt: time.Now()},
	}

	cases := map[string]struct {
		expectedData []entities.House
		expectedErr  error

		prepareMock func(mock sqlmock.Sqlmock)
	}{
		"Should return success": {
			expectedData: resp,
			prepareMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
				SELECT id, name, region, foundation_year, current_lord, created_at, updated_at
				FROM houses
				WHERE deleted_at is null
				ORDER BY created_at DESC;
				`)
				rows := test.NewRows("id", "name", "region", "foundation_year", "current_lord", "created_at", "updated_at").
					AddRow(resp[0].ID, resp[0].Name, resp[0].Region, resp[0].FoundationYear, resp[0].CurrentLord, resp[0].CreatedAt, nil).
					AddRow(resp[1].ID, resp[1].Name, resp[1].Region, resp[1].FoundationYear, resp[1].CurrentLord, resp[1].CreatedAt, nil)
				mock.ExpectQuery(query).
					WillReturnRows(rows)
			},
		},
		"Should return success without rows": {
			expectedData: []entities.House{},
			prepareMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
				SELECT id, name, region, foundation_year, current_lord, created_at, updated_at
				FROM houses
				WHERE deleted_at is null
				ORDER BY created_at DESC;
				`)
				mock.ExpectQuery(query).
					WillReturnError(sql.ErrNoRows)
			},
		},
		"Should return Error": {
			expectedErr: errors.New("problem to find houses"),
			prepareMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
				SELECT id, name, region, foundation_year, current_lord, created_at, updated_at
				FROM houses
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
	resp := entities.House{
		ID:             "id_123",
		Name:           "house Patrick",
		Region:         "sao paulo",
		FoundationYear: "2023",
		CurrentLord:    "id_1",
		CreatedAt:      time.Now(),
	}

	cases := map[string]struct {
		input        string
		expectedData entities.House
		expectedErr  error

		prepareMock func(mock sqlmock.Sqlmock)
	}{
		"Should return success": {
			input:        resp.ID,
			expectedData: resp,
			prepareMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
				SELECT id, name, region, foundation_year, current_lord, created_at, updated_at
				FROM houses
				WHERE id =$1 AND deleted_at is null;`)
				rows := test.NewRows("id", "name", "region", "foundation_year", "current_lord", "created_at", "updated_at").
					AddRow(resp.ID, resp.Name, resp.Region, resp.FoundationYear, resp.CurrentLord, resp.CreatedAt, nil)
				mock.ExpectQuery(query).
					WithArgs(resp.ID).
					WillReturnRows(rows)
			},
		},
		"Should return Error": {
			input:       resp.ID,
			expectedErr: errors.New("house is not found or deleted"),
			prepareMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
				SELECT id, name, region, foundation_year, current_lord, created_at, updated_at
				FROM houses
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

func Test_FindByName(t *testing.T) {
	resp := entities.House{
		ID:             "id_123",
		Name:           "house Patrick",
		Region:         "sao paulo",
		FoundationYear: "2023",
		CurrentLord:    "id_1",
		CreatedAt:      time.Now(),
	}

	cases := map[string]struct {
		input        string
		expectedData entities.House
		expectedErr  error

		prepareMock func(mock sqlmock.Sqlmock)
	}{
		"Should return success": {
			input:        resp.Name,
			expectedData: resp,
			prepareMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
				SELECT id, name, region, foundation_year, current_lord, created_at, updated_at
				FROM houses
				WHERE name=$1 AND deleted_at is null;`)
				rows := test.NewRows("id", "name", "region", "foundation_year", "current_lord", "created_at", "updated_at").
					AddRow(resp.ID, resp.Name, resp.Region, resp.FoundationYear, resp.CurrentLord, resp.CreatedAt, nil)
				mock.ExpectQuery(query).
					WithArgs(resp.Name).
					WillReturnRows(rows)
			},
		},
		"Should return Error": {
			input:       resp.Name,
			expectedErr: errors.New("house is not found or deleted"),
			prepareMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
				SELECT id, name, region, foundation_year, current_lord, created_at, updated_at
				FROM houses
				WHERE name=$1 AND deleted_at is null;`)
				mock.ExpectExec(query).
					WithArgs(resp.Name).
					WillReturnError(errors.New("Problem to execute query"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			db, mock := test.GetDB()

			cs.prepareMock(mock)

			repo := NewSqlx(logger.NewLogrusLogger(), db, db)

			data, err := repo.FindByName(context.Background(), cs.input)

			assert.Equal(t, cs.expectedErr, err)
			assert.Equal(t, cs.expectedData, data)
		})
	}
}

func Test_RemoveLord(t *testing.T) {
	input := "lord_Id_123"
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
			input: input,
			prepareMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
				UPDATE houses
				SET current_lord = '', updated_at = $1
				WHERE  current_lord= $2;
				`)
				mock.ExpectExec(query).
					WithArgs(now, input).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		"Should return Error": {
			input:       input,
			expectedErr: errors.New("failed to remove current_lord by house"),
			prepareMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
				UPDATE houses
				SET current_lord = '', updated_at = $1
				WHERE  current_lord= $2;
				`)
				mock.ExpectExec(query).
					WithArgs(now, input).
					WillReturnError(errors.New("Problem to execute query"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			db, mock := test.GetDB()

			cs.prepareMock(mock)

			repo := NewSqlx(logger.NewLogrusLogger(), db, db)

			err := repo.RemoveLord(context.Background(), cs.input)

			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_Update(t *testing.T) {
	now := time.Now()
	resp := entities.House{
		ID:             "id_123",
		Name:           "house Patrick",
		Region:         "sao paulo",
		FoundationYear: "2023",
		CurrentLord:    "id_1",
		CreatedAt:      time.Now(),
		UpdatedAt:      &now,
	}

	cases := map[string]struct {
		input       *entities.House
		expectedErr error

		prepareMock func(mock sqlmock.Sqlmock)
	}{
		"Should return success": {
			input: &resp,
			prepareMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
				UPDATE houses
				SET name = $1, region = $2, foundation_year = $3, current_lord = $4, updated_at = $5
				WHERE id = $6;
				`)
				mock.ExpectExec(query).
					WithArgs(resp.Name, resp.Region, resp.FoundationYear, resp.CurrentLord, resp.UpdatedAt, resp.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		"Should return Error": {
			input:       &resp,
			expectedErr: errors.New("failed to update house"),
			prepareMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
				UPDATE houses
				SET name = $1, region = $2, foundation_year = $3, current_lord = $4, updated_at = $5
				WHERE id = $6;
				`)
				mock.ExpectExec(query).
					WithArgs(resp.Name, resp.Region, resp.FoundationYear, resp.CurrentLord, resp.UpdatedAt, resp.ID).
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
				UPDATE houses
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
			expectedErr: errors.New("failed to delete house"),
			prepareMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(`
				UPDATE houses
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
