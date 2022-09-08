package rest

import (
	"net/http"
	"strconv"

	domain "github.com/LieAlbertTriAdrian/clean-arch-golang/domain"
	echo "github.com/labstack/echo/v4"
)

type todoHandler struct {
	service domain.ITodoUsecase
}

// InitTodoHandler will initialize the TodoHandler for REST HTTP
func InitTodoHandler(e *echo.Echo, service domain.ITodoUsecase) {
	h := &todoHandler{
		service: service,
	}

	e.POST("/todos", h.AddTodo)
	e.GET("/todos", h.FetchTodo)
}

func (t todoHandler) AddTodo(c echo.Context) (err error) {
	todo := domain.Todo{}
	err = c.Bind(&todo)
	if err != nil {
		return
	}

	err = t.service.AddTodo(c.Request().Context(), &todo)
	if err != nil {
		return
	}

	return c.JSON(http.StatusCreated, todo)
}

func (t todoHandler) FetchTodo(c echo.Context) (err error) {
	cursor := c.QueryParam("cursor")
	limitStr := c.QueryParam("limit")

	limit := 20
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return
		}
	}

	param := domain.FetchTodoParam{
		Cursor: cursor,
		Limit:  int64(limit),
	}

	res, nextCursor, err := t.service.Fetch(c.Request().Context(), param)
	if err != nil {
		return
	}

	c.Response().Header().Set("X-Cursor", nextCursor)
	return c.JSON(http.StatusOK, res)
}
