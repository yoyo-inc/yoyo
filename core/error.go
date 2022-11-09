package core

import (
	"reflect"
	"strings"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/go-playground/validator/v10"
	valid "github.com/yoyo-inc/yoyo/common/validator"
)

// ParameterError represents parameter error
type ParameterError struct {
	Message string
}

func (e ParameterError) Error() string {
	return e.Message
}

// NewParameterError returns parameter error
func NewParameterError(err interface{}) ParameterError {
	var message string
	t := reflect.TypeOf(err)
	switch t.Kind() {
	case reflect.String:
		message = err.(string)
	case reflect.Slice:
		errMessages := slice.Map(err.(validator.ValidationErrors), func(_ int, item validator.FieldError) string {
			return item.Translate(valid.Trans)
		})
		message = strings.Join(errMessages, ";")
	}
	return ParameterError{
		Message: message,
	}
}

// BusinessError represents bussiness error
type BusinessError struct {
	Code    string
	Message string
}

func (e BusinessError) Error() string {
	return e.Message
}

// NewBusinessError returns business error
func NewBusinessError(code string, message string) BusinessError {
	return BusinessError{
		Code:    code,
		Message: message,
	}
}
