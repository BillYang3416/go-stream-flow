package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type errorResponse struct {
	Message   string            `json:"message" example:"message"`
	FieldErrs map[string]string `json:"fieldErrs,omitempty"`
}

func sendErrorResponse(c *gin.Context, code int, msg string) {
	c.JSON(code, errorResponse{Message: msg})
}

func sendValidationErrorResponse(c *gin.Context, validationErrs validator.ValidationErrors) {
	errMessages := make(map[string]string)
	for _, err := range validationErrs {
		switch err.Tag() {
		case "required":
			errMessages[err.Field()] = fmt.Sprintf("%s is required", err.Field())
		case "min":
			errMessages[err.Field()] = fmt.Sprintf("%s must be greater than %s", err.Field(), err.Param())
		case "max":
			errMessages[err.Field()] = fmt.Sprintf("%s must be less than %s", err.Field(), err.Param())
		case "datetime":
			errMessages[err.Field()] = fmt.Sprintf("%s must be in the format of %s", err.Field(), err.Param())
		case "url":
			errMessages[err.Field()] = fmt.Sprintf("%s must be a valid url", err.Field())
		case "validategamestatus":
			errMessages[err.Field()] = fmt.Sprintf("%s is not a valid game status", err.Field())
		case "validategametype":
			errMessages[err.Field()] = fmt.Sprintf("%s is not a valid game type", err.Field())
		default:
			errMessages[err.Field()] = fmt.Sprintf("%s is not valid", err.Field())
		}
	}
	c.JSON(http.StatusBadRequest, errorResponse{Message: "validation error", FieldErrs: errMessages})
}
