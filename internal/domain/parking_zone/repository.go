package parking_zone

import (
	"spotsync/internal/domain/reservation_model"

	"gorm.io/gorm"
)

// Repository defines the data-access contract for the parking_zone domain.
type Repository interface {
	Create(z *ParkingZone) error
	FindAll() ([]ParkingZone, error)
	FindByID(id uint) (*ParkingZone, error)
	CountActiveReservations(zoneID uint) (int64, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository returns a Repository backed by the given *gorm.DB.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(z *ParkingZone) error {
	return r.db.Create(z).Error
}

func (r *repository) FindAll() ([]ParkingZone, error) {
	var zones []ParkingZone
	return zones, r.db.Find(&zones).Error
}

func (r *repository) FindByID(id uint) (*ParkingZone, error) {
	var z ParkingZone
	if err := r.db.First(&z, id).Error; err != nil {
		return nil, err
	}
	return &z, nil
}

// CountActiveReservations queries the reservations table to get a live count
// of active bookings for a zone. Used to compute available_spots.
func (r *repository) CountActiveReservations(zoneID uint) (int64, error) {
	var count int64
	err := r.db.Model(&reservation_model.Reservation{}).
		Where("zone_id = ? AND status = ?", zoneID, "active").
		Count(&count).Error
	return count, err
}
