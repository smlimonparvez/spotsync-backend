package reservation_model

import "time"

// Reservation is the GORM model for the "reservations" table.
type Reservation struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint      `gorm:"not null"                 json:"user_id"`
	ZoneID       uint      `gorm:"not null"                 json:"zone_id"`
	LicensePlate string    `gorm:"not null;size:15"         json:"license_plate"`
	Status       string    `gorm:"not null;default:'active'"json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
