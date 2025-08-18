package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ApiError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func errorResponse(err error) gin.H {
	return gin.H{"errors": err.Error()}
}

func notFoundResponse(ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
		"errors": err.Error(),
	})
}

func failedValidationError(ctx *gin.Context, err []ApiError) {
	ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
		"errors": err,
	})
}

func tagErrorMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "this field is required"
	case "min":
		return fmt.Sprintf("minimum of %s characters", fieldError.Param())
	case "max":
		return fmt.Sprintf("maximum of %s characters", fieldError.Param())
	default:
		return ""
	}
}
