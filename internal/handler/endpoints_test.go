package handler

// import (
// 	"database/sql"
// 	"fmt"
// 	"testing"
// 	"time"

// 	"net/http"
// 	"net/http/httptest"

// 	"github.com/google/uuid"
// 	"github.com/labstack/echo/v4"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/vldcreation/movie-fest/generated"
// 	"github.com/vldcreation/movie-fest/pkg/builderx/json/json_iter"
// 	"github.com/vldcreation/movie-fest/pkg/builderx/url"
// 	"github.com/vldcreation/movie-fest/pkg/types"
// 	"github.com/vldcreation/movie-fest/repository"
// 	"go.uber.org/mock/gomock"
// )

// func TestGetHello(t *testing.T) {
// 	s := NewServer(NewServerOptions{Repository: repository.NewMockRepositoryInterface(nil)})

// 	t.Run("GetHello", func(t *testing.T) {
// 		// create a params object
// 		params := generated.GetHelloParams{
// 			Id: 1,
// 		}

// 		expectedRes := generated.HelloResponse{
// 			Message: "Hello User 1",
// 		}

// 		// Create a new Echo instance
// 		e := echo.New()

// 		// Create a new request with the desired parameters
// 		path := fmt.Sprintf("/hello?%s", url.BuildQuery(params).Encode())
// 		req := httptest.NewRequest(http.MethodGet, path, nil)
// 		rec := httptest.NewRecorder()
// 		ctx := e.NewContext(req, rec)

// 		// Call the handler function
// 		err := s.GetHello(ctx, params)

// 		// Assertions
// 		assert.NoError(t, err)                   // Check that there was no error
// 		assert.Equal(t, http.StatusOK, rec.Code) // Check that the status code is 200

// 		// Check the response body
// 		expectedResponse, err := json_iter.Marshal(expectedRes)
// 		assert.NoError(t, err)
// 		assert.JSONEq(t, string(expectedResponse), rec.Body.String()) // Check the response body
// 	})
// }

// func TestSaveEstate(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockRepository := repository.NewMockRepositoryInterface(ctrl)
// 	mockRepository.EXPECT().SaveEstate(gomock.Any(), gomock.Any()).Return(repository.SaveEstateOutput{
// 		Id: uuid.New().String(),
// 	}, nil)
// 	s := NewServer(NewServerOptions{Repository: mockRepository})

// 	t.Run("SaveEstate 400 Bad Request", func(t *testing.T) {
// 		// Create a new Echo instance
// 		e := echo.New()

// 		// Create a new request with the desired parameters
// 		req := httptest.NewRequest(http.MethodPost, "/estate", nil)
// 		rec := httptest.NewRecorder()
// 		ctx := e.NewContext(req, rec)

// 		// Call the handler function
// 		err := s.SaveEstate(ctx)

// 		// Assertions
// 		assert.NoError(t, err)                           // Check that there was no error
// 		assert.Equal(t, http.StatusBadRequest, rec.Code) // Check that the status code is 400
// 	})

// 	t.Run("SaveEstate 201 Created", func(t *testing.T) {
// 		// create a params object
// 		params := generated.SaveEstateJSONBody{
// 			Length: types.IntPtr(10),
// 			Width:  types.IntPtr(20),
// 		}

// 		// Create a new Echo instance
// 		e := echo.New()

// 		// Create body request
// 		data, err := json_iter.Marshal(params)
// 		assert.NoError(t, err)

// 		// Create a new request with the desired parameters
// 		req := httptest.NewRequest(http.MethodPost, "/estate", json_iter.GenerateReader(data))
// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		ctx := e.NewContext(req, rec)

// 		// Call the handler function
// 		err = s.SaveEstate(ctx)
// 		assert.NoError(t, err)

// 		res := repository.SaveEstateOutput{}
// 		err = json_iter.Unmarshal(rec.Body.Bytes(), &res)
// 		assert.NoError(t, err)

// 		// Assertions
// 		assert.NoError(t, err)                        // Check that there was no error
// 		assert.Equal(t, http.StatusCreated, rec.Code) // Check that the status code is 201

// 		assert.NotEmpty(t, res.Id) // Check that the response body is not empty
// 		estateID, err := uuid.Parse(res.Id)
// 		assert.NoError(t, err)                          // Check that there was no error
// 		assert.NotEqual(t, estateID, uuid.Nil.String()) // Check that the response body is not empty
// 	})
// }

// func TestAddTreeToEstate(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	t.Run("AddTreeToEstate 400 Bad Request", func(t *testing.T) {
// 		// Create a new Echo instance
// 		e := echo.New()
// 		mockRepository := repository.NewMockRepositoryInterface(ctrl)
// 		s := NewServer(NewServerOptions{Repository: mockRepository})

// 		testID := uuid.Nil // force to get invalid request

