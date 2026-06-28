package reservation

import (
	"errors"

	"spotsync/internal/domain/parking_zone"
	"spotsync/internal/domain/reservation/dto"
	"spotsync/internal/domain/reservation_model"

	"gorm.io/gorm"
)

// Service defines the business-logic contract for the reservation domain.
type Service interface {
	CreateReservation(userID uint, req *dto.CreateReservationRequest) (*dto.ReservationResponse, error)
	GetMyReservations(userID uint) ([]dto.MyReservationResponse, error)
	CancelReservation(reservationID, userID uint, role string) error
	GetAllReservations() ([]dto.AdminReservationResponse, error)
}

type service struct {
	repo     Repository
	zoneRepo parking_zone.Repository
}

// NewService returns a Service with the given repositories.
func NewService(repo Repository, zoneRepo parking_zone.Repository) Service {
	return &service{repo: repo, zoneRepo: zoneRepo}
}

func (s *service) CreateReservation(userID uint, req *dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	// Verify the zone exists before entering the transaction
	if _, err := s.zoneRepo.FindByID(req.ZoneID); err != nil {
		return nil, errors.New("parking zone not found")
	}

	res := &reservation_model.Reservation{
		UserID:       userID,
		ZoneID:       req.ZoneID,
		LicensePlate: req.LicensePlate,
		Status:       "active",
	}

	if err := s.repo.CreateWithLock(res); err != nil {
		if errors.Is(err, ErrZoneFull) {
			return nil, ErrZoneFull
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("parking zone not found")
		}
		return nil, err
	}

	return &dto.ReservationResponse{
		ID:           res.ID,
		UserID:       res.UserID,
		ZoneID:       res.ZoneID,
		LicensePlate: res.LicensePlate,
		Status:       res.Status,
		CreatedAt:    res.CreatedAt,
		UpdatedAt:    res.UpdatedAt,
	}, nil
}

func (s *service) GetMyReservations(userID uint) ([]dto.MyReservationResponse, error) {
	list, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MyReservationResponse, 0, len(list))
	for _, r := range list {
		responses = append(responses, dto.MyReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			Zone:         dto.ZoneSummary{ID: r.Zone.ID, Name: r.Zone.Name, Type: r.Zone.Type},
			CreatedAt:    r.CreatedAt,
		})
	}
	return responses, nil
}

func (s *service) CancelReservation(reservationID, userID uint, role string) error {
	res, err := s.repo.FindByID(reservationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("reservation not found")
		}
		return err
	}

	// Drivers may only cancel their own reservations
	if role != "admin" && res.UserID != userID {
		return errors.New("forbidden")
	}

	if res.Status == "cancelled" {
		return errors.New("already cancelled")
	}

	return s.repo.Cancel(reservationID)
}

func (s *service) GetAllReservations() ([]dto.AdminReservationResponse, error) {
	list, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.AdminReservationResponse, 0, len(list))
	for _, r := range list {
		responses = append(responses, dto.AdminReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			User:         dto.UserSummary{ID: r.User.ID, Name: r.User.Name, Email: r.User.Email, Role: r.User.Role},
			Zone:         dto.ZoneSummary{ID: r.Zone.ID, Name: r.Zone.Name, Type: r.Zone.Type},
			CreatedAt:    r.CreatedAt,
			UpdatedAt:    r.UpdatedAt,
		})
	}
	return responses, nil
}
