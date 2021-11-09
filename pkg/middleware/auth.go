package middleware

import (
	"simple-jwt-go/pkg/config"
	"simple-jwt-go/pkg/domain"
	"simple-jwt-go/pkg/utils"
	"strings"

	"github.com/labstack/echo/v4"
)

func Auth(cfg *config.Config, userUseCase domain.UserUseCase, log utils.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := verifyAccessToken(cfg, c, log); err != nil {
				log.ErrorFormat("verifyAccessToken: %v", err)
				return echo.ErrUnauthorized
			}

			if err := utils.VerifyRefreshToken(c, cfg.Server.JwtRefreshSecret, log); err != nil {
				log.ErrorFormat("verifyRefreshToken: %v", err)
				return echo.ErrUnauthorized
			}

			return next(c)
		}
	}
}

func verifyAccessToken(cfg *config.Config, c echo.Context, log utils.Logger) error {
	tokenName := "access"
	bearerToken := c.Request().Header.Get("Authorization")
	if bearerToken != "" {
		arr := strings.Split(bearerToken, " ")
		if len(arr) == 2 && arr[0] == "Bearer" {
			token := arr[1]

			if err := utils.ValidateToken(
				c,
				tokenName,
				token,
				cfg.Server.JwtSecret,
			); err != nil {
				log.ErrorFormat("validateToken: %v", err)
				return err
			}

			return nil
		}
	}

	accessCookie, err := c.Cookie("access_token")
	if err != nil {
		log.ErrorFormat("c.Cookie: %v", err)
		return err
	}

	if err = utils.ValidateToken(
		c,
		tokenName,
		accessCookie.Value,
		cfg.Server.JwtSecret,
	); err != nil {
		log.ErrorFormat("validateToken: %v", err)
		return err
	}

	return nil
}
