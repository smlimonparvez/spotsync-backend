package httpresponse

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

// Envelope is the standard JSON wrapper for every API response.
type Envelope struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func Success(c echo.Context, status int, message string, data interface{}) error {
	return c.JSON(status, Envelope{Success: true, Message: message, Data: data})
}

func Error(c echo.Context, status int, message string, errs interface{}) error {
	return c.JSON(status, Envelope{Success: false, Message: message, Errors: errs})
}

// Convenience helpers for common HTTP error statuses
func BadRequest(c echo.Context, message string, errs interface{}) error {
	return Error(c, http.StatusBadRequest, message, errs)
}

func Unauthorized(c echo.Context, message string) error {
	return Error(c, http.StatusUnauthorized, message, nil)
}

func Forbidden(c echo.Context, message string) error {
	return Error(c, http.StatusForbidden, message, nil)
}

func NotFound(c echo.Context, message string) error {
	return Error(c, http.StatusNotFound, message, nil)
}

func Conflict(c echo.Context, message string) error {
	return Error(c, http.StatusConflict, message, nil)
}

func InternalServerError(c echo.Context, message string) error {
	return Error(c, http.StatusInternalServerError, message, nil)
}
