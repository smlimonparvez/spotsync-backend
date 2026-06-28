package server

import (
	"fmt"
	"log"

	"spotsync/internal/auth"
	"spotsync/internal/config"
	"spotsync/internal/domain/parking_zone"
	"spotsync/internal/domain/reservation"
	"spotsync/internal/domain/reservation_model"
	"spotsync/internal/domain/user"
	"spotsync/internal/middlewares"

	"github.com/labstack/echo/v4"
	echoMW "github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

// Server holds Echo and app config.
type Server struct {
	echo *echo.Echo
	cfg  *config.Config
}

// New wires all dependencies, auto-migrates the DB, and returns a ready Server.
func New(cfg *config.Config, db *gorm.DB) *Server {
	// Auto-migrate all domain entities
	if err := db.AutoMigrate(
		&user.User{},
		&parking_zone.ParkingZone{},
		&reservation_model.Reservation{},
	); err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}
	log.Println("✅ Database auto-migrated")

	// ─── Dependency Injection ────────────────────────────────────────────────
	jwtSvc := auth.NewJWTService(cfg.JWTSecret)

	// User domain
	userRepo := user.NewRepository(db)
	userSvc := user.NewService(userRepo, jwtSvc)
	userHandler := user.NewHandler(userSvc)

	// Parking zone domain
	zoneRepo := parking_zone.NewRepository(db)
	zoneSvc := parking_zone.NewService(zoneRepo)
	zoneHandler := parking_zone.NewHandler(zoneSvc)

	// Reservation domain
	resRepo := reservation.NewRepository(db)
	resSvc := reservation.NewService(resRepo, zoneRepo)
	resHandler := reservation.NewHandler(resSvc)
	// ─────────────────────────────────────────────────────────────────────────

	// JWT middleware factory (carries jwtSvc)
	jwtMW := middlewares.JWT(jwtSvc)

	// ─── Echo setup ──────────────────────────────────────────────────────────
	e := echo.New()
	e.HideBanner = true

	e.Use(echoMW.Logger())
	e.Use(echoMW.Recover())
	e.Use(echoMW.CORSWithConfig(echoMW.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
	}))

	// ─── Route registration ──────────────────────────────────────────────────
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello World SpotSync API running!")
	})

	api := e.Group("/api/v1")

	user.RegisterRoutes(api, userHandler)
	parking_zone.RegisterRoutes(api, zoneHandler, jwtMW)
	reservation.RegisterRoutes(api, resHandler, jwtMW)
	// ─────────────────────────────────────────────────────────────────────────

	return &Server{echo: e, cfg: cfg}
}

// Start begins listening on the configured port.
func (s *Server) Start() {
	addr := fmt.Sprintf(":%s", s.cfg.Port)
	fmt.Printf("SpotSync API running on %s\n", addr)
	if err := s.echo.Start(addr); err != nil {
		log.Fatal("Server failed:", err)
	}
}
