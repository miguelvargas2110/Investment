package errors

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AppError represents a standard API error response
// @AppError
type AppError struct {
	// HTTP status code
	// Required: true
	// Example: 400
	Code int `json:"code"`

	// Descriptive error message
	// Required: true
	// Example: "Invalid request parameters"
	Message string `json:"message"`

	// Original error (not exposed in API response)
	// swagger:ignore
	Err error `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new AppError instance
// @Summary Creates a new application error
// @Description Creates a structured error for API responses
// @Return *AppError
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// ErrorHandler is a global error handling middleware
// @Summary Global error handler
// @Description Consistently captures and formats errors for the entire API
// @Accept json
// @Produce json
// @Param c body gin.Context true "Gin context"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} AppError
// @Failure 401 {object} AppError
// @Failure 404 {object} AppError
// @Failure 500 {object} AppError
// @Router / [get]
func ErrorHandler(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		err := c.Errors.Last()

		var appErr *AppError
		if errors.As(err, &appErr) {
			c.JSON(appErr.Code, appErr)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Internal Server Error",
			})
		}
	}
}
