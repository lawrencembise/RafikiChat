package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"rafikichat/internal/adapters/controllers"
	"rafikichat/internal/adapters/http"
	"rafikichat/internal/infrastructure/config"
)

// Server is the struct for server engine and configuration
type Server struct {
	engine *gin.Engine
	config config.Config
}

// NewServer is the function for creating new server
func NewServer(cfg config.Config) *Server {
	app := gin.Default()
	controllers.SetConfig(cfg)
	http.Routers(app)
	return &Server{
		engine: app,
		config: cfg,
	}
}

// Start is the function for starting server
func (s *Server) Start() error {
	return s.engine.Run(fmt.Sprintf(":%s", s.config.Port))
}