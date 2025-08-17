package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

type getBlogURI struct {
	ID string `uri:"id" binding:"required"`
}

func (s Server) getBlogHandler(ctx *gin.Context) {
	var uri getBlogURI

	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ID, err := uuid.Parse(uri.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	blog, err := s.store.GetBlogByID(ctx, pgtype.UUID{Bytes: ID, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			notFoundResponse(ctx, fmt.Errorf("blog with ID %s could not be found", ID))
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, envelope{"data": blog})

}
