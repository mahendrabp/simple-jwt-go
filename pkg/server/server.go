package server

import (
	"simple-jwt-go/pkg/config"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type Server struct {
	cfg    *config.Config
	router *echo.Echo
	db     *sqlx.DB
}

func New(cfg *config.Config, db *sqlx.DB) *Server {
	return &Server{
		cfg:    cfg,
		router: echo.New(),
		db:     db,
	}
}

func (s *Server) Run() error {
	s.middleware()
	s.handlers()

	return s.router.Start(s.cfg.Server.Addr)
}
