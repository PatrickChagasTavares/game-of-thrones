package houses

import (
	"context"
	"errors"
	"testing"

	"github.com/PatrickChagastavares/game-of-thrones/internal/entities"
	"github.com/PatrickChagastavares/game-of-thrones/internal/repositories"
	"github.com/PatrickChagastavares/game-of-thrones/internal/repositories/database/houses"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/logger"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Create(t *testing.T) {
	data := entities.HouseRequest{
		Name:           "house Patrick",
		Region:         "sao paulo",
		FoundationYear: "2023",
		CurrentLord:    "",
	}

	cases := map[string]struct {
		input entities.HouseRequest

		expectedErr error
		prepareMock func(mock *houses.MockIRepository)
	}{
		"Should return success": {
			input: data,
			prepareMock: func(mock *houses.MockIRepository) {
				mock.EXPECT().
					FindByName(gomock.Any(), data.Name).
					Times(1).
					Return(entities.House{}, errors.New("not found"))

				mock.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(entities.HouseRequest{})).
					Times(1).
					Return(nil)
			},
		},
		"Should return error name already used": {
			input:       data,
			expectedErr: ErrNameUsed,
			prepareMock: func(mock *houses.MockIRepository) {
				mock.EXPECT().
					FindByName(gomock.Any(), data.Name).
					Times(1).
					Return(entities.House{}, nil)
			},
		},
		"Should return error": {
			input:       data,
			expectedErr: errors.New("problem to create house"),
			prepareMock: func(mock *houses.MockIRepository) {
				mock.EXPECT().
					FindByName(gomock.Any(), data.Name).
					Times(1).
					Return(entities.House{}, errors.New("not found"))

				mock.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(entities.HouseRequest{})).
					Times(1).
					Return(errors.New("problem to create house"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			mock := houses.NewMockIRepository(ctrl)

			cs.prepareMock(mock)

			srv := New(&repositories.Container{
				Database: repositories.SqlContainer{House: mock}},
				logger.NewLogrusLogger(),
			)

			_, err := srv.Create(ctx, cs.input)

			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_Find(t *testing.T) {
	data := []entities.House{
		{ID: "id_1", Name: "house Patrick", Region: "sao paulo", FoundationYear: "2023", CurrentLord: ""},
		{ID: "id_1", Name: "house Patrick Chagas", Region: "sao paulo", FoundationYear: "2023", CurrentLord: ""},
	}

	cases := map[string]struct {
		input        string
		expectedData []entities.House
		expectedErr  error
		prepareMock  func(mock *houses.MockIRepository)
	}{
		"Should return success with name": {
			input:        "house Patrick",
			expectedData: []entities.House{data[0]},
			prepareMock: func(mock *houses.MockIRepository) {
				mock.EXPECT().
					FindByName(gomock.Any(), data[0].Name).
					Times(1).
					Return(data[0], nil)
			},
		},
		"Should return success without name": {
			expectedData: data,
			prepareMock: func(mock *houses.MockIRepository) {
				mock.EXPECT().
					Find(gomock.Any()).
					Times(1).
					Return(data, nil)
			},
		},
		"Should return error on FindByName": {
			input:       "house Patrick",
			expectedErr: ErrFind,
			prepareMock: func(mock *houses.MockIRepository) {
				mock.EXPECT().
					FindByName(gomock.Any(), data[0].Name).
					Times(1).
					Return(entities.House{}, errors.New("problem to query"))
			},
		},
		"Should return error on Find": {
			expectedErr: ErrFind,
			prepareMock: func(mock *houses.MockIRepository) {
				mock.EXPECT().
					Find(gomock.Any()).
					Times(1).
					Return(nil, errors.New("problem to query"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			mock := houses.NewMockIRepository(ctrl)

			cs.prepareMock(mock)

			srv := New(&repositories.Container{
				Database: repositories.SqlContainer{House: mock}},
				logger.NewLogrusLogger(),
			)

			data, err := srv.Find(ctx, cs.input)

			assert.Equal(t, cs.expectedErr, err)
			assert.Equal(t, cs.expectedData, data)
		})
	}
}

func Test_FindByID(t *testing.T) {
	data := entities.House{
		ID:             "id_1",
		Name:           "Patrick",
		Region:         "sao paulo",
		FoundationYear: "2023",
		CurrentLord:    "",
	}

	cases := map[string]struct {
		input        string
		expectedData entities.House
		expectedErr  error
		prepareMock  func(mock *houses.MockIRepository)
	}{
		"Should return success": {
			input:        data.ID,
			expectedData: data,
			prepareMock: func(mock *houses.MockIRepository) {
				mock.EXPECT().
					FindByID(gomock.Any(), data.ID).
					Times(1).
					Return(data, nil)
			},
		},
		"Should return error": {
			input:       data.ID,
			expectedErr: ErrHouseNotFound,
			prepareMock: func(mock *houses.MockIRepository) {
				mock.EXPECT().
					FindByID(gomock.Any(), data.ID).
					Times(1).
					Return(entities.House{}, errors.New("problem to query"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			mock := houses.NewMockIRepository(ctrl)

			cs.prepareMock(mock)

			srv := New(&repositories.Container{
				Database: repositories.SqlContainer{House: mock}},
				logger.NewLogrusLogger(),
			)

			data, err := srv.FindByID(ctx, cs.input)

			assert.Equal(t, cs.expectedErr, err)
			assert.Equal(t, cs.expectedData, data)
		})
	}
}

func Test_Update(t *testing.T) {
	req := entities.HouseRequest{
		ID:             "id_1",
		Name:           "house Patrick",
		Region:         "sao paulo",
		FoundationYear: "2023",
		CurrentLord:    "",
	}

	cases := map[string]struct {
		input        entities.HouseRequest
		expectedData entities.House
		expectedErr  error
		prepareMock  func(mock *houses.MockIRepository)
	}{
		"Should return success": {
			input: req,
			prepareMock: func(mock *houses.MockIRepository) {
				mock.EXPECT().
					FindByID(gomock.Any(), req.ID).
					Times(1).
					Return(entities.House{
						ID:             "id_1",
						Name:           "house Patrick chagas",
						Region:         "sao paulo",
						FoundationYear: "2023",
						CurrentLord:    "",
					}, nil)

				mock.EXPECT().
					FindByName(gomock.Any(), req.Name).
					Times(1).
					Return(entities.House{}, errors.New("not found"))

				mock.EXPECT().
					Update(gomock.Any(), gomock.AssignableToTypeOf(&entities.House{})).
					Times(1).
					Return(nil)
			},
		},
		"Should return error find": {
			input:       req,
			expectedErr: ErrHouseNotFound,
			prepareMock: func(mock *houses.MockIRepository) {
				mock.EXPECT().
					FindByID(gomock.Any(), req.ID).
					Times(1).
					Return(entities.House{}, ErrHouseNotFound)
			},
		},
		"Should return error name already used": {
			input:       req,
			expectedErr: ErrNameUsed,
			prepareMock: func(mock *houses.MockIRepository) {
				mock.EXPECT().
					FindByID(gomock.Any(), req.ID).
					Times(1).
					Return(entities.House{
						ID:             "id_1",
						Name:           "house Patrick chagas",
						Region:         "sao paulo",
						FoundationYear: "2023",
						CurrentLord:    "",
					}, nil)

				mock.EXPECT().
					FindByName(gomock.Any(), req.Name).
					Times(1).
					Return(entities.House{}, nil)
			},
		},
		"Should return error update": {
			input:       req,
			expectedErr: errors.New("problem to query"),
			prepareMock: func(mock *houses.MockIRepository) {
				mock.EXPECT().
					FindByID(gomock.Any(), req.ID).
					Times(1).
					Return(entities.House{
						ID:             "id_1",
						Name:           "house Patrick chagas",
						Region:         "sao paulo",
						FoundationYear: "2023",
						CurrentLord:    "",
					}, nil)

				mock.EXPECT().
					FindByName(gomock.Any(), req.Name).
					Times(1).
					Return(entities.House{}, errors.New("not found"))

				mock.EXPECT().
					Update(gomock.Any(), gomock.AssignableToTypeOf(&entities.House{})).
					Times(1).
					Return(errors.New("problem to query"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			mock := houses.NewMockIRepository(ctrl)

			cs.prepareMock(mock)

			srv := New(&repositories.Container{
				Database: repositories.SqlContainer{House: mock}},
				logger.NewLogrusLogger(),
			)

			_, err := srv.Update(ctx, cs.input)

			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_Delete(t *testing.T) {
	id := "id_123"
	cases := map[string]struct {
		input       string
		expectedErr error
		prepareMock func(mock *houses.MockIRepository)
	}{
		"Should return success": {
			input: id,
			prepareMock: func(mock *houses.MockIRepository) {
				mock.EXPECT().
					FindByID(gomock.Any(), id).
					Times(1).
					Return(entities.House{}, nil)

				mock.EXPECT().
					Delete(gomock.Any(), id).
					Times(1).
					Return(nil)
			},
		},
		"Should return error find": {
			input:       id,
			expectedErr: ErrHouseNotFound,
			prepareMock: func(mock *houses.MockIRepository) {
				mock.EXPECT().
					FindByID(gomock.Any(), id).
					Times(1).
					Return(entities.House{}, ErrHouseNotFound)
			},
		},
		"Should return error delete": {
			input:       id,
			expectedErr: errors.New("problem to query"),
			prepareMock: func(mock *houses.MockIRepository) {
				mock.EXPECT().
					FindByID(gomock.Any(), id).
					Times(1).
					Return(entities.House{}, nil)

				mock.EXPECT().
					Delete(gomock.Any(), id).
					Times(1).
					Return(errors.New("problem to query"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			mock := houses.NewMockIRepository(ctrl)

			cs.prepareMock(mock)

			srv := New(&repositories.Container{
				Database: repositories.SqlContainer{House: mock}},
				logger.NewLogrusLogger(),
			)

			err := srv.Delete(ctx, cs.input)

			assert.Equal(t, cs.expectedErr, err)
		})
	}
}
