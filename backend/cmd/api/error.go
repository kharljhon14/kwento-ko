package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func notFoundResponse(ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
		"error": err.Error(),
	})
}
