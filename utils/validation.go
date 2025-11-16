package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// SuccessResponse represents a standardized success response
type SuccessResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

// ValidationError represents validation errors for fields
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// RespondError responds with an error
func RespondError(c *gin.Context, statusCode int, message string, details interface{}) {
	c.JSON(statusCode, ErrorResponse{
		Status:  statusCode,
		Message: message,
		Details: details,
	})
}

// RespondSuccess responds with success
func RespondSuccess(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, SuccessResponse{
		Status: statusCode,
		Data:   data,
	})
}

// RespondValidationError responds with validation errors
func RespondValidationError(c *gin.Context, errors []ValidationError) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: "Validation failed",
		Details: errors,
	})
}

// ValidateStruct validates a struct and returns validation errors
func ValidateStruct(data interface{}) []ValidationError {
	validate := validator.New()
	err := validate.Struct(data)

	if err == nil {
		return nil
	}

	var validationErrors []ValidationError
	for _, err := range err.(validator.ValidationErrors) {
		validationErrors = append(validationErrors, ValidationError{
			Field:   err.Field(),
			Message: getValidationMessage(err),
		})
	}

	return validationErrors
}

// getValidationMessage returns a human-readable validation message
func getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", err.Field(), err.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", err.Field(), err.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", err.Field(), err.Param())
	default:
		return fmt.Sprintf("%s is invalid", err.Field())
	}
}

// RecoverPanic middleware to recover from panics
func RecoverPanic() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			RespondError(c, http.StatusInternalServerError, "Internal Server Error", err)
		} else {
			RespondError(c, http.StatusInternalServerError, "Internal Server Error", nil)
		}
	})
}
