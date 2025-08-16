package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/markbates/goth/gothic"
	db "www.github.com/kharljhon14/kwento-ko/db/sqlc"
	"www.github.com/kharljhon14/kwento-ko/internal/token"
	"www.github.com/kharljhon14/kwento-ko/internal/utils"
)

func (s Server) signInWithProviderHandler(ctx *gin.Context) {
	provider := ctx.Param("provider")
	q := ctx.Request.URL.Query()

	q.Add("provider", provider)
	ctx.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

func (s Server) signInCallbackHandler(ctx *gin.Context) {
	provider := ctx.Param("provider")
	q := ctx.Request.URL.Query()
	q.Add("provider", provider)
	ctx.Request.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = s.store.GetUser(ctx, user.Email)
	if err == nil {
		duration := (24 * time.Hour) * 7
		token, err := s.tokenMaker.CreateToken(user.Email, duration)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, token)
		return
	}

	newUser := db.CreateUserParams{
		Name:         utils.GenerateRandomName(),
		Email:        user.Email,
		GoogleID:     user.UserID,
		ProfileImage: pgtype.Text{String: user.AvatarURL, Valid: true},
	}

	_, err = s.store.CreateUser(ctx, newUser)
	if err != nil {

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	duration := (24 * time.Hour) * 7
	token, err := s.tokenMaker.CreateToken(user.Email, duration)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, token)
}

type updateUserRequest struct {
	Name string `json:"name" binding:"required"`
}

func (s Server) updateUserHandler(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req updateUserRequest

	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	updatedName, err := s.store.UpdateUser(ctx, db.UpdateUserParams{Name: req.Name, Email: authPayload.Email})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, envelope{"data": updatedName})

}
