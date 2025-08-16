package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"www.github.com/kharljhon14/kwento-ko/internal/token"
)

func (s Server) healthCheckHandler(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	fmt.Println(authPayload)
	ctx.JSON(http.StatusOK, authPayload)
}