// 		// Create body request
// 		params := generated.AddTreeToEstateJSONBody{
// 			X:      types.IntPtr(20),
// 			Y:      types.IntPtr(10),
// 			Height: types.IntPtr(30),
// 		}
// 		data, err := json_iter.Marshal(params)
// 		assert.NoError(t, err)

// 		// mocking repoisitory layer

// 		// Create a new request with the desired parameters
// 		req := httptest.NewRequest(http.MethodPost, "/estate/"+testID.String()+"/tree", json_iter.GenerateReader(data))
// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		ctx := e.NewContext(req, rec)

// 		// Call the handler function
// 		err = s.AddTreeToEstate(ctx, testID)

// 		// Assertions
// 		assert.NoError(t, err)                           // Check that there was no error
// 		assert.Equal(t, http.StatusBadRequest, rec.Code) // Check that the status code is 404
// 	})

// 	t.Run("AddTreeToEstate 404 Not Found", func(t *testing.T) {
// 		// Create a new Echo instance
// 		e := echo.New()

// 		mockRepository := repository.NewMockRepositoryInterface(ctrl)
// 		s := NewServer(NewServerOptions{Repository: mockRepository})

// 		testID := uuid.New()

// 		// Create body request
// 		params := generated.AddTreeToEstateJSONBody{
// 			X:      types.IntPtr(20),
// 			Y:      types.IntPtr(10),
// 			Height: types.IntPtr(30),
// 		}
// 		data, err := json_iter.Marshal(params)
// 		assert.NoError(t, err)

// 		// mocking repoisitory layer
// 		mockRepository.EXPECT().GetEstate(gomock.Any(), gomock.Any()).Return(repository.GetEstateOutput{}, sql.ErrNoRows)

// 		// Create a new request with the desired parameters
// 		req := httptest.NewRequest(http.MethodPost, "/estate/"+testID.String()+"/tree", json_iter.GenerateReader(data))
// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		ctx := e.NewContext(req, rec)

// 		// Call the handler function
// 		err = s.AddTreeToEstate(ctx, testID)

// 		// Assertions
// 		assert.NoError(t, err)                         // Check that there was no error
// 		assert.Equal(t, http.StatusNotFound, rec.Code) // Check that the status code is 404
// 	})

// 	t.Run("AddTreeToEstate 201 Created", func(t *testing.T) {
// 		// Create a new Echo instance
// 		e := echo.New()
// 		mockRepository := repository.NewMockRepositoryInterface(ctrl)
// 		s := NewServer(NewServerOptions{Repository: mockRepository})

// 		testID := uuid.New()

// 		mockEstate := repository.GetEstateOutput{
// 			Id: testID.String(),
// 		}

// 		// mocking repoisitory layer
// 		// mocking creating a new instance of estate
// 		mockRepository.EXPECT().GetEstate(gomock.Any(), gomock.Any()).Return(mockEstate, nil)
// 		// mocking add tree to estate
// 		mockRepository.EXPECT().AddTreeToEstate(gomock.Any(), gomock.Any()).Return(repository.AddTreeToEstateOutput{
// 			Id: uuid.New().String(),
// 		}, nil)

// 		// create a params object
// 		params := generated.AddTreeToEstateJSONBody{
// 			X:      types.IntPtr(20),
// 			Y:      types.IntPtr(10),
// 			Height: types.IntPtr(30),
// 		}
// 		data, err := json_iter.Marshal(params)
// 		assert.NoError(t, err)

// 		// Create a new request with the desired parameters
// 		req := httptest.NewRequest(http.MethodPost, "/estate/"+testID.String()+"/tree", json_iter.GenerateReader(data))
// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		ctx := e.NewContext(req, rec)

// 		// Call the handler function
// 		err = s.AddTreeToEstate(ctx, testID)
// 		assert.NoError(t, err)

// 		res := repository.SaveEstateOutput{}
// 		err = json_iter.Unmarshal(rec.Body.Bytes(), &res)
// 		assert.NoError(t, err)

// 		// Assertions
// 		assert.NoError(t, err)                        // Check that there was no error
// 		assert.Equal(t, http.StatusCreated, rec.Code) // Check that the status code is 201

// 		assert.NotEmpty(t, res.Id) // Check that the response body is not empty
// 		estateID, err := uuid.Parse(res.Id)
// 		assert.NoError(t, err)                          // Check that there was no error
// 		assert.NotEqual(t, estateID, uuid.Nil.String()) // Check that the response body is not empty
// 	})
// }

// func TestGetEstateStats(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	t.Run("GetEstateStats 400 Bad Request", func(t *testing.T) {
// 		// Create a new Echo instance
// 		e := echo.New()
// 		mockRepository := repository.NewMockRepositoryInterface(ctrl)
// 		s := NewServer(NewServerOptions{Repository: mockRepository})

