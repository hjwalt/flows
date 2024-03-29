package adapter_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/hjwalt/flows/adapter"
	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/test_helper"
	"github.com/stretchr/testify/assert"
)

func TestRouteProducerBodyMapConversion(t *testing.T) {
	testcases := []struct {
		name   string
		input  flow.Message[[]byte, []byte]
		output flow.Message[[]byte, []byte]
		empty  bool
		err    string
	}{
		{
			name: "basic conversion",
			input: flow.Message[[]byte, []byte]{
				Key:   []byte("k"),
				Value: []byte("v"),
			},
			output: flow.Message[[]byte, []byte]{
				Key:   []byte("k-updated"),
				Value: []byte("v-updated"),
			},
			err:   "",
			empty: false,
		},
		{
			name: "empty result",
			input: flow.Message[[]byte, []byte]{
				Key:   []byte("empty"),
				Value: []byte("v"),
			},
			err:   "",
			empty: true,
		},
		{
			name: "error input conversion",
			input: flow.Message[[]byte, []byte]{
				Key:   []byte("error"),
				Value: []byte("error"),
			},
			err:   "error",
			empty: false,
		},
		{
			name: "error execute",
			input: flow.Message[[]byte, []byte]{
				Key:   []byte("stupid error"),
				Value: []byte("v"),
			},
			err:   "stupid error",
			empty: false,
		},
		{
			name: "error output conversion",
			input: flow.Message[[]byte, []byte]{
				Key:   []byte("output error"),
				Value: []byte("v"),
			},
			err:   "error",
			empty: false,
		},
	}
	crappyStringFormat := test_helper.CrappyStringFormat()
	source := func(ctx context.Context, m flow.Message[[]byte, string]) (*flow.Message[string, string], error) {
		if strings.ToLower(string(m.Key)) == "stupid error" {
			return &flow.Message[string, string]{}, errors.New(string(m.Key))
		}

		if strings.ToLower(string(m.Key)) == "output error" {
			return &flow.Message[string, string]{
				Key:   "error",
				Value: "error",
			}, nil
		}

		if strings.ToLower(string(m.Key)) == "empty" {
			return nil, nil
		}

		return &flow.Message[string, string]{
			Key:   string(m.Key) + "-updated",
			Value: m.Value + "-updated",
		}, nil
	}

	converted := adapter.RouteProduceBodyMapConvert(source, crappyStringFormat, crappyStringFormat, crappyStringFormat)

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			assert := assert.New(t)

			output, err := converted(context.Background(), testcase.input)
			if len(testcase.err) > 0 {
				assert.Nil(output)
				assert.Contains(err.Error(), testcase.err)
				return
			}

			if testcase.empty {
				assert.Nil(output)
				return
			}

			assert.NotNil(output)
			assert.Equal(testcase.output, *output)
		})
	}
}
