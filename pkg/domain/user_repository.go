package domain

import (
	"github.com/google/uuid"
)

type (
	PostgreUserRepository interface {
		GetByID(id uuid.UUID) (User, error)
		FindByEmail(email string) (User, error)
		Store(u *User) (*User, error)
		Update(u *User) (*User, error)
		Delete(id uuid.UUID) error
	}
)
