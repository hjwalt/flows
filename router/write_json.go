package router

import (
	"net/http"

	"github.com/hjwalt/flows/format"
)

func WriteJson[R any](w http.ResponseWriter, httpStatus int, r R, f format.Format[R]) error {

	jsonBody, jsonErr := f.ToJson(r)
	if jsonErr != nil {
		return jsonErr
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	// hidden response is length of bytes written
	_, writeErr := w.Write(jsonBody)
	if writeErr != nil {
		return writeErr
	}

	return nil
}
