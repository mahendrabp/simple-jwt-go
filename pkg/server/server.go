package server

import (
	"simple-jwt-go/pkg/config"
	"simple-jwt-go/pkg/utils"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	cfg    *config.Config
	router *echo.Echo
	db     *sqlx.DB
	log    utils.Logger
}

func New(cfg *config.Config, db *sqlx.DB, log utils.Logger) *Server {
	return &Server{
		cfg:    cfg,
		router: echo.New(),
		db:     db,
		log:    log,
	}
}

func (s *Server) Run() error {
	s.router.Use(middleware.Logger())
	s.middleware()
	s.handlers()

	return s.router.Start(s.cfg.Server.Addr)
}
