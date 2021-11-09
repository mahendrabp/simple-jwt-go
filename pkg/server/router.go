package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) middleware() {
	s.router.Pre(middleware.RemoveTrailingSlash())
	s.router.Use(middleware.CORS())
}

func (s *Server) handlers() {
	s.router.GET("/ping", func(c echo.Context) (err error) {
		return c.JSON(http.StatusOK, "ping")
	})
}
