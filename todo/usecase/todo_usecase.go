package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/LieAlbertTriAdrian/clean-arch-golang/internal/sqlx"
	repository "github.com/LieAlbertTriAdrian/clean-arch-golang/todo/repository/postgres"
	"github.com/sirupsen/logrus"

	domain "github.com/LieAlbertTriAdrian/clean-arch-golang/domain"
	"github.com/LieAlbertTriAdrian/clean-arch-golang/ebus"
)

type service struct {
	repo domain.ITodoRepository
}

// Events name
const (
	EventTodoCreated = "TODO_CREATED"
)

// NewService will initilaize the Todo service
func NewService(repo domain.ITodoRepository) domain.ITodoUsecase {
	return &service{
		repo: repo,
	}
}

func NewTxService(db *sql.DB) domain.ITodoUsecase {
	if db == nil {
		panic("missing db")
	}
	delegate := service{
		repo: repository.NewTodoRepository(),
	}
	return transactionalService{
		delegate: delegate,
		db:       db,
	}
}

func (s service) AddTodo(ctx context.Context, todo *domain.Todo) (err error) {
	err = s.repo.AddTodo(ctx, todo)
	if err != nil {
		return
	}
	event := ebus.Event{
		Name: EventTodoCreated,
		Data: todo,
	}
	// Event publishing should be done in the background
	eventCtx := context.Background()
	go ebus.Publish(eventCtx, event)
	return
}

func (s service) Fetch(ctx context.Context, param domain.FetchTodoParam) (res []domain.Todo, cursor string, err error) {
	res, cursor, err = s.repo.Fetch(ctx, param)
	return
}

type transactionalService struct {
	delegate service
	db       *sql.DB
}

func (s transactionalService) AddTodo(ctx context.Context, todo *domain.Todo) error {
	ctx, tx, err := sqlx.WithTx(ctx, s.db, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("could not begin transaction - %s", err)
	}

	err = s.delegate.AddTodo(ctx, todo)
	if err != nil {
		rErr := tx.Rollback()
		if rErr != nil && !errors.Is(rErr, sql.ErrTxDone) {
			logrus.WithError(rErr).Warnf("could not safely rollback transaction during #AddTodo")
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("could not finish committing - %s", err)
	}
	return nil
}

func (s transactionalService) Fetch(ctx context.Context, param domain.FetchTodoParam) ([]domain.Todo, string, error) {
	ctx, tx, err := sqlx.WithTx(ctx, s.db, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, "", fmt.Errorf("could not begin transaction - %s", err)
	}

	items, cursor, err := s.delegate.Fetch(ctx, param)
	rErr := tx.Rollback()
	if rErr != nil && !errors.Is(rErr, sql.ErrTxDone) {
		logrus.WithError(rErr).Warn("could not properly clean up transaction during #Fetch")
	}
	return items, cursor, err
}
