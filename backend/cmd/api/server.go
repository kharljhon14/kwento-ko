package api

import (
	"os"

	"github.com/gin-gonic/gin"
	db "www.github.com/kharljhon14/kwento-ko/db/sqlc"
	"www.github.com/kharljhon14/kwento-ko/internal/token"
)

type Server struct {
	router     *gin.Engine
	store      db.Store
	tokenMaker *token.JWTMaker
}

func NewServer(store db.Store) (*Server, error) {

	tokenMaker, err := token.NewJWTMaker(os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		return nil, err
	}

	server := &Server{store: store, tokenMaker: tokenMaker}

	server.mountRouter()
	return server, nil
}

func (s *Server) mountRouter() {

	router := gin.Default()

	router.GET("/api/v1/health", s.healthCheckHandler)
	router.GET("/api/v1/auth/:provider", s.signInWithProviderHandler)
	router.GET("/api/v1/auth/:provider/callback", s.signInCallbackHandler)

	authRoutes := router.Group("/").Use(authMiddleware(*s.tokenMaker))

	authRoutes.GET("/api/v1/users/:id", s.getUser)
	authRoutes.PATCH("/api/v1/users", s.updateUserHandler)

	authRoutes.POST("/api/v1/tags", s.createTagHandler)
	authRoutes.GET("/api/v1/tags/:id", s.getTagHandler)
	authRoutes.GET("/api/v1/tags", s.getTagsHandler)
	authRoutes.PATCH("/api/v1/tags/:id", s.updateTagHandler)
	authRoutes.DELETE("/api/v1/tags/:id", s.deleteTagHandler)

	authRoutes.GET("/api/v1/blogs", s.getBlogsHandler)
	authRoutes.POST("/api/v1/blogs", s.createBlogHandler)
	authRoutes.GET("/api/v1/blogs/:id", s.getBlogHandler)
	authRoutes.PATCH("/api/v1/blogs/:id", s.updateBlogHandler)

	s.router = router
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

type envelope gin.H
