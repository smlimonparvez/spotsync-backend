package reservation

import (
	"errors"

	"spotsync/internal/domain/parking_zone"
	"spotsync/internal/domain/reservation_model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ErrZoneFull is returned when a parking zone has reached its total capacity.
var ErrZoneFull = errors.New("parking zone is at full capacity")

// Repository defines the data-access contract for the reservation domain.
type Repository interface {
	CreateWithLock(r *reservation_model.Reservation) error
	FindByUserID(userID uint) ([]Reservation, error)
	FindByID(id uint) (*reservation_model.Reservation, error)
	FindAll() ([]Reservation, error)
	Cancel(id uint) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository returns a Repository backed by the given *gorm.DB.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateWithLock(res *reservation_model.Reservation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var zone parking_zone.ParkingZone
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&zone, res.ZoneID).Error; err != nil {
			return err
		}

		var activeCount int64
		if err := tx.Model(&reservation_model.Reservation{}).
			Where("zone_id = ? AND status = ?", res.ZoneID, "active").
			Count(&activeCount).Error; err != nil {
			return err
		}

		if int(activeCount) >= zone.TotalCapacity {
			return ErrZoneFull
		}

		return tx.Create(res).Error
	})
}

func (r *repository) FindByUserID(userID uint) ([]Reservation, error) {
	var list []Reservation
	err := r.db.Preload("Zone").
		Where("user_id = ?", userID).
		Find(&list).Error
	return list, err
}

func (r *repository) FindByID(id uint) (*reservation_model.Reservation, error) {
	var res reservation_model.Reservation
	if err := r.db.First(&res, id).Error; err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *repository) FindAll() ([]Reservation, error) {
	var list []Reservation
	err := r.db.Preload("User").Preload("Zone").Find(&list).Error
	return list, err
}

func (r *repository) Cancel(id uint) error {
	return r.db.Model(&reservation_model.Reservation{}).
		Where("id = ?", id).
		Update("status", "cancelled").Error
}
