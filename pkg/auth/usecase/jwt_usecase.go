package usecase

import (
	"simple-jwt-go/pkg/domain"
	"simple-jwt-go/pkg/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const authPrefix = "auth"

func (u *usecase) Auth(user *domain.User) (*domain.AuthUser, error) {
	td, err := utils.GenerateToken(
		&utils.JWTConfig{
			JWTSecret:           u.cfg.Server.JwtSecret,
			JWTRefreshSecret:    u.cfg.Server.JwtRefreshSecret,
			AccessTokenExpires:  u.cfg.Auth.AccessToken.MaxAge,
			RefreshTokenExpires: u.cfg.Auth.RefreshToken.MaxAge,
		},
		user.ID,
	)

	if err != nil {
		u.log.ErrorFormat("generateToken: %v", err)
		return nil, echo.ErrInternalServerError
	}

	return &domain.AuthUser{
		User:         user,
		TokenType:    "Bearer",
		ExpiresIn:    u.cfg.Auth.AccessToken.MaxAge,
		AccessToken:  td.AccessToken,
		RefreshToken: td.RefreshToken,
	}, nil
}

func (u *usecase) GetToken(id uuid.UUID, tokenID uuid.UUID) (uuid.UUID, error) {
	panic("impelement me")
}

func (u *usecase) DeleteToken(id uuid.UUID, tokenID uuid.UUID) error {
	panic("impelement me")
}

func (u *usecase) Logout(id uuid.UUID, td *utils.TokenDetails) error {
	panic("impelement me")
}

func (u *usecase) LogoutAll(id uuid.UUID) error {
	panic("impelement me")
}
