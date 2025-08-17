package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "www.github.com/kharljhon14/kwento-ko/db/sqlc"
	"www.github.com/kharljhon14/kwento-ko/internal/token"
)

type createBlogRequest struct {
	Title   string   `json:"title" binding:"required"`
	Content string   `json:"content" binding:"required"`
	Tags    []string `json:"tags" binding:"required"`
}

func (s Server) createBlogHandler(ctx *gin.Context) {
	var req createBlogRequest

	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	args := db.CreateBlogParams{
		Title:   req.Title,
		Content: req.Content,
		Author:  authPayload.ID,
	}

	newBlog, err := s.store.CreateBlog(ctx, args)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, envelope{"data": newBlog})
}
