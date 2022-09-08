package rest_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/LieAlbertTriAdrian/clean-arch-golang/domain"
	"github.com/LieAlbertTriAdrian/clean-arch-golang/internal/rest/middleware"
	"github.com/LieAlbertTriAdrian/clean-arch-golang/mocks"
	"github.com/LieAlbertTriAdrian/clean-arch-golang/todo/delivery/rest"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var e *echo.Echo

func TestMain(m *testing.M) {
	e = echo.New()
	e.Use(
		middleware.ErrorMiddleware(),
		middleware.LogErrorMiddleware(),
	)
	os.Exit(m.Run())
}

func TestAddTodo(t *testing.T) {
	testCases := []struct {
		desc               string
		service            domain.ITodoUsecase
		URL                string
		method             string
		reqBody            io.Reader
		expectedStatusCode int
	}{
		{
			desc: "success",
			service: func() domain.ITodoUsecase {
				mockService := new(mocks.ITodoUsecase)
				mockService.On("AddTodo", mock.Anything, mock.AnythingOfType("*domain.Todo")).Return(nil).
					Once()
				return mockService
			}(),
			URL:    "/todos",
			method: http.MethodPost,
			reqBody: strings.NewReader(`
				{
					"text": "Sleep well at 10 PM",
					"status": "PENDING"
				}
			`),
			expectedStatusCode: http.StatusCreated,
		},
		{
			desc: "service error",
			service: func() domain.ITodoUsecase {
				mockService := new(mocks.ITodoUsecase)
				mockService.On("AddTodo", mock.Anything, mock.AnythingOfType("*domain.Todo")).Return(errors.New("Database Error")).
					Once()
				return mockService
			}(),
			URL:    "/todos",
			method: http.MethodPost,
			reqBody: strings.NewReader(`
				{
					"text": "Sleep well at 10 PM",
					"status": "PENDING"
				}
			`),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			rest.InitTodoHandler(e, tC.service)
			req := httptest.NewRequest(tC.method, tC.URL, tC.reqBody)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			require.Equal(t, tC.expectedStatusCode, rec.Code)
		})
	}
}

func TestFetchTodo(t *testing.T) {
	testCases := []struct {
		desc               string
		service            domain.ITodoUsecase
		URL                string
		method             string
		reqBody            io.Reader
		expectedStatusCode int
	}{
		{
			desc: "success",
			service: func() domain.ITodoUsecase {
				mockService := new(mocks.ITodoUsecase)
				res := []domain.Todo{}
				mockService.On("Fetch", mock.Anything, mock.AnythingOfType("domain.FetchTodoParam")).Return(res, "next_cursor", nil).Once()
				return mockService
			}(),
			URL:                "/todos?cursor=fake_cursor&limit=100",
			method:             http.MethodGet,
			reqBody:            nil,
			expectedStatusCode: http.StatusOK,
		},
		{
			desc: "service error",
			service: func() domain.ITodoUsecase {
				mockService := new(mocks.ITodoUsecase)
				mockService.On("Fetch", mock.Anything, mock.AnythingOfType("domain.FetchTodoParam")).Return(nil, "", errors.New("Database Error")).Once()
				return mockService
			}(),
			URL:                "/todos?cursor=fake_cursor&limit=100",
			method:             http.MethodGet,
			reqBody:            nil,
			expectedStatusCode: http.StatusInternalServerError,
		},
		// TODO: add more testcase including response body checking
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			rest.InitTodoHandler(e, tC.service)
			req := httptest.NewRequest(tC.method, tC.URL, tC.reqBody)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			require.Equal(t, tC.expectedStatusCode, rec.Code)
		})
	}
}
