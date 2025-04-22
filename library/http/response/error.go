package response

import (
	"fmt"
	"log"
	"net/http"

	"case-study-kredit-plus/library/types"

	"github.com/gin-gonic/gin"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/pkg/errors"
)

// FieldError represents error message for each field
//
//swagger:model
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorResponse represents error message
//
//swagger:model
type ErrorResponse struct {
	Code    string        `json:"code"`
	Status  string        `json:"Status"`
	Message string        `json:"Message"`
	Fields  []*FieldError `json:"fields"`
	Data    *DataError    `json:"Data"`
}

type DataError struct {
	Message string `json:"Message"`
	Status  int    `json:"Status"`
}

// MakeFieldError create field error object
func MakeFieldError(field string, message string) *FieldError {
	return &FieldError{
		Field:   field,
		Message: message,
	}
}

// Error writes error http response
func Error(c *gin.Context, data string, status int, err types.Error) {
	var errorCode string

	if status == 0 {
		status = http.StatusInternalServerError
	}

	switch status {
	case http.StatusUnauthorized:
		errorCode = "Unauthorized"
	case http.StatusNotFound:
		errorCode = "NotFound"
	case http.StatusBadRequest:
		errorCode = "BadRequest"
	case http.StatusUnprocessableEntity:
		errorCode = "ValidationError"
	case http.StatusInternalServerError:
		errorCode = "InternalServerError"
	case http.StatusNotImplemented:
		errorCode = "NotImplemented"
	}

	errorFields := []*FieldError{}

	switch err.Error.(type) {
	case validator.ValidationErrors:
		for _, err := range err.Error.(validator.ValidationErrors) {
			e := MakeFieldError(
				err.Field(),
				err.ActualTag())

			errorFields = append(errorFields, e)
		}

		data = "Unprocessable Entity"
		errorCode = "UnprocessableEntity"
		status = http.StatusUnprocessableEntity
	}

	c.JSON(status, ErrorResponse{
		Code:    errorCode,
		Status:  "Warning",
		Message: data,
		Fields:  errorFields,
		Data:    nil,
	})

	if err.Error != nil {
		log.Printf("INFO: %v\n", err.Error.Error())
		log.Printf("DETAIL [%s - %s]: %s\n", err.Path, err.Type, err.Message)
		type stackTracer interface {
			StackTrace() errors.StackTrace
		}

		var st errors.StackTrace
		if err, ok := err.Error.(stackTracer); ok {
			st = err.StackTrace()
			fmt.Printf("INFO: %+v\n", st[0])
		}
	}
}
