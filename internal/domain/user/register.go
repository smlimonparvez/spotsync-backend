package user

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes mounts the user/auth routes onto the given Echo group.
// Public routes (no auth required):
//
//	POST /api/v1/auth/register
//	POST /api/v1/auth/login
func RegisterRoutes(g *echo.Group, h *Handler) {
	auth := g.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
}
