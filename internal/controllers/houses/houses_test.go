package houses

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
	"github.com/PatrickChagastavares/game-of-thrones/internal/services/houses"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/httpRouter"
	"github.com/PatrickChagastavares/game-of-thrones/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Create(t *testing.T) {
	endpoint := "/houses"
	idCreated := "id_123"
	cases := map[string]struct {
		inputBody    func() io.Reader
		expectedCode int
		expectedData func() string
		prepareMock  func(mock *houses.MockIService)
	}{
		"Should return success": {
			inputBody: func() io.Reader {
				data := entities.HouseRequest{
					Name:           "house Patrick",
					Region:         "sao paulo",
					FoundationYear: "2023",
					CurrentLord:    "",
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
			prepareMock: func(mock *houses.MockIService) {
				mock.EXPECT().
					Create(gomock.Any(), entities.HouseRequest{
						Name:           "house Patrick",
						Region:         "sao paulo",
						FoundationYear: "2023",
						CurrentLord:    "",
					}).
					Times(1).
					Return(idCreated, nil)
			},
		},
		"Should return error decode": {
			inputBody: func() io.Reader {
				bt := []byte(`{"name":123,"region":"sp","foundation_year":123}`)
				return bytes.NewReader(bt)
			},
			expectedCode: http.StatusBadRequest,
			expectedData: func() string {
				bt, _ := json.Marshal(entities.ErrDecode)
				return string(bt)
			},
			prepareMock: func(mock *houses.MockIService) {},
		},
		"Should return error validate": {
			inputBody: func() io.Reader {
				data := entities.HouseRequest{
					Name: "Pa",
				}
				bt, _ := json.Marshal(data)
				return bytes.NewReader(bt)
			},
			expectedCode: http.StatusBadRequest,
			expectedData: func() string {
				bt := []byte(`{"http_code":400,"message":"invalid_payload","detail":[{"field":"name","error":"min","value":"Pa"},{"field":"region","error":"required","value":""},{"field":"foundation_year","error":"required","value":""}]}`)
				return string(bt)
			},
			prepareMock: func(mock *houses.MockIService) {},
		},
		"Should return error service": {
			inputBody: func() io.Reader {
				data := entities.HouseRequest{
					Name:           "house Patrick",
					Region:         "sao paulo",
					FoundationYear: "2023",
					CurrentLord:    "",
				}
				bt, _ := json.Marshal(data)
				return bytes.NewReader(bt)
			},
			expectedCode: http.StatusInternalServerError,
			expectedData: func() string {
				return "\"failed to create house\""
			},
			prepareMock: func(mock *houses.MockIService) {
				mock.EXPECT().
					Create(gomock.Any(), entities.HouseRequest{
						Name:           "house Patrick",
						Region:         "sao paulo",
						FoundationYear: "2023",
						CurrentLord:    "",
					}).
					Times(1).
					Return("", errors.New("failed to create house"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// ============ TEST CONTROLLER ============
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			// ============ MOCKS ============
			mock := houses.NewMockIService(ctrl)
			cs.prepareMock(mock)

			// ============ START CONTROLLER ============
			ctr := New(
				&services.Container{House: mock},
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
	endpoint := "/houses"
	data := []entities.House{
		{ID: "id_1", Name: "House Algood", Region: "sao paulo", FoundationYear: "2023", CurrentLord: ""},
		{ID: "id_1", Name: "house Patrick Chagas", Region: "sao paulo", FoundationYear: "2023", CurrentLord: ""},
	}
	cases := map[string]struct {
		inputPath    string
		expectedCode int
		expectedData func() string
		prepareMock  func(mock *houses.MockIService)
	}{
		"Should return success with name": {
			inputPath:    "?name=House%20Algood",
			expectedCode: http.StatusOK,
			expectedData: func() string {
				bt, _ := json.Marshal([]entities.House{data[0]})
				return string(bt)
			},
			prepareMock: func(mock *houses.MockIService) {
				mock.EXPECT().
					Find(gomock.Any(), "House Algood").
					Times(1).
					Return([]entities.House{data[0]}, nil)
			},
		},
		"Should return success without name": {
			expectedCode: http.StatusOK,
			expectedData: func() string {
				bt, _ := json.Marshal(data)
				return string(bt)
			},
			prepareMock: func(mock *houses.MockIService) {
				mock.EXPECT().
					Find(gomock.Any(), "").
					Times(1).
					Return(data, nil)
			},
		},
		"Should return error service ": {
			expectedCode: http.StatusBadRequest,
			expectedData: func() string {
				resp := entities.NewHttpErr(http.StatusBadRequest, houses.ErrFind.Error(), nil)
				bt, _ := json.Marshal(resp)
				return string(bt)
			},
			prepareMock: func(mock *houses.MockIService) {
				mock.EXPECT().
					Find(gomock.Any(), "").
					Times(1).
					Return(nil, houses.ErrFind)
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// ============ TEST CONTROLLER ============
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			// ============ MOCKS ============
			mock := houses.NewMockIService(ctrl)
			cs.prepareMock(mock)

			// ============ START CONTROLLER ============
			ctr := New(
				&services.Container{House: mock},
				logger.NewLogrusLogger(),
			)

			// ============ START ROUTER ============
			router := httpRouter.NewGinRouter()
			router.Get(endpoint, ctr.Find)

			// ============ START MOCK REQUEST ============
			request := httptest.NewRequest(http.MethodGet, endpoint+cs.inputPath, nil).WithContext(ctx)
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
	endpoint := "/houses/"
	data := entities.House{
		ID:             "id_1",
		Name:           "Patrick",
		Region:         "sao paulo",
		FoundationYear: "2023",
		CurrentLord:    "",
	}
	cases := map[string]struct {
		inputPath    string
		expectedCode int
		expectedData func() string
		prepareMock  func(mock *houses.MockIService)
	}{
		"Should return success": {
			inputPath:    data.ID,
			expectedCode: http.StatusOK,
			expectedData: func() string {
				bt, _ := json.Marshal(data)
				return string(bt)
			},
			prepareMock: func(mock *houses.MockIService) {
				mock.EXPECT().
					FindByID(gomock.Any(), data.ID).
					Times(1).
					Return(data, nil)
			},
		},
		"Should return error service": {
			inputPath:    data.ID,
			expectedCode: http.StatusBadRequest,
			expectedData: func() string {
				resp := entities.NewHttpErr(http.StatusBadRequest, houses.ErrHouseNotFound.Error(), nil)
				bt, _ := json.Marshal(resp)
				return string(bt)
			},
			prepareMock: func(mock *houses.MockIService) {
				mock.EXPECT().
					FindByID(gomock.Any(), data.ID).
					Times(1).
					Return(entities.House{}, houses.ErrHouseNotFound)
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// ============ TEST CONTROLLER ============
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			// ============ MOCKS ============
			mock := houses.NewMockIService(ctrl)
			cs.prepareMock(mock)

			// ============ START CONTROLLER ============
			ctr := New(
				&services.Container{House: mock},
				logger.NewLogrusLogger(),
			)

			// ============ START ROUTER ============
			router := httpRouter.NewGinRouter()
			router.Get(endpoint+":id", ctr.FindByID)

			// ============ START MOCK REQUEST ============
			request := httptest.NewRequest(http.MethodGet, endpoint+cs.inputPath, nil).WithContext(ctx)
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
	endpoint := "/houses/"

	resp := entities.House{
		ID:             "id_1",
		Name:           "house Chagas",
		Region:         "sao paulo",
		FoundationYear: "2023",
		CurrentLord:    "",
	}
	cases := map[string]struct {
		inputPath    string
		inputBody    func() io.Reader
		expectedCode int
		expectedData func() string
		prepareMock  func(mock *houses.MockIService)
	}{
		"Should return success": {
			inputPath: resp.ID,
			inputBody: func() io.Reader {
				data := entities.HouseRequest{
					Name:           "house Chagas",
					Region:         "sao paulo",
					FoundationYear: "2023",
					CurrentLord:    "",
				}
				bt, _ := json.Marshal(data)
				return bytes.NewReader(bt)
			},
			expectedCode: http.StatusOK,
			expectedData: func() string {
				bt, _ := json.Marshal(resp)
				return string(bt)
			},
			prepareMock: func(mock *houses.MockIService) {
				mock.EXPECT().
					Update(gomock.Any(), entities.HouseRequest{
						ID:             resp.ID,
						Name:           "house Chagas",
						Region:         "sao paulo",
						FoundationYear: "2023",
						CurrentLord:    "",
					}).
					Times(1).
					Return(resp, nil)
			},
		},
		"Should return error decode": {
			inputPath: resp.ID,
			inputBody: func() io.Reader {
				bt := []byte(`{"name":123,"region":"sp","foundation_year":123}`)
				return bytes.NewReader(bt)
			},
			expectedCode: http.StatusBadRequest,
			expectedData: func() string {
				bt, _ := json.Marshal(entities.ErrDecode)
				return string(bt)
			},
			prepareMock: func(mock *houses.MockIService) {},
		},
		"Should return error validate": {
			inputPath: resp.ID,
			inputBody: func() io.Reader {
				data := entities.HouseRequest{
					Name: "Pa",
				}
				bt, _ := json.Marshal(data)
				return bytes.NewReader(bt)
			},
			expectedCode: http.StatusBadRequest,
			expectedData: func() string {
				bt := []byte(`{"http_code":400,"message":"invalid_payload","detail":[{"field":"name","error":"min","value":"Pa"},{"field":"region","error":"required","value":""},{"field":"foundation_year","error":"required","value":""}]}`)
				return string(bt)
			},
			prepareMock: func(mock *houses.MockIService) {},
		},
		"Should return error service": {
			inputPath: resp.ID,
			inputBody: func() io.Reader {
				data := entities.HouseRequest{
					Name:           "house Chagas",
					Region:         "sao paulo",
					FoundationYear: "2023",
					CurrentLord:    "",
				}
				bt, _ := json.Marshal(data)
				return bytes.NewReader(bt)
			},
			expectedCode: http.StatusInternalServerError,
			expectedData: func() string {
				return "\"failed to update house\""
			},
			prepareMock: func(mock *houses.MockIService) {
				mock.EXPECT().
					Update(gomock.Any(), entities.HouseRequest{
						ID:             resp.ID,
						Name:           "house Chagas",
						Region:         "sao paulo",
						FoundationYear: "2023",
						CurrentLord:    "",
					}).
					Times(1).
					Return(entities.House{}, errors.New("failed to update house"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// ============ TEST CONTROLLER ============
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			// ============ MOCKS ============
			mock := houses.NewMockIService(ctrl)
			cs.prepareMock(mock)

			// ============ START CONTROLLER ============
			ctr := New(
				&services.Container{House: mock},
				logger.NewLogrusLogger(),
			)

			// ============ START ROUTER ============
			router := httpRouter.NewGinRouter()
			router.Put(endpoint+":id", ctr.Update)

			// ============ START MOCK REQUEST ============
			request := httptest.NewRequest(http.MethodPut, endpoint+cs.inputPath, cs.inputBody()).WithContext(ctx)
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
	endpoint := "/houses/"
	cases := map[string]struct {
		paramInput   string
		expectedCode int
		expectedData func() string
		prepareMock  func(mock *houses.MockIService)
	}{
		"Should return success": {
			paramInput:   "33c55a43-f163-4a67-9f6c-75161410f376",
			expectedCode: http.StatusNoContent,
			expectedData: func() string {
				return ""
			},
			prepareMock: func(mock *houses.MockIService) {
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
				return "\"failed to delete house\""
			},
			prepareMock: func(mock *houses.MockIService) {
				mock.EXPECT().
					Delete(gomock.Any(), "33c55a43-f163-4a67-9f6c-75161410f376").
					Times(1).
					Return(errors.New("failed to delete house"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// ============ TEST CONTROLLER ============
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			// ============ MOCKS ============
			mock := houses.NewMockIService(ctrl)
			cs.prepareMock(mock)

			// ============ START CONTROLLER ============
			ctr := New(
				&services.Container{House: mock},
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
