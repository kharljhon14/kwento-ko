package api

import (
	"github.com/gin-gonic/gin"
	db "www.github.com/kharljhon14/kwento-ko/db/sqlc"
)

type Server struct {
	router *gin.Engine
	store  db.Store
}

func NewServer(store db.Store) (*Server, error) {
	server := &Server{store: store}

	server.mountRouter()
	return server, nil
}

func (s *Server) mountRouter() {

	router := gin.Default()

	router.GET("/api/v1/health", s.healthCheckHandler)
	router.GET("/api/v1/auth/:provider", s.signInWithProviderHandler)
	router.GET("/api/v1/auth/:provider/callback", s.signInCallbackHandler)
	s.router = router
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
