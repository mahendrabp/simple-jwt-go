package usecase

import (
	"net/http"
	"simple-jwt-go/pkg/config"
	"simple-jwt-go/pkg/domain"
	"simple-jwt-go/pkg/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type usecase struct {
	cfg          *config.Config
	pgRepository domain.PostgreUserRepository
	log          utils.Logger
}

const cacheDuration = 3600

func New(cfg *config.Config, pg domain.PostgreUserRepository, log utils.Logger) domain.UserUseCase {
	return &usecase{
		cfg:          cfg,
		pgRepository: pg,
		log:          log,
	}
}

func (u *usecase) GetByID(id uuid.UUID) (domain.User, error) {
	res, err := u.pgRepository.GetByID(id)
	if err != nil {
		u.log.ErrorFormat("auth.pgRepository.GetByID: %v", err)
		return res, err
	}

	res.RemovePassword()

	return res, nil
}

func (u *usecase) Login(user *domain.User) (*domain.AuthUser, error) {
	if err := user.Validate(); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := user.ValidatePassword(); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res, err := u.pgRepository.FindByEmail(user.Email)
	if err != nil {
		u.log.ErrorFormat("auth.pgRepository.FindByEmail: %v", err)
		return nil, err
	}

	if err = res.ComparePassword(user.Password); err != nil {
		return nil, echo.ErrUnauthorized
	}

	res.RemovePassword()

	return u.Auth(&res)
}

func (u *usecase) Store(user *domain.User) (*domain.AuthUser, error) {
	if err := user.Validate(); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := user.ValidatePassword(); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := user.HashPassword(); err != nil {
		return nil, echo.ErrInternalServerError
	}

	res, err := u.pgRepository.Store(user)
	if err != nil {
		u.log.ErrorFormat("auth.pgRepository.Store: %v", err)
		return nil, err
	}

	res.RemovePassword()

	return u.Auth(res)
}

func (u *usecase) Update(user *domain.User) (*domain.User, error) {
	if err := user.Validate(); err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if user.Password != "" {
		if err := user.HashPassword(); err != nil {
			return nil, echo.ErrBadRequest
		}
	}

	res, err := u.pgRepository.Update(user)
	if err != nil {
		u.log.ErrorFormat("auth.pgRepository.Update: %v", err)
		return nil, err
	}

	res.RemovePassword()

	return res, nil
}

func (u *usecase) Delete(id uuid.UUID) error {
	if err := u.pgRepository.Delete(id); err != nil {
		u.log.ErrorFormat("auth.pgRepository.Delete: %v", err)
		return err
	}

	return nil
}
