package parking_zone

import (
	"net/http"
	"strconv"

	"spotsync/internal/domain/parking_zone/dto"
	"spotsync/internal/httpresponse"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Handler holds HTTP handlers for the parking_zone domain.
type Handler struct {
	svc      Service
	validate *validator.Validate
}

// NewHandler returns a Handler with the given Service.
func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc, validate: validator.New()}
}

// CreateZone handles POST /api/v1/zones  (admin only)
func (h *Handler) CreateZone(c echo.Context) error {
	var req dto.CreateZoneRequest
	if err := c.Bind(&req); err != nil {
		return httpresponse.BadRequest(c, "Invalid request body", err.Error())
	}
	if err := h.validate.Struct(req); err != nil {
		return httpresponse.BadRequest(c, "Validation failed", err.Error())
	}

	zone, err := h.svc.CreateZone(&req)
	if err != nil {
		return httpresponse.InternalServerError(c, "Failed to create parking zone")
	}
	return httpresponse.Success(c, http.StatusCreated, "Parking zone created successfully", zone)
}

// GetAllZones handles GET /api/v1/zones  (public)
func (h *Handler) GetAllZones(c echo.Context) error {
	zones, err := h.svc.GetAllZones()
	if err != nil {
		return httpresponse.InternalServerError(c, "Failed to retrieve parking zones")
	}
	return httpresponse.Success(c, http.StatusOK, "Parking zones retrieved successfully", zones)
}

// GetZoneByID handles GET /api/v1/zones/:id  (public)
func (h *Handler) GetZoneByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return httpresponse.BadRequest(c, "Invalid zone ID", nil)
	}

	zone, err := h.svc.GetZoneByID(uint(id))
	if err != nil {
		return httpresponse.NotFound(c, "Parking zone not found")
	}
	return httpresponse.Success(c, http.StatusOK, "Parking zone retrieved successfully", zone)
}
