package rest

import (
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

}
