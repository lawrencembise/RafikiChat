package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"rafikichat/internal/adapters/http"
	"rafikichat/internal/infrastructure/config"
	"rafikichat/internal/adapters/controllers"
)

type Server struct {
	engine *gin.Engine
	config config.Config
}

func NewServer(cfg config.Config) *Server {
	app := gin.Default()
	controllers.SetConfig(cfg)
	http.Routers(app)
	return &Server{
		engine: app,
		config: cfg,
	}
}

func (s *Server) Start() error {
	return s.engine.Run(fmt.Sprintf(":%s", s.config.Port))
}