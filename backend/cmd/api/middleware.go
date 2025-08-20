package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"www.github.com/kharljhon14/kwento-ko/internal/token"
)

const (
	authorizationHeaderKey  = "access_token"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.JWTMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		accessToken, err := ctx.Cookie(authorizationHeaderKey)
		if err != nil {
			err := errors.New("access token is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
