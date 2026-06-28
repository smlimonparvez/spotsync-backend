package dto

import "time"

// CreateReservationRequest is the JSON body for POST /api/v1/reservations
type CreateReservationRequest struct {
	ZoneID       uint   `json:"zone_id"       validate:"required,gt=0"`
	LicensePlate string `json:"license_plate" validate:"required,max=15"`
}

// ReservationResponse is returned after a successful booking.
type ReservationResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	ZoneID       uint      `json:"zone_id"`
	LicensePlate string    `json:"license_plate"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ZoneSummary is the minimal zone info embedded in reservation responses.
type ZoneSummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// UserSummary is the minimal user info embedded in admin reservation responses.
type UserSummary struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// MyReservationResponse is the driver's view of their own reservations.
type MyReservationResponse struct {
	ID           uint        `json:"id"`
	LicensePlate string      `json:"license_plate"`
	Status       string      `json:"status"`
	Zone         ZoneSummary `json:"zone"`
	CreatedAt    time.Time   `json:"created_at"`
}

// AdminReservationResponse is the admin's view including full user + zone info.
type AdminReservationResponse struct {
	ID           uint        `json:"id"`
	LicensePlate string      `json:"license_plate"`
	Status       string      `json:"status"`
	User         UserSummary `json:"user"`
	Zone         ZoneSummary `json:"zone"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}
