package reservation

import (
	"errors"
	"net/http"
	"strconv"

	"spotsync/internal/domain/reservation/dto"
	"spotsync/internal/httpresponse"
	"spotsync/internal/middlewares"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Handler holds HTTP handlers for the reservation domain.
type Handler struct {
	svc      Service
	validate *validator.Validate
}

// NewHandler returns a Handler with the given Service.
func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc, validate: validator.New()}
}

// CreateReservation handles POST /api/v1/reservations
func (h *Handler) CreateReservation(c echo.Context) error {
	userID := middlewares.GetUserID(c)

	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return httpresponse.BadRequest(c, "Invalid request body", err.Error())
	}
	if err := h.validate.Struct(req); err != nil {
		return httpresponse.BadRequest(c, "Validation failed", err.Error())
	}

	res, err := h.svc.CreateReservation(userID, &req)
	if err != nil {
		if errors.Is(err, ErrZoneFull) {
			return httpresponse.Conflict(c, "Parking zone is at full capacity")
		}
		if err.Error() == "parking zone not found" {
			return httpresponse.NotFound(c, "Parking zone not found")
		}
		return httpresponse.InternalServerError(c, "Failed to create reservation")
	}

	return httpresponse.Success(c, http.StatusCreated, "Reservation confirmed successfully", res)
}

// GetMyReservations handles GET /api/v1/reservations/my-reservations
func (h *Handler) GetMyReservations(c echo.Context) error {
	userID := middlewares.GetUserID(c)

	list, err := h.svc.GetMyReservations(userID)
	if err != nil {
		return httpresponse.InternalServerError(c, "Failed to retrieve reservations")
	}
	return httpresponse.Success(c, http.StatusOK, "My reservations retrieved successfully", list)
}

// CancelReservation handles DELETE /api/v1/reservations/:id
func (h *Handler) CancelReservation(c echo.Context) error {
	userID := middlewares.GetUserID(c)
	role := middlewares.GetRole(c)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return httpresponse.BadRequest(c, "Invalid reservation ID", nil)
	}

	if err := h.svc.CancelReservation(uint(id), userID, role); err != nil {
		switch err.Error() {
		case "reservation not found":
			return httpresponse.NotFound(c, "Reservation not found")
		case "forbidden":
			return httpresponse.Forbidden(c, "You can only cancel your own reservations")
		case "already cancelled":
			return httpresponse.Conflict(c, "Reservation is already cancelled")
		default:
			return httpresponse.InternalServerError(c, "Failed to cancel reservation")
		}
	}

	return httpresponse.Success(c, http.StatusOK, "Reservation cancelled successfully", nil)
}

// GetAllReservations handles GET /api/v1/reservations  (admin only)
func (h *Handler) GetAllReservations(c echo.Context) error {
	list, err := h.svc.GetAllReservations()
	if err != nil {
		return httpresponse.InternalServerError(c, "Failed to retrieve reservations")
	}
	return httpresponse.Success(c, http.StatusOK, "All reservations retrieved successfully", list)
}
