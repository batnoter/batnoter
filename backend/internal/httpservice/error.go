package httpservice

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AppError represents the application error.
// Error of this type should be used in case when client needs to be informed with the failure details.
type AppError struct {
	code    string
	message string
	cause   error
}

// NewAppError creates and returns a new app error with provided code and message.
func NewAppError(code string, message string) error {
	return &AppError{
		code:    code,
		message: message,
	}
}

// NewAppErrorWithCause creates and returns a new app error with cause.
func NewAppErrorWithCause(code string, message string, cause error) error {
	return &AppError{
		code:    code,
		message: message,
		cause:   cause,
	}
}

func (a *AppError) Error() string {
	if a.cause != nil {
		return fmt.Sprintf("error code: %s, message: %s, cause: %v", a.code, a.message, a.cause)
	}
	return fmt.Sprintf("error code: %s, message: %s ", a.code, a.message)
}

// ErrorResponse represents the response payload of an app error.
type ErrorResponse struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

const (
	// ErrorCodeInvalidRequest error code for invalid http request.
	ErrorCodeInvalidRequest = "invalid_request"

	// ErrorCodeValidationFailed error code for validation fails on http request payload.
	ErrorCodeValidationFailed = "validation_failed"

	// ErrorCodeInternalServerError error code for internal server error.
	ErrorCodeInternalServerError = "internal_server_error"
)

func abortRequestWithError(c *gin.Context, err error) {
	var appErr *AppError
	errors.As(err, &appErr)
	if appErr != nil {
		logrus.WithField("error_code", appErr.code).WithField("error_message", appErr.message).Error("bad request")
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			Code:    appErr.code,
			Message: appErr.message,
		})
	} else {
		logrus.WithField("error_message", err.Error()).Error("request failed due to internal server error")
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
			Code:    ErrorCodeInternalServerError,
			Message: "something went wrong.",
		})
	}
}
