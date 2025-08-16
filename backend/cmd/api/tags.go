package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type createTagRequest struct {
	Name string `json:"name" binding:"required,max=60"`
}

func (s Server) createTagHandler(ctx *gin.Context) {
	var req createTagRequest

	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
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
