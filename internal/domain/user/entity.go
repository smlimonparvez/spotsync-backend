package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User is the GORM model that maps to the "users" table.
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"not null"                 json:"name"`
	Email     string    `gorm:"uniqueIndex;not null"     json:"email"`
	Password  string    `gorm:"not null"                 json:"-"` // never serialised
	Role      string    `gorm:"not null;default:'driver'"json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// HashPassword replaces the plain-text password with a bcrypt hash (cost 12).
func (u *User) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

// CheckPassword compares a plain-text candidate against the stored bcrypt hash.
func (u *User) CheckPassword(plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain)) == nil
}
