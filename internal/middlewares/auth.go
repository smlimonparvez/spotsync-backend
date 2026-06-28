package middlewares

import (
	"strings"

	"spotsync/internal/auth"
	"spotsync/internal/httpresponse"

	"github.com/labstack/echo/v5"
)

// JWT returns an Echo middleware that validates Bearer tokens and injects
// user_id and role into the Echo context for downstream handlers.
func JWT(jwtSvc *auth.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				return httpresponse.Unauthorized(c, "Missing or invalid Authorization header")
			}

			tokenStr := strings.TrimPrefix(header, "Bearer ")
			claims, err := jwtSvc.ParseToken(tokenStr)
			if err != nil {
				return httpresponse.Unauthorized(c, "Invalid or expired token")
			}

			c.Set("user_id", claims.UserID)
			c.Set("role", claims.Role)
			return next(c)
		}
	}
}

// AdminOnly is an Echo middleware that allows only users with role "admin".
// Must be placed after JWT middleware.
func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, _ := c.Get("role").(string)
		if role != "admin" {
			return httpresponse.Forbidden(c, "Access denied: admin role required")
		}
		return next(c)
	}
}

// GetUserID is a helper to safely extract the authenticated user's ID from context.
func GetUserID(c echo.Context) uint {
	id, _ := c.Get("user_id").(uint)
	return id
}

// GetRole is a helper to extract the authenticated user's role from context.
func GetRole(c echo.Context) string {
	role, _ := c.Get("role").(string)
	return role
}
