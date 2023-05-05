package router_test

import (
	"net/http/httptest"
	"testing"

	"github.com/hjwalt/flows/router"
	"github.com/hjwalt/flows/test_helper"
	"github.com/stretchr/testify/assert"
)

type CrappyResponseStruct struct {
	Message string
	Error   bool
}

func (crs CrappyResponseStruct) ShouldError() bool {
	return crs.Error
}

func TestWriteJson(t *testing.T) {
	cases := []struct {
		name     string
		status   int
		response CrappyResponseStruct
		err      error
	}{
		{
			name:     "write json",
			status:   200,
			response: CrappyResponseStruct{Message: "test", Error: false},
			err:      nil,
		},
		{
			name:     "format error",
			status:   200,
			response: CrappyResponseStruct{Message: "error", Error: true},
			err:      router.ErrorWriteJsonFormat,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert := assert.New(t)

			rr := httptest.NewRecorder()
			err := router.WriteJson(rr, c.status, c.response, test_helper.CrappyJson[CrappyResponseStruct]())

			if c.err == nil {
				assert.NoError(err)
				assert.Equal(c.status, rr.Code)
				assert.Equal("{\"Message\":\""+c.response.Message+"\",\"Error\":false}", rr.Body.String())
			} else {
				assert.ErrorIs(err, c.err)
			}
		})
	}
}
