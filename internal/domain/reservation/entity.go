package reservation

import (
	"spotsync/internal/domain/parking_zone"
	"spotsync/internal/domain/reservation_model"
	"spotsync/internal/domain/user"
)

// Reservation re-exports the shared model and adds GORM-loadable relations.
// The struct lives in reservation_model to avoid circular imports; here we
// compose it with User and ParkingZone for Preload-based queries.
type Reservation struct {
	reservation_model.Reservation
	User user.User                `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Zone parking_zone.ParkingZone `gorm:"foreignKey:ZoneID" json:"zone,omitempty"`
}
