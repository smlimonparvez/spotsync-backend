package parking_zone

import (
	"spotsync/internal/domain/parking_zone/dto"
)

// Service defines the business-logic contract for the parking_zone domain.
type Service interface {
	CreateZone(req *dto.CreateZoneRequest) (*dto.ZoneResponse, error)
	GetAllZones() ([]dto.ZoneResponse, error)
	GetZoneByID(id uint) (*dto.ZoneResponse, error)
}

type service struct {
	repo Repository
}

// NewService returns a Service backed by the given Repository.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateZone(req *dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	z := &ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}
	if err := s.repo.Create(z); err != nil {
		return nil, err
	}
	// Newly created zone has full capacity available
	return toZoneResponse(z, z.TotalCapacity), nil
}

func (s *service) GetAllZones() ([]dto.ZoneResponse, error) {
	zones, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ZoneResponse, 0, len(zones))
	for _, z := range zones {
		active, err := s.repo.CountActiveReservations(z.ID)
		if err != nil {
			return nil, err
		}
		available := z.TotalCapacity - int(active)
		if available < 0 {
			available = 0
		}
		responses = append(responses, *toZoneResponse(&z, available))
	}
	return responses, nil
}

func (s *service) GetZoneByID(id uint) (*dto.ZoneResponse, error) {
	z, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	active, err := s.repo.CountActiveReservations(z.ID)
	if err != nil {
		return nil, err
	}
	available := z.TotalCapacity - int(active)
	if available < 0 {
		available = 0
	}
	return toZoneResponse(z, available), nil
}

func toZoneResponse(z *ParkingZone, available int) *dto.ZoneResponse {
	return &dto.ZoneResponse{
		ID:             z.ID,
		Name:           z.Name,
		Type:           z.Type,
		TotalCapacity:  z.TotalCapacity,
		AvailableSpots: available,
		PricePerHour:   z.PricePerHour,
		CreatedAt:      z.CreatedAt,
	}
}
