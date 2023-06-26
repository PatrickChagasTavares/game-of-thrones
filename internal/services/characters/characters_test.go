package characters

import (
	"context"
	"errors"
	"testing"

	"github.com/PatrickChagastavares/game-of-thrones/internal/entities"
	"github.com/PatrickChagastavares/game-of-thrones/internal/repositories"
	"github.com/PatrickChagastavares/game-of-thrones/internal/repositories/database/characters"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var (
	characterID   string
	characterTime string
)

func Test_Create(t *testing.T) {
	data := entities.CharacterRequest{
		Name:     "character Patrick",
		TVSeries: pq.StringArray{"session 1", "session 2"},
	}

	cases := map[string]struct {
		input entities.CharacterRequest

		expectedErr error
		prepareMock func(mock *characters.MockIRepository)
	}{
		"Should return success": {
			input: data,
			prepareMock: func(mock *characters.MockIRepository) {
				mock.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(entities.CharacterRequest{})).
					Times(1).
					Return(nil)
			},
		},
		"Should return error": {
			input:       data,
			expectedErr: errors.New("problem to create user"),
			prepareMock: func(mock *characters.MockIRepository) {
				mock.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(entities.CharacterRequest{})).
					Times(1).
					Return(errors.New("problem to create user"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			mock := characters.NewMockIRepository(ctrl)

			cs.prepareMock(mock)

			srv := New(&repositories.Container{Database: repositories.SqlContainer{Character: mock}}, logger.NewLogrusLogger())

			_, err := srv.Create(ctx, cs.input)

			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func Test_Find(t *testing.T) {
	data := []entities.Character{
		{ID: "id_1", Name: "character Patrick", TVSeries: pq.StringArray{"session 1", "session 2"}},
		{ID: "id_2", Name: "character Patrick", TVSeries: pq.StringArray{"session 1", "session 2"}},
	}

	cases := map[string]struct {
		expectedData []entities.Character
		expectedErr  error
		prepareMock  func(mock *characters.MockIRepository)
	}{
		"Should return success": {
			expectedData: data,
			prepareMock: func(mock *characters.MockIRepository) {
				mock.EXPECT().
					Find(gomock.Any()).
					Times(1).
					Return(data, nil)
			},
		},
		"Should return error": {
			expectedErr: ErrFind,
			prepareMock: func(mock *characters.MockIRepository) {
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

			mock := characters.NewMockIRepository(ctrl)

			cs.prepareMock(mock)

			srv := New(&repositories.Container{Database: repositories.SqlContainer{Character: mock}}, logger.NewLogrusLogger())

			data, err := srv.Find(ctx)

			assert.Equal(t, cs.expectedErr, err)
			assert.Equal(t, cs.expectedData, data)
		})
	}
}

func Test_FindByID(t *testing.T) {
	data := entities.Character{
		ID:       "id_1",
		Name:     "character Patrick",
		TVSeries: pq.StringArray{"session 1", "session 2"},
	}

	cases := map[string]struct {
		input        string
		expectedData entities.Character
		expectedErr  error
		prepareMock  func(mock *characters.MockIRepository)
	}{
		"Should return success": {
			input:        data.ID,
			expectedData: data,
			prepareMock: func(mock *characters.MockIRepository) {
				mock.EXPECT().
					FindByID(gomock.Any(), data.ID).
					Times(1).
					Return(data, nil)
			},
		},
		"Should return error": {
			input:       data.ID,
			expectedErr: ErrCharacterNotFound,
			prepareMock: func(mock *characters.MockIRepository) {
				mock.EXPECT().
					FindByID(gomock.Any(), data.ID).
					Times(1).
					Return(entities.Character{}, errors.New("problem to query"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			mock := characters.NewMockIRepository(ctrl)

			cs.prepareMock(mock)

			srv := New(&repositories.Container{Database: repositories.SqlContainer{Character: mock}}, logger.NewLogrusLogger())

			data, err := srv.FindByID(ctx, cs.input)

			assert.Equal(t, cs.expectedErr, err)
			assert.Equal(t, cs.expectedData, data)
		})
	}
}

func Test_Update(t *testing.T) {
	req := entities.CharacterRequest{
		ID:       "id_1",
		Name:     "character 13",
		TVSeries: pq.StringArray{"session 1", "session 2"},
	}

	cases := map[string]struct {
		input        entities.CharacterRequest
		expectedData entities.Character
		expectedErr  error
		prepareMock  func(mock *characters.MockIRepository)
	}{
		"Should return success": {
			input: req,
			prepareMock: func(mock *characters.MockIRepository) {
				mock.EXPECT().
					FindByID(gomock.Any(), req.ID).
					Times(1).
					Return(entities.Character{
						ID:       "id_1",
						Name:     "character Patrick",
						TVSeries: pq.StringArray{"session 1", "session 2"},
					}, nil)

				mock.EXPECT().
					Update(gomock.Any(), gomock.AssignableToTypeOf(&entities.Character{})).
					Times(1).
					Return(nil)
			},
		},
		"Should return error find": {
			input:       req,
			expectedErr: ErrCharacterNotFound,
			prepareMock: func(mock *characters.MockIRepository) {
				mock.EXPECT().
					FindByID(gomock.Any(), req.ID).
					Times(1).
					Return(entities.Character{}, ErrCharacterNotFound)
			},
		},
		"Should return error update": {
			input:       req,
			expectedErr: errors.New("problem to query"),
			prepareMock: func(mock *characters.MockIRepository) {
				mock.EXPECT().
					FindByID(gomock.Any(), req.ID).
					Times(1).
					Return(entities.Character{
						ID:       "id_1",
						Name:     "character Patrick",
						TVSeries: pq.StringArray{"session 1", "session 2"},
					}, nil)

				mock.EXPECT().
					Update(gomock.Any(), gomock.AssignableToTypeOf(&entities.Character{})).
					Times(1).
					Return(errors.New("problem to query"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			mock := characters.NewMockIRepository(ctrl)

			cs.prepareMock(mock)

			srv := New(&repositories.Container{Database: repositories.SqlContainer{Character: mock}}, logger.NewLogrusLogger())

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
		prepareMock func(mock *characters.MockIRepository)
	}{
		"Should return success": {
			input: id,
			prepareMock: func(mock *characters.MockIRepository) {
				mock.EXPECT().
					FindByID(gomock.Any(), id).
					Times(1).
					Return(entities.Character{}, nil)

				mock.EXPECT().
					Delete(gomock.Any(), id).
					Times(1).
					Return(nil)
			},
		},
		"Should return error find": {
			input:       id,
			expectedErr: ErrCharacterNotFound,
			prepareMock: func(mock *characters.MockIRepository) {
				mock.EXPECT().
					FindByID(gomock.Any(), id).
					Times(1).
					Return(entities.Character{}, ErrCharacterNotFound)
			},
		},
		"Should return error delete": {
			input:       id,
			expectedErr: errors.New("problem to query"),
			prepareMock: func(mock *characters.MockIRepository) {
				mock.EXPECT().
					FindByID(gomock.Any(), id).
					Times(1).
					Return(entities.Character{}, nil)

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

			mock := characters.NewMockIRepository(ctrl)

			cs.prepareMock(mock)

			srv := New(&repositories.Container{Database: repositories.SqlContainer{Character: mock}}, logger.NewLogrusLogger())

			err := srv.Delete(ctx, cs.input)

			assert.Equal(t, cs.expectedErr, err)
		})
	}
}
