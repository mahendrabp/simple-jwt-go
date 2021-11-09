package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	User struct {
		ID        uuid.UUID `json:"id" db:"id"`
		Email     string    `json:"email" db:"email" validate:"required,email" example:"test@email.test"`
		Password  string    `json:"password,omitempty" db:"password" validate:"omitempty,min=6,max=100" swaggerignore:"true"`
		UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
		CreatedAt time.Time `json:"created_at" db:"created_at"`
	}

	AuthUser struct {
		User         *User  `json:"user"`
		TokenType    string `json:"token_type" validate:"required" example:"Bearer"`
		ExpiresIn    int    `json:"expires_in" validate:"required" example:"360"`
		AccessToken  string `json:"access_token" validate:"required"`
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
)

func (u *User) Validate() error {
	validate := validator.New()

	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)

	return validate.Struct(u)
}

func (u *User) ValidatePassword() error {
	if u.Password == "" {
		return errors.New("password can not be empty")
	}

	return nil
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(u.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}

func (u *User) ComparePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}

	return nil
}
