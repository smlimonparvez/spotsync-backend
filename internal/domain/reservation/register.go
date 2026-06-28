package reservation

import (
	"spotsync/internal/middlewares"

	"github.com/labstack/echo/v4"
)

// Authenticated (driver + admin):
//	POST   /api/v1/reservations
//	GET    /api/v1/reservations/my-reservations
//	DELETE /api/v1/reservations/:id

// Admin only:
//	GET    /api/v1/reservations

func RegisterRoutes(g *echo.Group, h *Handler, jwtMW echo.MiddlewareFunc) {
	r := g.Group("/reservations", jwtMW)
	r.POST("", h.CreateReservation)
	r.GET("/my-reservations", h.GetMyReservations)
	r.DELETE("/:id", h.CancelReservation)
	r.GET("", h.GetAllReservations, middlewares.AdminOnly)
}
