package rest

import (
	"fmt"
	"net/http"
	"simple-jwt-go/pkg/config"
	"simple-jwt-go/pkg/domain"
	"simple-jwt-go/pkg/middleware"
	"simple-jwt-go/pkg/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type handler struct {
	cfg         *config.Config
	userUseCase domain.UserUseCase
	log         utils.Logger
}

func newHandler(cfg *config.Config, uuc domain.UserUseCase, log utils.Logger) *handler {
	return &handler{cfg, uuc, log}
}

func Init(cfg *config.Config, e *echo.Group, uuc domain.UserUseCase, log utils.Logger) {
	h := newHandler(cfg, uuc, log)
	auth := middleware.Auth(cfg, uuc, log)

	authGroup := e.Group("/auth")
	authGroup.GET("/me", h.Me, auth)
	authGroup.POST("/login", h.Login)
	authGroup.POST("/register", h.Register)
	authGroup.POST("/refresh", h.Refresh)

}

func (h *handler) setCookies(c echo.Context, user *domain.AuthUser) {
	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    user.AccessToken,
		Path:     "/",
		MaxAge:   h.cfg.Auth.AccessToken.MaxAge,
		Secure:   h.cfg.Auth.AccessToken.Secure,
		HttpOnly: h.cfg.Auth.AccessToken.HttpOnly,
		SameSite: http.SameSiteStrictMode,
	})

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    user.RefreshToken,
		Path:     "/",
		MaxAge:   h.cfg.Auth.RefreshToken.MaxAge,
		Secure:   h.cfg.Auth.RefreshToken.Secure,
		HttpOnly: h.cfg.Auth.RefreshToken.HttpOnly,
		SameSite: http.SameSiteStrictMode,
	})
}

func (h *handler) Me(c echo.Context) error {
	fmt.Println(c.Get("user_id"))
	userID := c.Get("user_id").(uuid.UUID)

	res, err := h.userUseCase.GetByID(userID)
	if err != nil {
		h.log.ErrorFormat("auth.UseCase.GetByID: %v", err)
		return err
	}

	return c.JSON(http.StatusOK, res)
}

func (h *handler) Login(c echo.Context) error {
	u := new(domain.User)

	if err := c.Bind(u); err != nil {
		return echo.ErrBadRequest
	}

	user, err := h.userUseCase.Login(u)
	if err != nil {
		h.log.ErrorFormat("auth.UseCase.Login: %v", err)
		return err
	}

	h.setCookies(c, user)

	return c.JSON(http.StatusOK, user)
}

func (h *handler) Register(c echo.Context) error {
	u := new(domain.User)

	if err := c.Bind(u); err != nil {
		return echo.ErrBadRequest
	}

	createdUser, err := h.userUseCase.Store(u)
	if err != nil {
		h.log.ErrorFormat("auth.UseCase.Store: %v", err)
		return err
	}

	h.setCookies(c, createdUser)

	return c.JSON(http.StatusCreated, createdUser)
}

func (h *handler) Refresh(c echo.Context) error {

	mapToken := map[string]string{}

	if err := c.Bind(&mapToken); err != nil {
		return echo.ErrBadRequest
	}

	refreshToken := mapToken["refresh_token"]

	if refreshToken != "" {
		if err := utils.VerifyRefreshTokenInBody(c, h.cfg.Server.JwtRefreshSecret, h.log, refreshToken); err != nil {
			h.log.ErrorFormat("verifyRefreshToken: %v", err)
			return echo.ErrUnauthorized
		}
	} else {
		if err := utils.VerifyRefreshToken(c, h.cfg.Server.JwtRefreshSecret, h.log); err != nil {
			h.log.ErrorFormat("verifyRefreshToken: %v", err)
			return echo.ErrUnauthorized
		}
	}

	userID := c.Get("user_id").(uuid.UUID)

	u, err := h.userUseCase.GetByID(userID)
	if err != nil {
		h.log.ErrorFormat("auth.UseCase.GetByID: %v", err)
		return err
	}

	user, err := h.userUseCase.Auth(&u)
	if err != nil {
		h.log.ErrorFormat("auth.UseCase.Auth: %v", err)
		return err
	}

	h.setCookies(c, user)

	return c.JSON(http.StatusOK, user)
}
