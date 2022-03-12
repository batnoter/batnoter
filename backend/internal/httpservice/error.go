package httpservice

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AppError struct {
	code    string
	message string
	cause   error
}

func NewAppError(code string, message string) error {
	return &AppError{
		code:    code,
		message: message,
	}
}

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

type ErrorResponse struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

const (
	ErrorCodeInvalidRequest      = "invalid_request"
	ErrorCodeValidationFailed    = "validation_failed"
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
			Message: "something went wrong. contact support",
		})
	}
}
