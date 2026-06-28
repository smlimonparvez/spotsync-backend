package user

import (
	"net/http"

	"spotsync/internal/domain/user/dto"
	"spotsync/internal/httpresponse"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Handler holds HTTP handlers for the user domain.
type Handler struct {
	svc      Service
	validate *validator.Validate
}

// NewHandler returns a Handler with the given Service.
func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc, validate: validator.New()}
}

// Register handles POST /api/v1/auth/register
func (h *Handler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return httpresponse.BadRequest(c, "Invalid request body", err.Error())
	}
	if err := h.validate.Struct(req); err != nil {
		return httpresponse.BadRequest(c, "Validation failed", err.Error())
	}

	resp, err := h.svc.Register(&req)
	if err != nil {
		if err.Error() == "email already registered" {
			return httpresponse.Conflict(c, err.Error())
		}
		return httpresponse.InternalServerError(c, "Failed to register user")
	}

	return httpresponse.Success(c, http.StatusCreated, "User registered successfully", resp)
}

// Login handles POST /api/v1/auth/login
func (h *Handler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return httpresponse.BadRequest(c, "Invalid request body", err.Error())
	}
	if err := h.validate.Struct(req); err != nil {
		return httpresponse.BadRequest(c, "Validation failed", err.Error())
	}

	resp, err := h.svc.Login(&req)
	if err != nil {
		return httpresponse.Error(c, http.StatusUnauthorized, err.Error(), nil)
	}

	return httpresponse.Success(c, http.StatusOK, "Login successful", resp)
}
