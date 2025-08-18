package api

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	db "www.github.com/kharljhon14/kwento-ko/db/sqlc"
	"www.github.com/kharljhon14/kwento-ko/internal/token"
)

func (s Server) getUser(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	user, err := s.store.GetUserByID(ctx, authPayload.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			notFoundResponse(ctx, fmt.Errorf("user with ID %s could not be found", authPayload.ID))
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, envelope{"data": user})
}

type updateUserRequest struct {
	Name string `json:"name" binding:"required,min=5,max=30"`
}

func (s Server) updateUserHandler(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req updateUserRequest

	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		if errors.Is(err, io.EOF) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(errors.New("body must not be empty")))
			return
		}

		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			failedValidationError(ctx, ve)
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	updatedName, err := s.store.UpdateUser(ctx, db.UpdateUserParams{Name: req.Name, Email: authPayload.Email})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, envelope{"data": updatedName})

}
