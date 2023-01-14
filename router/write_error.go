package router

import (
	"net/http"

	"github.com/hjwalt/flows/format"
)

type ErrorResponse struct {
	Message string
}

func WriteError(w http.ResponseWriter, httpStatus int, err error) error {
	f := format.Json[ErrorResponse]()
	r := ErrorResponse{
		Message: err.Error(),
	}
	return WriteJson(w, httpStatus, r, f)
}
