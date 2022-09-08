package postgres

import (
	"context"
	"time"

	"github.com/LieAlbertTriAdrian/clean-arch-golang/internal/sqlx"

	domain "github.com/LieAlbertTriAdrian/clean-arch-golang/domain"
	"github.com/LieAlbertTriAdrian/clean-arch-golang/internal/postgres"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type todoHandler struct {
}

// NewTodoRepository will return
func NewTodoRepository() domain.ITodoRepository {
	return &todoHandler{}
}

func (t todoHandler) AddTodo(ctx context.Context, todo *domain.Todo) (err error) {
	tx, ok := sqlx.TxFrom(ctx)
	if !ok {
		return sqlx.ErrMissingTx
	}

	query := "INSERT INTO todo (id, text, status, created_at) VALUES ($1, $2, $3, $4)"
	now := time.Now()
	uuid, err := uuid.NewRandom()
	if err != nil {
		return
	}
	todo.ID = uuid.String()
	if todo.CreatedAt.IsZero() {
		todo.CreatedAt = now
	}
	res, err := tx.ExecContext(ctx, query, todo.ID, todo.Text, todo.Status, todo.CreatedAt)
	if err != nil {
		return
	}
	if _, err = res.RowsAffected(); err != nil {
		return
	}
	return
}

func (t todoHandler) Fetch(ctx context.Context, param domain.FetchTodoParam) (res []domain.Todo, cursor string, err error) {
	tx, ok := sqlx.TxFrom(ctx)
	if !ok {
		return []domain.Todo{}, "", sqlx.ErrMissingTx
	}

	builder := squirrel.Select("id", "text", "status", "created_at").PlaceholderFormat(squirrel.Dollar)
	builder = builder.OrderBy("created_at DESC").From("todo")
	if param.Limit > 0 {
		builder = builder.Limit(uint64(param.Limit))
	}

	if param.Cursor != "" {
		tmstmp, er := postgres.DecodeCursor(param.Cursor)
		if er != nil {
			err = er
			return
		}
		createdAt := time.Unix(tmstmp, 0)
		builder = builder.Where(squirrel.LtOrEq{
			"created_at": createdAt,
		})
	}
	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return
	}

	rows, err := tx.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	res = []domain.Todo{}
	for rows.Next() {
		t := domain.Todo{}
		err = rows.Scan(
			&t.ID,
			&t.Text,
			&t.Status,
			&t.CreatedAt,
		)
		res = append(res, t)
	}

	cursor = param.Cursor
	if len(res) > 0 {
		cursor = postgres.EncodeCursor(res[len(res)-1].CreatedAt.Unix())
	}
	return
}
