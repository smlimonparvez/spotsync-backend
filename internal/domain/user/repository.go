package user

import (
	"gorm.io/gorm"
)

// Repository defines the data-access contract for the user domain.
type Repository interface {
	Create(u *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository returns a Repository backed by the given *gorm.DB.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(u *User) error {
	return r.db.Create(u).Error
}

func (r *repository) FindByEmail(email string) (*User, error) {
	var u User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) FindByID(id uint) (*User, error) {
	var u User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
