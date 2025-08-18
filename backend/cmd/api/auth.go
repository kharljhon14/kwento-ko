package api

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/markbates/goth/gothic"
	db "www.github.com/kharljhon14/kwento-ko/db/sqlc"
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

	gotUser, err := s.store.GetUserByEmail(ctx, user.Email)
	if err == nil {
		duration := (24 * time.Hour) * 7
		token, err := s.tokenMaker.CreateToken(gotUser.ID, user.Email, duration)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.Header(authorizationHeaderKey, token)

		ctx.SetCookie("access_token", token, int(duration.Seconds()), "/", "localhost:8080", true, true)
		ctx.Redirect(http.StatusTemporaryRedirect, os.Getenv("CLIENT_URL"))
		return
	}

	newUser := db.CreateUserParams{
		Name:         utils.GenerateRandomName(),
		Email:        user.Email,
		GoogleID:     user.UserID,
		ProfileImage: pgtype.Text{String: user.AvatarURL, Valid: true},
	}

	userID, err := s.store.CreateUser(ctx, newUser)
	if err != nil {

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	duration := (24 * time.Hour) * 7
	token, err := s.tokenMaker.CreateToken(userID, user.Email, duration)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Header(authorizationHeaderKey, token)
	ctx.SetCookie("access_token", token, int(duration.Seconds()), "/", "localhost:8080", true, true)
	ctx.Redirect(http.StatusTemporaryRedirect, os.Getenv("CLIENT_URL"))
}
