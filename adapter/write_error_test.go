package adapter_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/hjwalt/flows/adapter"
	"github.com/stretchr/testify/assert"
)

func TestWriteError(t *testing.T) {
	cases := []struct {
		name string
		err  error
	}{
		{
			name: "write error string into message",
			err:  errors.New("test error"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert := assert.New(t)

			rr := httptest.NewRecorder()
			err := adapter.WriteError(rr, 400, c.err)

			assert.NoError(err)
			assert.Equal(400, rr.Code)
			assert.Equal("{\"message\":\""+c.err.Error()+"\"}", rr.Body.String())
		})
	}
}