// 		testID := uuid.Nil // force to get invalid request

// 		// mocking repoisitory layer

// 		// Create a new request with the desired parameters
// 		req := httptest.NewRequest(http.MethodGet, "/estate/"+testID.String()+"/stats", nil)
// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		ctx := e.NewContext(req, rec)

// 		// Call the handler function
// 		err := s.GetEstateStats(ctx, testID)

// 		// Assertions
// 		assert.NoError(t, err)                           // Check that there was no error
// 		assert.Equal(t, http.StatusBadRequest, rec.Code) // Check that the status code is 404
// 	})

// 	t.Run("GetEstateStats 404 Not Found", func(t *testing.T) {
// 		// Create a new Echo instance
// 		e := echo.New()

// 		mockRepository := repository.NewMockRepositoryInterface(ctrl)
// 		s := NewServer(NewServerOptions{Repository: mockRepository})

// 		testID := uuid.New()

// 		// mocking repoisitory layer
// 		mockRepository.EXPECT().GetStatsEstate(gomock.Any(), gomock.Any()).Return(repository.GetStatsEstateOutput{}, sql.ErrNoRows)

// 		// Create a new request with the desired parameters
// 		req := httptest.NewRequest(http.MethodGet, "/estate/"+testID.String()+"/stats", nil)
// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		ctx := e.NewContext(req, rec)

// 		// Call the handler function
// 		err := s.GetEstateStats(ctx, testID)

// 		// Assertions
// 		assert.NoError(t, err)                         // Check that there was no error
// 		assert.Equal(t, http.StatusNotFound, rec.Code) // Check that the status code is 404
// 	})

// 	t.Run("GetEstateStats 200 Ok", func(t *testing.T) {
// 		// Create a new Echo instance
// 		e := echo.New()
// 		mockRepository := repository.NewMockRepositoryInterface(ctrl)
// 		s := NewServer(NewServerOptions{Repository: mockRepository})

// 		testID := uuid.New()

// 		mockResponse := repository.GetStatsEstateOutput{
// 			Count:  5,
// 			Min:    10,
// 			Max:    20,
// 			Median: 15,
// 		}

// 		// mocking repoisitory layer
// 		// mocking creating a new instance of estate
// 		mockRepository.EXPECT().GetStatsEstate(gomock.Any(), gomock.Any()).Return(mockResponse, nil)

// 		// Create a new request with the desired parameters
// 		req := httptest.NewRequest(http.MethodGet, "/estate/"+testID.String()+"/stats", nil)
// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		ctx := e.NewContext(req, rec)

// 		// Call the handler function
// 		err := s.GetEstateStats(ctx, testID)
// 		assert.NoError(t, err)

// 		res := repository.GetStatsEstateOutput{}
// 		err = json_iter.Unmarshal(rec.Body.Bytes(), &res)
// 		assert.NoError(t, err)

// 		// Assertions
// 		assert.NoError(t, err)                   // Check that there was no error
// 		assert.Equal(t, http.StatusOK, rec.Code) // Check that the status code is 200

// 		assert.NotNil(t, res) // Check that the response body is not empty
// 		assert.ObjectsAreEqual(mockResponse, res)
// 	})
// }

// func TestGetDronePlan(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	t.Run("GetDronePlan 400 Bad Request", func(t *testing.T) {
// 		// Create a new Echo instance
// 		e := echo.New()
// 		mockRepository := repository.NewMockRepositoryInterface(ctrl)
// 		s := NewServer(NewServerOptions{Repository: mockRepository})

// 		testID := uuid.Nil // force to get invalid request

// 		// mocking repoisitory layer

// 		// Create a new request with the desired parameters
// 		req := httptest.NewRequest(http.MethodGet, "/estate/"+testID.String()+"/drone-plan", nil)
// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		ctx := e.NewContext(req, rec)

// 		// Call the handler function
// 		err := s.GetDronePlan(ctx, testID, generated.GetDronePlanParams{})

// 		// Assertions
// 		assert.NoError(t, err)                           // Check that there was no error
// 		assert.Equal(t, http.StatusBadRequest, rec.Code) // Check that the status code is 404
// 	})

// 	t.Run("GetDronePlan 404 Not Found", func(t *testing.T) {
// 		// Create a new Echo instance
// 		e := echo.New()

// 		mockRepository := repository.NewMockRepositoryInterface(ctrl)
// 		s := NewServer(NewServerOptions{Repository: mockRepository})

// 		testID := uuid.New()

// 		// mocking repoisitory layer
// 		mockRepository.EXPECT().GetEstate(gomock.Any(), gomock.Any()).Return(repository.GetEstateOutput{}, sql.ErrNoRows)

