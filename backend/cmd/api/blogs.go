package api

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "www.github.com/kharljhon14/kwento-ko/db/sqlc"
	"www.github.com/kharljhon14/kwento-ko/internal/filter"
	"www.github.com/kharljhon14/kwento-ko/internal/token"
)

type createBlogRequest struct {
	Title   string   `json:"title" binding:"required,min=1,max=30"`
	Content string   `json:"content" binding:"required,min=1"`
	Tags    []string `json:"tags" binding:"required,min=1"`
}

func (s Server) createBlogHandler(ctx *gin.Context) {
	var req createBlogRequest

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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	args := db.CreateBlogParams{
		Title:   req.Title,
		Content: req.Content,
		Author:  authPayload.ID,
	}

	if len(req.Tags) > 5 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(errors.New("tags must be max of 5")))
		return
	}

	var tagIDs []pgtype.UUID
	for _, id := range req.Tags {
		ID, err := uuid.Parse(id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		tagIDs = append(tagIDs, pgtype.UUID{Bytes: ID, Valid: true})
	}

	newBlog, err := s.store.CreateBlog(ctx, args)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = s.store.AddBlogTags(ctx,
		db.AddBlogTagsParams{
			BlogID:  newBlog.ID,
			Column2: tagIDs,
		},
	)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	tags, err := s.store.GetBlogTags(ctx, newBlog.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	env := envelope{
		"data": map[string]any{
			"id":         newBlog.ID,
			"author":     newBlog.Author,
			"title":      newBlog.Title,
			"content":    newBlog.Content,
			"tags":       tags,
			"created_at": newBlog.CreatedAt,
			"version":    newBlog.Version,
		},
	}

	ctx.JSON(http.StatusCreated, envelope{"data": env})
}

type getBlogURI struct {
	ID string `uri:"id" binding:"required"`
}

func (s Server) getBlogHandler(ctx *gin.Context) {
	var uri getBlogURI

	if err := ctx.ShouldBindUri(&uri); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			failedValidationError(ctx, ve)
			return
		}
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

	tags, err := s.store.GetBlogTags(ctx, pgtype.UUID{Bytes: ID, Valid: true})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	env := envelope{
		"data": map[string]any{
			"id":         blog.ID,
			"author":     blog.Name,
			"title":      blog.Title,
			"content":    blog.Content,
			"tags":       tags,
			"created_at": blog.CreatedAt,
			"version":    blog.Version,
		},
	}

	ctx.JSON(http.StatusOK, envelope{"data": env})
}

type getGlogsQuery struct {
	Page     int32  `form:"page" binding:"omitempty,min=1"`
	PageSize int32  `form:"page_size" binding:"omitempty,min=5,max=20"`
	Sort     string `form:"sort" binding:"omitempty"`
}

func (s Server) getBlogsHandler(ctx *gin.Context) {
	var query getGlogsQuery

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

	blogCount, err := s.store.GetBlogCount(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	args := db.GetBlogsParams{
		Limit:  filters.Limit(),
		Offset: filters.Offset(),
	}

	blogs, err := s.store.GetBlogs(ctx, args)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	metadata := filter.CalaculateMetadata(int(blogCount), int(filters.Page), int(filters.PageSize))

	ctx.JSON(http.StatusOK, envelope{"data": blogs, "metadata": metadata})
}

type updateBlogRequest struct {
	Title   *string  `json:"title" binding:"omitempty,min=1,max=30"`
	Content *string  `json:"content" binding:"omitempty,min=1"`
	Tags    []string `json:"tags" binding:"omitempty,min=1,max=5"`
}

func (s Server) updateBlogHandler(ctx *gin.Context) {
	var uri getBlogURI

	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req updateBlogRequest
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.Title == nil && req.Content == nil && len(req.Tags) == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(errors.New("body must not be empty")))
		return
	}

	var tagIDs []pgtype.UUID
	for _, id := range req.Tags {
		ID, err := uuid.Parse(id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		tagIDs = append(tagIDs, pgtype.UUID{Bytes: ID, Valid: true})
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

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

	if blog.AuthorID != authPayload.ID {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("unauthorized to update this blog"))
		return
	}

	args := db.UpdateBlogParams{
		Title:   blog.Title,
		Content: blog.Content,
		ID:      blog.ID,
	}

	if req.Title != nil {
		args.Title = *req.Title
	}

	if req.Content != nil {
		args.Content = *req.Content
	}

	updatedBlog, err := s.store.UpdateBlog(ctx, args)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if req.Tags != nil {
		err := s.store.RemoveBlogTags(ctx, blog.ID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		err = s.store.AddBlogTags(ctx, db.AddBlogTagsParams{
			BlogID:  blog.ID,
			Column2: tagIDs,
		})
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	tags, err := s.store.GetBlogTags(ctx, blog.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	env := envelope{
		"data": map[string]any{
			"id":         updatedBlog.ID,
			"author":     updatedBlog.Author,
			"title":      updatedBlog.Title,
			"content":    updatedBlog.Content,
			"tags":       tags,
			"created_at": updatedBlog.CreatedAt,
			"version":    updatedBlog.Version,
		},
	}

	ctx.JSON(http.StatusOK, env)
}

func (s Server) deleteBlogHandler(ctx *gin.Context) {
	var uri getBlogURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	ID, err := uuid.Parse(uri.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	blog, err := s.store.GetBlogByID(ctx, pgtype.UUID{Bytes: ID, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrConnDone) {
			notFoundResponse(ctx, fmt.Errorf("blog with ID %s could not be found", ID))
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if blog.AuthorID != authPayload.ID {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("unauthorized to update this blog"))
		return
	}

	err = s.store.RemoveBlogTags(ctx, blog.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = s.store.DeleteBlog(ctx, db.DeleteBlogParams{ID: blog.ID, Author: blog.AuthorID})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
