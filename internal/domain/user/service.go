package user

import (
	"errors"

	"spotsync/internal/auth"
	"spotsync/internal/domain/user/dto"

	"gorm.io/gorm"
)

// Service defines the business-logic contract for the user domain.
type Service interface {
	Register(req *dto.RegisterRequest) (*dto.UserResponse, error)
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
}

type service struct {
	repo   Repository
	jwtSvc *auth.JWTService
}

// NewService returns a Service with the given repository and JWT service.
func NewService(repo Repository, jwtSvc *auth.JWTService) Service {
	return &service{repo: repo, jwtSvc: jwtSvc}
}

func (s *service) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
	// Reject duplicate emails early
	if _, err := s.repo.FindByEmail(req.Email); err == nil {
		return nil, errors.New("email already registered")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	role := req.Role
	if role == "" {
		role = "driver"
	}

	u := &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     role,
	}

	if err := u.HashPassword(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(u); err != nil {
		return nil, err
	}

	return toUserResponse(u), nil
}

func (s *service) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	u, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !u.CheckPassword(req.Password) {
		return nil, errors.New("invalid email or password")
	}

	token, err := s.jwtSvc.GenerateToken(u.ID, u.Role)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User:  *toUserResponse(u),
	}, nil
}

// toUserResponse maps a User entity to the safe DTO (no password).
func toUserResponse(u *User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
