package middleware

import (
	"context"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/LieAlbertTriAdrian/clean-arch-golang/ebus"
)

// LogErrorMiddleware will log every error for every HTTP Request made by client
func LogErrorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"method": c.Request().Method,
					"uri":    c.Request().RequestURI,
					"err":    err,
				}).Error("Got http error")
			}
			return err
		}
	}
}

// ErrorMiddleware is a function to generate http status code
func ErrorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			message := err.Error()
			err = errors.Cause(err)
			switch err {
			case context.DeadlineExceeded, context.Canceled:
				return echo.NewHTTPError(http.StatusRequestTimeout, err.Error())
			}

			if echoErr, ok := err.(*echo.HTTPError); ok {
				return echoErr
			}

			// Add another Error handling here
			return echo.NewHTTPError(http.StatusInternalServerError, message)
		}
	}
}

// EbusInjectorToRequestContext is a function to inject ebus to every request context
func EbusInjectorToRequestContext(pub ebus.Publisher) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			// merge request context with user context
			ctx := ebus.ContextWithPublisher(c.Request().Context(), pub)
			// recompose request with a new context
			httpReq := c.Request().WithContext(ctx)
			c.SetRequest(httpReq)
			return next(c)
		}
	}
}
