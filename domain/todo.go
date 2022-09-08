package domain

import (
	"context"
	"encoding/json"
	"time"
)

// ITodoUsecase interface is the contract for the Services
type ITodoUsecase interface {
	AddTodo(ctx context.Context, todo *Todo) (err error)
	Fetch(ctx context.Context, param FetchTodoParam) (res []Todo, cursor string, err error)
}

// FetchTodoParam used for fetching param
type FetchTodoParam struct {
	Limit  int64
	Cursor string
}

// ITodoRepository interface is the contract for the repositories
type ITodoRepository interface {
	AddTodo(ctx context.Context, todo *Todo) (err error)
	Fetch(ctx context.Context, param FetchTodoParam) (res []Todo, cursor string, err error)
}

// Todo represent the Todo data structure
type Todo struct {
	ID        string    `json:"id"`
	Text      string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

// MarshalJSON overide the default JSON Marshal
func (t Todo) MarshalJSON() (res []byte, err error) {
	item := struct {
		ID        string `json:"id"`
		Text      string `json:"message"`
		Status    string `json:"status"`
		CreatedAt string `json:"createdAt"`
	}{
		ID:        t.ID,
		Text:      t.Text,
		Status:    t.Status,
		CreatedAt: t.CreatedAt.Format(time.RFC3339),
	}

	return json.Marshal(item)
}