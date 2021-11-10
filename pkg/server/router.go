package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	authDelivery "simple-jwt-go/pkg/auth/delivery/rest"
	authRepository "simple-jwt-go/pkg/auth/repository"
	authUseCase "simple-jwt-go/pkg/auth/usecase"

	echoSwagger "github.com/swaggo/echo-swagger"
)

func (s *Server) middleware() {

	s.router.Pre(middleware.RemoveTrailingSlash())
	s.router.Use(middleware.Logger())
	s.router.Use(middleware.Recover())
	s.router.Use(middleware.CORS())
}

func (s *Server) handlers() {
	s.router.GET("/ping", func(c echo.Context) (err error) {
		return c.JSON(http.StatusOK, "ping")
	})

	authRepo := authRepository.NewPostgreRepository(s.db)
	authUC := authUseCase.New(s.cfg, authRepo, s.log)

	if s.cfg.Server.Debug {
		s.router.GET("/swagger/*", echoSwagger.WrapHandler)
	}
	api := s.router.Group("/api")

	authDelivery.Init(s.cfg, api, authUC, s.log)
}
