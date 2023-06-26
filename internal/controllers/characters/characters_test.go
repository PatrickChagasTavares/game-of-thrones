package characters

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PatrickChagastavares/game-of-thrones/internal/entities"
	"github.com/PatrickChagastavares/game-of-thrones/internal/services"
	"github.com/PatrickChagastavares/game-of-thrones/internal/services/characters"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/httpRouter"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func Test_Create(t *testing.T) {
	endpoint := "/characters"
	idCreated := "id_123"
	cases := map[string]struct {
		inputBody    func() io.Reader
		expectedCode int
		expectedData func() string
		prepareMock  func(mock *characters.MockIService)
	}{
		"Should return success": {
			inputBody: func() io.Reader {
				data := entities.CharacterRequest{
					Name:     "Patrick Chagas",
					TVSeries: pq.StringArray{"session 1", "session 2"},
				}
				bt, _ := json.Marshal(data)
				return bytes.NewReader(bt)
			},
			expectedCode: http.StatusCreated,
			expectedData: func() string {
				data := map[string]any{
					"id": idCreated,
				}
				bt, _ := json.Marshal(data)
				return string(bt)
			},
			prepareMock: func(mock *characters.MockIService) {
				mock.EXPECT().
					Create(gomock.Any(), entities.CharacterRequest{
						Name:     "Patrick Chagas",
						TVSeries: pq.StringArray{"session 1", "session 2"},
					}).
					Times(1).
					Return(idCreated, nil)
			},
		},
		"Should return error decode": {
			inputBody: func() io.Reader {
				bt := []byte(`{"name":123,"tv_series":"Patrick"}`)
				return bytes.NewReader(bt)
			},
			expectedCode: http.StatusBadRequest,
			expectedData: func() string {
				bt, _ := json.Marshal(entities.ErrDecode)
				return string(bt)
			},
			prepareMock: func(mock *characters.MockIService) {},
		},
		"Should return error validate": {
			inputBody: func() io.Reader {
				data := entities.CharacterRequest{
					Name: "Pa",
				}
				bt, _ := json.Marshal(data)
				return bytes.NewReader(bt)
			},
			expectedCode: http.StatusBadRequest,
			expectedData: func() string {
				bt := []byte(`{"http_code":400,"message":"invalid_payload","detail":[{"field":"name","error":"min","value":"Pa"},{"field":"tv_series","error":"required","value":null}]}`)
				return string(bt)
			},
			prepareMock: func(mock *characters.MockIService) {},
		},
		"Should return error service": {
			inputBody: func() io.Reader {
				data := entities.CharacterRequest{
					Name:     "Patrick Chagas",
					TVSeries: pq.StringArray{"session 1", "session 2"},
				}
				bt, _ := json.Marshal(data)
				return bytes.NewReader(bt)
			},
			expectedCode: http.StatusInternalServerError,
			expectedData: func() string {
				return "\"failed to create character\""
			},
			prepareMock: func(mock *characters.MockIService) {
				mock.EXPECT().
					Create(gomock.Any(), entities.CharacterRequest{
						Name:     "Patrick Chagas",
						TVSeries: pq.StringArray{"session 1", "session 2"},
					}).
					Times(1).
					Return("", errors.New("failed to create character"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// ============ TEST CONTROLLER ============
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			// ============ MOCKS ============
			mock := characters.NewMockIService(ctrl)
			cs.prepareMock(mock)

			// ============ START CONTROLLER ============
			ctr := New(
				&services.Container{Character: mock},
				logger.NewLogrusLogger(),
			)

			// ============ START ROUTER ============
			router := httpRouter.NewGinRouter()
			router.Post(endpoint, ctr.Create)

			// ============ START MOCK REQUEST ============
			request := httptest.NewRequest(http.MethodPost, endpoint, cs.inputBody()).WithContext(ctx)
			writer := httptest.NewRecorder()
			request.Header.Set("Content-Type", "application/json")

			// ============ START SERVER HTTP ============
			router.ServeHTTP(writer, request)

			// ============ START VALIDATION ============
			responseData, _ := ioutil.ReadAll(writer.Body)
			assert.Equal(t, cs.expectedData(), string(responseData))
			assert.Equal(t, cs.expectedCode, writer.Code)
		})
	}
}

func Test_Find(t *testing.T) {
	endpoint := "/characters"
	data := []entities.Character{
		{ID: "id_1", Name: "character Patrick", TVSeries: pq.StringArray{"session 1", "session 2"}},
		{ID: "id_2", Name: "character Patrick", TVSeries: pq.StringArray{"session 1", "session 2"}},
	}
	cases := map[string]struct {
		expectedCode int
		expectedData func() string
		prepareMock  func(mock *characters.MockIService)
	}{
		"Should return success": {
			expectedCode: http.StatusOK,
			expectedData: func() string {
				bt, _ := json.Marshal(data)
				return string(bt)
			},
			prepareMock: func(mock *characters.MockIService) {
				mock.EXPECT().
					Find(gomock.Any()).
					Times(1).
					Return(data, nil)
			},
		},
		"Should return error service": {
			expectedCode: http.StatusBadRequest,
			expectedData: func() string {
				resp := entities.NewHttpErr(http.StatusBadRequest, characters.ErrFind.Error(), nil)
				bt, _ := json.Marshal(resp)
				return string(bt)
			},
			prepareMock: func(mock *characters.MockIService) {
				mock.EXPECT().
					Find(gomock.Any()).
					Times(1).
					Return(nil, characters.ErrFind)
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// ============ TEST CONTROLLER ============
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			// ============ MOCKS ============
			mock := characters.NewMockIService(ctrl)
			cs.prepareMock(mock)

			// ============ START CONTROLLER ============
			ctr := New(
				&services.Container{Character: mock},
				logger.NewLogrusLogger(),
			)

			// ============ START ROUTER ============
			router := httpRouter.NewGinRouter()
			router.Get(endpoint, ctr.Find)

			// ============ START MOCK REQUEST ============
			request := httptest.NewRequest(http.MethodGet, endpoint, nil).WithContext(ctx)
			writer := httptest.NewRecorder()
			request.Header.Set("Content-Type", "application/json")

			// ============ START SERVER HTTP ============
			router.ServeHTTP(writer, request)

			// ============ START VALIDATION ============
			responseData, _ := ioutil.ReadAll(writer.Body)
			assert.Equal(t, cs.expectedData(), string(responseData))
			assert.Equal(t, cs.expectedCode, writer.Code)
		})
	}
}

func Test_FindByID(t *testing.T) {
	endpoint := "/characters/"
	data := entities.Character{
		ID:       "id_123",
		Name:     "character Patrick",
		TVSeries: pq.StringArray{"session 1", "session 2"},
	}
	cases := map[string]struct {
		expectedCode int
		expectedData func() string
		prepareMock  func(mock *characters.MockIService)
	}{
		"Should return success": {
			expectedCode: http.StatusOK,
			expectedData: func() string {
				bt, _ := json.Marshal(data)
				return string(bt)
			},
			prepareMock: func(mock *characters.MockIService) {
				mock.EXPECT().
					FindByID(gomock.Any(), data.ID).
					Times(1).
					Return(data, nil)
			},
		},
		"Should return error service": {
			expectedCode: http.StatusBadRequest,
			expectedData: func() string {
				resp := entities.NewHttpErr(http.StatusBadRequest, characters.ErrCharacterNotFound.Error(), nil)
				bt, _ := json.Marshal(resp)
				return string(bt)
			},
			prepareMock: func(mock *characters.MockIService) {
				mock.EXPECT().
					FindByID(gomock.Any(), data.ID).
					Times(1).
					Return(entities.Character{}, characters.ErrCharacterNotFound)
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// ============ TEST CONTROLLER ============
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			// ============ MOCKS ============
			mock := characters.NewMockIService(ctrl)
			cs.prepareMock(mock)

			// ============ START CONTROLLER ============
			ctr := New(
				&services.Container{Character: mock},
				logger.NewLogrusLogger(),
			)

			// ============ START ROUTER ============
			router := httpRouter.NewGinRouter()
			router.Get(endpoint+":id", ctr.FindByID)

			// ============ START MOCK REQUEST ============
			request := httptest.NewRequest(http.MethodGet, endpoint+data.ID, nil).WithContext(ctx)
			writer := httptest.NewRecorder()
			request.Header.Set("Content-Type", "application/json")

			// ============ START SERVER HTTP ============
			router.ServeHTTP(writer, request)

			// ============ START VALIDATION ============
			responseData, _ := ioutil.ReadAll(writer.Body)
			assert.Equal(t, cs.expectedData(), string(responseData))
			assert.Equal(t, cs.expectedCode, writer.Code)
		})
	}
}

func Test_Update(t *testing.T) {
	endpoint := "/characters/"

	resp := entities.Character{
		ID:       "id_123",
		Name:     "House Patrick",
		TVSeries: pq.StringArray{"session 1"},
	}
	cases := map[string]struct {
		inputBody    func() io.Reader
		expectedCode int
		expectedData func() string
		prepareMock  func(mock *characters.MockIService)
	}{
		"Should return success": {
			inputBody: func() io.Reader {
				data := entities.CharacterRequest{
					Name:     "house Chagas",
					TVSeries: pq.StringArray{"session 1"},
				}
				bt, _ := json.Marshal(data)
				return bytes.NewReader(bt)
			},
			expectedCode: http.StatusOK,
			expectedData: func() string {
				bt, _ := json.Marshal(resp)
				return string(bt)
			},
			prepareMock: func(mock *characters.MockIService) {
				mock.EXPECT().
					Update(gomock.Any(), entities.CharacterRequest{
						ID:       resp.ID,
						Name:     "house Chagas",
						TVSeries: pq.StringArray{"session 1"},
					}).
					Times(1).
					Return(resp, nil)
			},
		},
		"Should return error decode": {
			inputBody: func() io.Reader {
				bt := []byte(`{"name":123,"tv_series":"Patrick"}`)
				return bytes.NewReader(bt)
			},
			expectedCode: http.StatusBadRequest,
			expectedData: func() string {
				bt, _ := json.Marshal(entities.ErrDecode)
				return string(bt)
			},
			prepareMock: func(mock *characters.MockIService) {},
		},
		"Should return error validate": {
			inputBody: func() io.Reader {
				data := entities.CharacterRequest{
					Name: "Pa",
				}
				bt, _ := json.Marshal(data)
				return bytes.NewReader(bt)
			},
			expectedCode: http.StatusBadRequest,
			expectedData: func() string {
				bt := []byte(`{"http_code":400,"message":"invalid_payload","detail":[{"field":"name","error":"min","value":"Pa"},{"field":"tv_series","error":"required","value":null}]}`)
				return string(bt)
			},
			prepareMock: func(mock *characters.MockIService) {},
		},
		"Should return error service": {
			inputBody: func() io.Reader {
				data := entities.CharacterRequest{
					Name:     "house Chagas",
					TVSeries: pq.StringArray{"session 1"},
				}
				bt, _ := json.Marshal(data)
				return bytes.NewReader(bt)
			},
			expectedCode: http.StatusInternalServerError,
			expectedData: func() string {
				return "\"failed to create character\""
			},
			prepareMock: func(mock *characters.MockIService) {
				mock.EXPECT().
					Update(gomock.Any(), entities.CharacterRequest{
						ID:       resp.ID,
						Name:     "house Chagas",
						TVSeries: pq.StringArray{"session 1"},
					}).
					Times(1).
					Return(entities.Character{}, errors.New("failed to create character"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// ============ TEST CONTROLLER ============
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			// ============ MOCKS ============
			mock := characters.NewMockIService(ctrl)
			cs.prepareMock(mock)

			// ============ START CONTROLLER ============
			ctr := New(
				&services.Container{Character: mock},
				logger.NewLogrusLogger(),
			)

			// ============ START ROUTER ============
			router := httpRouter.NewGinRouter()
			router.Put(endpoint+":id", ctr.Update)

			// ============ START MOCK REQUEST ============
			request := httptest.NewRequest(http.MethodPut, endpoint+resp.ID, cs.inputBody()).WithContext(ctx)
			writer := httptest.NewRecorder()
			request.Header.Set("Content-Type", "application/json")

			// ============ START SERVER HTTP ============
			router.ServeHTTP(writer, request)

			// ============ START VALIDATION ============
			responseData, _ := ioutil.ReadAll(writer.Body)
			assert.Equal(t, cs.expectedData(), string(responseData))
			assert.Equal(t, cs.expectedCode, writer.Code)
		})
	}
}

func Test_Delete(t *testing.T) {
	endpoint := "/characters/"
	cases := map[string]struct {
		paramInput   string
		expectedCode int
		expectedData func() string
		prepareMock  func(mock *characters.MockIService)
	}{
		"Should return success": {
			paramInput:   "33c55a43-f163-4a67-9f6c-75161410f376",
			expectedCode: http.StatusNoContent,
			expectedData: func() string {
				return ""
			},
			prepareMock: func(mock *characters.MockIService) {
				mock.EXPECT().
					Delete(gomock.Any(), "33c55a43-f163-4a67-9f6c-75161410f376").
					Times(1).
					Return(nil)
			},
		},
		"Should return error service": {
			paramInput:   "33c55a43-f163-4a67-9f6c-75161410f376",
			expectedCode: http.StatusInternalServerError,
			expectedData: func() string {
				return "\"failed to delete character\""
			},
			prepareMock: func(mock *characters.MockIService) {
				mock.EXPECT().
					Delete(gomock.Any(), "33c55a43-f163-4a67-9f6c-75161410f376").
					Times(1).
					Return(errors.New("failed to delete character"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// ============ TEST CONTROLLER ============
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			// ============ MOCKS ============
			mock := characters.NewMockIService(ctrl)
			cs.prepareMock(mock)

			// ============ START CONTROLLER ============
			ctr := New(
				&services.Container{Character: mock},
				logger.NewLogrusLogger(),
			)

			// ============ START ROUTER ============
			router := httpRouter.NewGinRouter()
			router.Delete(endpoint+":id", ctr.Delete)

			// ============ START MOCK REQUEST ============
			request := httptest.NewRequest(http.MethodDelete, endpoint+cs.paramInput, nil).WithContext(ctx)
			writer := httptest.NewRecorder()
			request.Header.Set("Content-Type", "application/json")

			// ============ START SERVER HTTP ============
			router.ServeHTTP(writer, request)

			// ============ START VALIDATION ============
			responseData, _ := ioutil.ReadAll(writer.Body)
			assert.Equal(t, cs.expectedData(), string(responseData))
			assert.Equal(t, cs.expectedCode, writer.Code)
		})
	}
}
