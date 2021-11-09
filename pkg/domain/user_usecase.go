package domain

import (
	"simple-jwt-go/pkg/utils"

	"github.com/google/uuid"
)

type (
	JWTUseCase interface {
		Auth(user *User) (*AuthUser, error)
		GetToken(id uuid.UUID, tokenID uuid.UUID) (uuid.UUID, error)
		DeleteToken(id uuid.UUID, tokenID uuid.UUID) error
		Logout(id uuid.UUID, tokenID *utils.TokenDetails) error
	}

	UserUseCase interface {
		GetByID(id uuid.UUID) (User, error)
		Login(user *User) (*AuthUser, error)
		Store(user *User) (*AuthUser, error)
		Update(user *User) (*User, error)
		Delete(id uuid.UUID) error
		JWTUseCase
	}
)
