package parking_zone

import (
	"spotsync/internal/middlewares"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes mounts the parking zone routes onto the given Echo group.
//
// Public:
//
//	GET  /api/v1/zones
//	GET  /api/v1/zones/:id
//
// Admin only:
//
//	POST /api/v1/zones
func RegisterRoutes(g *echo.Group, h *Handler, jwtMW echo.MiddlewareFunc) {
	zones := g.Group("/zones")
	zones.GET("", h.GetAllZones)
	zones.GET("/:id", h.GetZoneByID)
	zones.POST("", h.CreateZone, jwtMW, middlewares.AdminOnly)
}
