package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/LieAlbertTriAdrian/clean-arch-golang/domain"
	"github.com/LieAlbertTriAdrian/clean-arch-golang/mocks"
	_todoUsecase "github.com/LieAlbertTriAdrian/clean-arch-golang/todo/usecase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddTodo(t *testing.T) {
	mockTodo := domain.Todo{
		Text:   "Eating breakfast",
		Status: "in-progress",
	}
	testCases := []struct {
		desc        string
		repo        domain.ITodoRepository
		expectedRes error
		ctx         context.Context
		todo        *domain.Todo
	}{
		{
			desc: "created",
			repo: func() domain.ITodoRepository {
				mockRepo := new(mocks.ITodoRepository)
				mockRepo.On("AddTodo", mock.Anything, mock.AnythingOfType("*domain.Todo")).Return(nil).Once()
				return mockRepo
			}(),
			expectedRes: nil,
			ctx:         context.Background(),
			todo:        &mockTodo,
		},
		{
			desc: "repo error",
			repo: func() domain.ITodoRepository {
				mockRepo := new(mocks.ITodoRepository)
				mockRepo.On("AddTodo", mock.Anything, mock.AnythingOfType("*domain.Todo")).Return(errors.New("DB Error")).Once()
				return mockRepo
			}(),
			expectedRes: errors.New("DB Error"),
			ctx:         context.Background(),
			todo:        &mockTodo,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			srv := _todoUsecase.NewService(tC.repo)
			err := srv.AddTodo(tC.ctx, tC.todo)
			require.Equal(t, tC.expectedRes, err)
		})
	}
}

func TestFetchTodo(t *testing.T) {
	mockFetchTodoParam := domain.FetchTodoParam{
		Limit:  100,
		Cursor: "mock_cursor",
	}
	mockTodos := []domain.Todo{{
		Text:   "Brewing coffee",
		Status: "in-progress",
	}}
	mockNextCursor := "next_cursor"

	type resType struct {
		Res    []domain.Todo
		Cursor string
		Err    error
	}
	testCases := []struct {
		desc           string
		repo           domain.ITodoRepository
		expectedRes    resType
		ctx            context.Context
		fetchTodoParam domain.FetchTodoParam
	}{
		{
			desc: "fetched",
			repo: func() domain.ITodoRepository {
				mockRepo := new(mocks.ITodoRepository)
				mockRepo.On("Fetch", mock.Anything, mock.AnythingOfType("domain.FetchTodoParam")).Return(mockTodos, mockNextCursor, nil).Once()
				return mockRepo
			}(),
			expectedRes: resType{
				Res:    mockTodos,
				Cursor: mockNextCursor,
			},
			ctx:            context.Background(),
			fetchTodoParam: mockFetchTodoParam,
		},
		{
			desc: "repo error",
			repo: func() domain.ITodoRepository {
				mockRepo := new(mocks.ITodoRepository)
				mockRepo.On("Fetch", mock.Anything, mock.AnythingOfType("domain.FetchTodoParam")).Return(nil, "", errors.New("DB error")).Once()
				return mockRepo
			}(),
			expectedRes: resType{
				Err: errors.New("DB error"),
			},
			ctx:            context.Background(),
			fetchTodoParam: mockFetchTodoParam,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			srv := _todoUsecase.NewService(tC.repo)
			res, cur, err := srv.Fetch(tC.ctx, tC.fetchTodoParam)
			require.Equal(t, tC.expectedRes.Res, res)
			require.Equal(t, tC.expectedRes.Cursor, cur)
			require.Equal(t, tC.expectedRes.Err, err)
		})
	}
}
