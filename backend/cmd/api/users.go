package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "www.github.com/kharljhon14/kwento-ko/db/sqlc"
	"www.github.com/kharljhon14/kwento-ko/internal/token"
)

type getUserURI struct {
	ID string `uri:"id" binding:"required"`
}

func (s Server) getUser(ctx *gin.Context) {
	var uri getUserURI

	err := ctx.ShouldBindUri(&uri)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ID, err := uuid.Parse(uri.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.store.GetUserByID(ctx, pgtype.UUID{Bytes: ID, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			notFoundResponse(ctx, fmt.Errorf("user with ID %s could not be found", uri.ID))
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
		var ve validator.ValidationErrors

		if errors.As(err, &ve) {
			apiErrors := make([]ApiError, len(ve))

			for i, fe := range ve {
				apiErrors[i] = ApiError{
					strings.ToLower(fe.Field()),
					tagErrorMessage(fe),
				}
			}
			failedValidationError(ctx, apiErrors)
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
