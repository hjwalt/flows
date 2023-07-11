package router

import (
	"net/http"

	"github.com/hjwalt/runway/format"
)

func WriteError(w http.ResponseWriter, httpStatus int, err error) error {
	f := format.Json[BasicResponse]()
	r := BasicResponse{
		Message: err.Error(),
	}
	return WriteJson(w, httpStatus, r, f)
}
