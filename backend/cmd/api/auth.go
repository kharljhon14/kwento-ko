package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/markbates/goth/gothic"
	db "www.github.com/kharljhon14/kwento-ko/db/sqlc"
	utils "www.github.com/kharljhon14/kwento-ko/internal"
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
		ctx.AbortWithError(http.StatusInternalServerError, err)
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

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "users_email_key":
				ctx.JSON(http.StatusOK, user.IDToken)
				return
			}
		} else {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	ctx.JSON(http.StatusOK, user.IDToken)
}

type updateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email:" binding:"required"`
}

func (s Server) updateUserHandler(ctx *gin.Context) {
	var req updateUserRequest

	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
	}

	updatedName, err := s.store.UpdateUser(ctx, db.UpdateUserParams{Name: req.Name, Email: req.Email})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, updatedName)

}
