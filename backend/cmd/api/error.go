package api

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

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

func failedValidationError(ctx *gin.Context, validationErrors validator.ValidationErrors) {
	apiErrors := make([]ApiError, len(validationErrors))

	for i, fe := range validationErrors {
		apiErrors[i] = ApiError{
			strings.ToLower(fe.Field()),
			tagErrorMessage(fe),
		}
	}

	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"errors": apiErrors,
	})
}

func tagErrorMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "this field is required"
	case "min":
		if fieldError.Kind() == reflect.Slice {
			return fmt.Sprintf("minimum of %s item(s)", fieldError.Param())
		}
		return fmt.Sprintf("minimum of %s characters", fieldError.Param())

	case "max":
		if fieldError.Kind() == reflect.Slice {
			return fmt.Sprintf("maximum of %s item(s)", fieldError.Param())
		}
		return fmt.Sprintf("maximum of %s characters", fieldError.Param())
	default:
		return ""
	}
}
