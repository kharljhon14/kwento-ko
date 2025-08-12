package api

import "github.com/gin-gonic/gin"

type Server struct {
	router *gin.Engine
}

func NewServer() (*Server, error) {
	server := &Server{}

	server.mountRouter()
	return server, nil
}

func (s *Server) mountRouter() {

	router := gin.Default()

	router.GET("/api/v1/health", s.healthCheckHandler)
	s.router = router
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