// 		// Create a new request with the desired parameters
// 		req := httptest.NewRequest(http.MethodGet, "/estate/"+testID.String()+"/drone-plan", nil)
// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		ctx := e.NewContext(req, rec)

// 		// Call the handler function
// 		err := s.GetDronePlan(ctx, testID, generated.GetDronePlanParams{})

// 		// Assertions
// 		assert.NoError(t, err)                         // Check that there was no error
// 		assert.Equal(t, http.StatusNotFound, rec.Code) // Check that the status code is 404
// 	})

// 	t.Run("GetDronePlan 200 Ok", func(t *testing.T) {
// 		// Create a new Echo instance
// 		e := echo.New()
// 		mockRepository := repository.NewMockRepositoryInterface(ctrl)
// 		s := NewServer(NewServerOptions{Repository: mockRepository})

// 		testID := uuid.New()

// 		mockEstateRes := repository.GetEstateOutput{
// 			Id:        testID.String(),
// 			Length:    10,
// 			Width:     20,
// 			CreatedAt: time.Now().Format(time.RFC3339),
// 		}
// 		mockGetTrees := []repository.GetTreeOutput{
// 			{
// 				Id:        uuid.New().String(),
// 				EstateID:  testID.String(),
// 				X:         10,
// 				Y:         20,
// 				Height:    30,
// 				CreatedAt: time.Now().Format(time.RFC3339),
// 			},
// 		}

// 		// mocking repoisitory layer
// 		mockRepository.EXPECT().GetEstate(gomock.Any(), gomock.Any()).Return(mockEstateRes, nil)
// 		mockRepository.EXPECT().GetTrees(gomock.Any(), gomock.Any()).Return(mockGetTrees, nil)

// 		// Create a new request with the desired parameters
// 		req := httptest.NewRequest(http.MethodGet, "/estate/"+testID.String()+"/drone-plan", nil)
// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		ctx := e.NewContext(req, rec)

// 		// Call the handler function
// 		err := s.GetDronePlan(ctx, testID, generated.GetDronePlanParams{})
// 		assert.NoError(t, err)

// 		res := repository.GetStatsEstateOutput{}
// 		err = json_iter.Unmarshal(rec.Body.Bytes(), &res)
// 		assert.NoError(t, err)

// 		// Assertions
// 		assert.NoError(t, err)                   // Check that there was no error
// 		assert.Equal(t, http.StatusOK, rec.Code) // Check that the status code is 200

// 		assert.NotNil(t, res) // Check that the response body is not empty
// 		assert.Greater(t, len(mockGetTrees), 0)
// 	})

// 	t.Run("GetDronePlan with max distance 200 Ok", func(t *testing.T) {
// 		// Create a new Echo instance
// 		e := echo.New()
// 		mockRepository := repository.NewMockRepositoryInterface(ctrl)
// 		s := NewServer(NewServerOptions{Repository: mockRepository})

// 		testID := uuid.New()

// 		mockEstateRes := repository.GetEstateOutput{
// 			Id:        testID.String(),
// 			Length:    10,
// 			Width:     20,
// 			CreatedAt: time.Now().Format(time.RFC3339),
// 		}
// 		mockGetTrees := []repository.GetTreeOutput{
// 			{
// 				Id:        uuid.New().String(),
// 				EstateID:  testID.String(),
// 				X:         10,
// 				Y:         20,
// 				Height:    30,
// 				CreatedAt: time.Now().Format(time.RFC3339),
// 			},
// 		}

// 		// mocking repoisitory layer
// 		mockRepository.EXPECT().GetEstate(gomock.Any(), gomock.Any()).Return(mockEstateRes, nil)
// 		mockRepository.EXPECT().GetTrees(gomock.Any(), gomock.Any()).Return(mockGetTrees, nil)

// 		// Create a new request with the desired parameters
// 		req := httptest.NewRequest(http.MethodGet, "/estate/"+testID.String()+"/drone-plan?max_distance=100", nil)
// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		ctx := e.NewContext(req, rec)

// 		// Call the handler function
// 		maxDistanceSample := 100
// 		err := s.GetDronePlan(ctx, testID, generated.GetDronePlanParams{
// 			MaxDistance: &maxDistanceSample,
// 		})
// 		assert.NoError(t, err)

// 		res := generated.GeneratedDronePlanWithMaxDistanceResponse{}
// 		err = json_iter.Unmarshal(rec.Body.Bytes(), &res)
// 		assert.NoError(t, err)

// 		// Assertions
// 		assert.NoError(t, err)                   // Check that there was no error
// 		assert.Equal(t, http.StatusOK, rec.Code) // Check that the status code is 200

// 		assert.NotNil(t, res) // Check that the response body is not empty
// 		assert.Greater(t, len(mockGetTrees), 0)
// 		assert.NotNil(t, res.Distance)
// 	})
// }
