package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	db "www.github.com/kharljhon14/kwento-ko/db/sqlc"
	"www.github.com/kharljhon14/kwento-ko/internal/filter"
)

type createTagRequest struct {
	Name string `json:"name" binding:"required,max=30"`
}

func (s Server) createTagHandler(ctx *gin.Context) {
	var req createTagRequest

	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			failedValidationError(ctx, ve)
			return
		}

		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tag, err := s.store.CreateTag(ctx, req.Name)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, envelope{"data": tag})
}

type getTagURI struct {
	ID string `uri:"id" binding:"required"`
}

func (s Server) getTagHandler(ctx *gin.Context) {
	var uri getTagURI

	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ID, err := uuid.Parse(uri.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tag, err := s.store.GetTag(ctx, pgtype.UUID{Bytes: ID, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, envelope{"data": tag})
}

type getTagsQuery struct {
	Page     int32  `form:"page" binding:"omitempty,min=1"`
	PageSize int32  `form:"page_size" binding:"omitempty,min=5,max=20"`
	Sort     string `form:"sort" binding:"omitempty"`
}

func (s Server) getTagsHandler(ctx *gin.Context) {
	var query getTagsQuery

	if err := ctx.ShouldBindQuery(&query); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			failedValidationError(ctx, ve)
			return
		}

		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if query.Page == 0 {
		query.Page = 1
	}

	if query.PageSize < 5 {
		query.PageSize = 5
	}

	if query.Sort == "" {
		query.Sort = "created_at"
	}

	filters := filter.Filter{
		Page:         query.Page,
		PageSize:     query.PageSize,
		Sort:         query.Sort,
		SortSafeList: []string{"title", "-title", "created_at", "-created_at"},
	}

	tagsCount, err := s.store.GetTagsCount(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	getTagsArgs := db.GetTagsParams{
		Offset: filters.Offset(),
		Limit:  filters.Limit(),
	}

	tags, err := s.store.GetTags(ctx, getTagsArgs)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	metadata := filter.CalaculateMetadata(int(tagsCount), int(filters.Page), int(filters.PageSize))

	ctx.JSON(http.StatusOK, envelope{"data": tags, "metadata": metadata})
}

type updateTagRequest struct {
	Name string `json:"name" binding:"required,min=1,max=30"`
}

func (s Server) updateTagHandler(ctx *gin.Context) {
	var req updateTagRequest

	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			failedValidationError(ctx, ve)
			return
		}

		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var uri getTagURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ID, err := uuid.Parse(uri.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	updatedTag, err := s.store.UpdateTag(ctx,
		db.UpdateTagParams{
			Name: req.Name,
			ID:   pgtype.UUID{Bytes: ID, Valid: true},
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, envelope{"data": updatedTag})
}

func (s Server) deleteTagHandler(ctx *gin.Context) {
	var uri getTagURI

	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ID, err := uuid.Parse(uri.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = s.store.GetTag(ctx, pgtype.UUID{Bytes: ID, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			notFoundResponse(ctx, fmt.Errorf("tag with ID %s could not be found", uri.ID))
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = s.store.DeleteTag(ctx, pgtype.UUID{Bytes: ID, Valid: true})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
