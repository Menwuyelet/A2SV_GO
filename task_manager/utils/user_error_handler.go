package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// handleValidationError parses the error and returns specific messages
func HandleValidationError(c *gin.Context, err error) {
	var ve validator.ValidationErrors

	// Check if the error is a ValidationErrors type
	if errors.As(err, &ve) {
		errorMessages := make(map[string]string)

		for _, e := range ve {
			// Customize messages based on the tag and field
			errorMessages[e.Field()] = getCustomMessage(e)
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": errorMessages,
		})
	} else {
		// Handle non-validation errors (e.g., malformed JSON)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
	}
}

// getCustomMessage returns a user-friendly string based on the validation tag
func getCustomMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "Value is too short (min " + fe.Param() + " characters)"
	case "max":
		return "Value is too long (max " + fe.Param() + " characters)"
	case "email":
		return "Must be a valid email address"
	default:
		return "Invalid value"
	}
}
