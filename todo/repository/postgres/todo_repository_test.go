package postgres_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/LieAlbertTriAdrian/clean-arch-golang/internal/sqlx"

	"github.com/LieAlbertTriAdrian/clean-arch-golang/domain"
	postgres "github.com/LieAlbertTriAdrian/clean-arch-golang/internal/postgres"
	repository "github.com/LieAlbertTriAdrian/clean-arch-golang/todo/repository/postgres"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
)

type todoTestSuite struct {
	postgres.Suite
}

func TestSuiteTodo(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip the Test Suite for Todo Repository")
	}

	dsn := os.Getenv("POSTGRES_TEST_URL")
	if dsn == "" {
		dsn = "user=postgres password=password dbname=testing host=localhost port=54320 sslmode=disable"
	}

	todoSuite := &todoTestSuite{
		postgres.Suite{
			DSN:                     dsn,
			MigrationLocationFolder: "../../../internal/postgres/migrations",
		},
	}

	suite.Run(t, todoSuite)
}

func (s todoTestSuite) BeforeTest(suiteName, testName string) {
	ok, err := s.Migration.Up()
	s.Require().NoError(err)
	s.Require().True(ok)
}

func (s todoTestSuite) AfterTest(suiteName, testName string) {
	ok, err := s.Migration.Down()
	s.Require().NoError(err)
	s.Require().True(ok)
}

func (s todoTestSuite) createContext() context.Context {
	ctx, _, err := sqlx.WithTx(context.TODO(), s.DBConn, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	s.Require().NoError(err)
	return ctx
}

func (s todoTestSuite) TestCreateTodo() {
	mockTodo := domain.Todo{
		Text:   "Eating Breakfast",
		Status: "todo",
	}

	testCases := []struct {
		desc           string
		repo           domain.ITodoRepository
		expectedResult error
		ctx            context.Context
		reqBody        *domain.Todo
	}{
		{
			desc: "insert-success",
			repo: func() domain.ITodoRepository {
				repo := repository.NewTodoRepository()
				return repo
			}(),
			expectedResult: nil,
			ctx:            s.createContext(),
			reqBody:        &mockTodo,
		},
		{
			desc: "context timeout. Too long to execute and already pass the limit context from parent",
			repo: func() domain.ITodoRepository {
				repo := repository.NewTodoRepository()
				return repo
			}(),
			expectedResult: context.DeadlineExceeded,
			ctx: func() context.Context {
				bCtx := s.createContext()
				// Context already expired for 1 hour
				ctx, cancel := context.WithDeadline(bCtx, time.Now().Add(-1*time.Hour))
				defer cancel()
				return ctx
			}(),
			reqBody: &mockTodo,
		},
		{
			desc: "context Canceled by the caller",
			repo: func() domain.ITodoRepository {
				repo := repository.NewTodoRepository()
				return repo
			}(),
			expectedResult: context.Canceled,
			ctx: func() context.Context {
				bCtx := s.createContext()
				// Context expired in 1 hour
				ctx, cancel := context.WithDeadline(bCtx, time.Now().Add(1*time.Hour))
				// Directly call cancel function
				defer cancel()
				return ctx
			}(),
			reqBody: &mockTodo,
		},
	}
	for _, tC := range testCases {
		s.T().Run(tC.desc, func(t *testing.T) {
			err := tC.repo.AddTodo(tC.ctx, tC.reqBody)
			if err == nil {
				s.Require().NoError(sqlx.Commit(tC.ctx))
			} else {
				s.Require().NoError(sqlx.Rollback(tC.ctx))
			}
			s.Require().Equal(tC.expectedResult, errors.Cause(err))
		})
	}

}

func (s todoTestSuite) TestFetchTodo() {
	mockTodos := []domain.Todo{
		{
			Text:      "Brewing coffee",
			Status:    "in-progress",
			CreatedAt: time.Now(),
		},
		{
			Text:      "Eating Breakfast",
			Status:    "PENDING",
			CreatedAt: time.Now().Add(-1 * time.Second),
		},
		{
			Text:      "Take a shower",
			Status:    "PENDING",
			CreatedAt: time.Now().Add(-2 * time.Second),
		},
	}
	// Seed the items
	repo := repository.NewTodoRepository()
	ctx := s.createContext()
	for _, item := range mockTodos {
		itemPointer := item
		err := repo.AddTodo(ctx, &itemPointer)
		s.Require().NoError(err)
	}
	s.Require().NoError(sqlx.Commit(ctx))

	// Testcase Table
	testCases := []struct {
		desc           string
		params         domain.FetchTodoParam
		expectedRes    []domain.Todo
		expectedErr    error
		expectedCursor string
		ctx            context.Context
	}{
		{
			desc: "fetch-success-simple",
			params: domain.FetchTodoParam{
				Limit:  1,
				Cursor: "",
			},
			expectedRes:    mockTodos[0:1],
			expectedErr:    nil,
			expectedCursor: postgres.EncodeCursor(mockTodos[0].CreatedAt.Unix()),
			ctx:            s.createContext(),
		},
		{
			desc: "fetch-success-with-cursor",
			params: domain.FetchTodoParam{
				Limit:  1,
				Cursor: postgres.EncodeCursor(mockTodos[0].CreatedAt.Unix()),
			},
			expectedRes:    mockTodos[1:2],
			expectedErr:    nil,
			expectedCursor: postgres.EncodeCursor(mockTodos[1].CreatedAt.Unix()),
			ctx:            s.createContext(),
		},
	}

	for _, tC := range testCases {
		s.T().Run(tC.desc, func(t *testing.T) {
			res, csr, err := repo.Fetch(tC.ctx, tC.params)
			s.Require().NoError(sqlx.Rollback(tC.ctx)) // we're not modifying anything, so just rollback

			s.Require().Len(res, len(tC.expectedRes))
			if len(tC.expectedRes) > 0 {
				for i, item := range res {
					s.Require().Equal(tC.expectedRes[i].Text, item.Text)
					s.Require().Equal(tC.expectedRes[i].Status, item.Status)
				}
			}
			s.Require().Equal(tC.expectedCursor, csr)
			s.Require().Equal(tC.expectedErr, err)
		})
	}
}
