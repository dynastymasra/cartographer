package config

import (
	"fmt"
	"net/http"

	"github.com/dynastymasra/cookbook"
)

type ServiceError struct {
	code    int
	key     string
	message string
}

func NewError(code int, key, message string) *ServiceError {
	return &ServiceError{
		code:    code,
		key:     key,
		message: message,
	}
}

func (s *ServiceError) Code() int {
	return s.code
}

func (s *ServiceError) Key() string {
	return s.key
}

func (s *ServiceError) Error() string {
	return s.message
}

func ParseToJSON(err *ServiceError, w http.ResponseWriter, requestID string) {
	if err.Code() >= 500 {
		w.WriteHeader(err.Code())
		fmt.Fprint(w, cookbook.ErrorResponse(err.message, requestID).Stringify())
	}

	if err.Code() >= 400 && err.Code() < 500 {
		w.WriteHeader(err.Code())
		fmt.Fprint(w, cookbook.FailResponse(&cookbook.JSON{
			err.Key(): err.Error(),
		}, requestID).Stringify())
	}
}
