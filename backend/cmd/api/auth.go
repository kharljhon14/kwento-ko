package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
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

	fmt.Println(user)

	ctx.Redirect(http.StatusTemporaryRedirect, "/success")
}
