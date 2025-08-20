package api

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func (s Server) healthCheckHandler(ctx *gin.Context) {
	log.Println(os.Getenv("CLIENT_URL"))
	ctx.JSON(http.StatusOK, nil)
}
