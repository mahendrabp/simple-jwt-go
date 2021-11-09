package rest

import (
	"net/http"
	"simple-jwt-go/pkg/config"
	"simple-jwt-go/pkg/domain"
	"simple-jwt-go/pkg/utils"

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

func Init(cfg *config.Config, e *echo.Group, uu domain.UserUseCase, log utils.Logger) {
	h := newHandler(cfg, uu, log)

	apiGroup := e.Group("/api")
	authGroup := apiGroup.Group("/auth")
	authGroup.POST("/login", h.Login)

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
