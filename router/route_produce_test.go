package router_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/router"
	"github.com/hjwalt/flows/test_helper"
	"github.com/stretchr/testify/assert"
)

func TestRouteProducer(t *testing.T) {
	producer := test_helper.NewChannelProducer()

	routerProducer := router.NewRouteProducer(
		router.WithRouteProducerRuntime(producer),
		router.WithRouteBodyMap(func(ctx context.Context, req message.Message[message.Bytes, message.Bytes]) (*message.Message[message.Bytes, message.Bytes], error) {
			return &req, nil
		}),
	)

	cases := []struct {
		name       string
		method     string
		url        string
		body       io.Reader
		assertions func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		{
			name:   "get",
			method: "GET",
			url:    "/test",
			body:   nil,
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert := assert.New(t)

				assert.Equal(200, rr.Code)
				assert.Equal(1, len(producer.Messages))

				m := <-producer.Messages

				assert.Equal("", string(m.Value))
			},
		},
		{
			name:   "post",
			method: "POST",
			url:    "/test",
			body:   bytes.NewBuffer([]byte("test")),
			assertions: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert := assert.New(t)

				assert.Equal(200, rr.Code)
				assert.Equal(1, len(producer.Messages))

				m := <-producer.Messages

				assert.Equal("test", string(m.Value))
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			req, _ := http.NewRequest(c.method, c.url, c.body)
			rr := httptest.NewRecorder()
			routerProducer.Handle(rr, req)

			c.assertions(t, rr)
		})
	}
}
