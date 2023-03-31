package stateful_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/mock"
	"github.com/hjwalt/flows/stateful"
	"github.com/stretchr/testify/assert"
)

func TestConvertOneToOne(t *testing.T) {
	testcases := []struct {
		name        string
		input       message.Message[[]byte, []byte]
		inputState  []byte
		output      message.Message[[]byte, []byte]
		outputState []byte
		empty       bool
		err         string
	}{
		{
			name: "basic conversion",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("k"),
				Value: []byte("v"),
			},
			output: message.Message[[]byte, []byte]{
				Key:   []byte("k-updated"),
				Value: []byte("v-updated"),
			},
			inputState:  []byte("state"),
			outputState: []byte("statev"),
			err:         "",
			empty:       false,
		},
		{
			name: "empty result",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("empty"),
				Value: []byte("v"),
			},
			inputState:  []byte("state"),
			outputState: []byte("state"),
			err:         "",
			empty:       true,
		},
		{
			name: "error input conversion",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("error"),
				Value: []byte("error"),
			},
			inputState:  []byte("state"),
			outputState: []byte("state"),
			err:         "error",
			empty:       false,
		},
		{
			name: "error input state conversion",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("k"),
				Value: []byte("v"),
			},
			inputState:  []byte("error"),
			outputState: []byte("error"),
			err:         "error",
			empty:       false,
		},
		{
			name: "error execute",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("stupid error"),
				Value: []byte("v"),
			},
			inputState:  []byte("state"),
			outputState: []byte("state"),
			err:         "stupid error",
			empty:       false,
		},
		{
			name: "error output conversion",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("output error"),
				Value: []byte("v"),
			},
			inputState:  []byte("state"),
			outputState: []byte("state"),
			err:         "error",
			empty:       false,
		},
		{
			name: "error output state conversion",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("k"),
				Value: []byte("r"),
			},
			inputState:  []byte("erro"),
			outputState: []byte("erro"),
			err:         "error",
			empty:       false,
		},
	}

	crappyStringFormat := mock.CrappyStringFormat()
	oneToOne := func(ctx context.Context, m message.Message[string, string], ss stateful.SingleState[string]) (*message.Message[string, string], stateful.SingleState[string], error) {

		if strings.ToLower(m.Key) == "stupid error" {
			return &message.Message[string, string]{},
				ss,
				errors.New(m.Key)
		}

		if strings.ToLower(m.Key) == "output error" {
			return &message.Message[string, string]{
					Key:   "error",
					Value: "error",
				},
				ss,
				nil
		}

		if strings.ToLower(m.Key) == "empty" {
			return nil,
				ss,
				nil
		}

		ss.Content += m.Value

		return &message.Message[string, string]{
				Key:   m.Key + "-updated",
				Value: m.Value + "-updated",
			},
			ss,
			nil
	}

	converted := stateful.ConvertOneToOne(oneToOne, crappyStringFormat, crappyStringFormat, crappyStringFormat, crappyStringFormat, crappyStringFormat)

	for _, testcase := range testcases {

		t.Run(testcase.name, func(t *testing.T) {
			assert := assert.New(t)

			output, outputState, err := converted(context.Background(), testcase.input, stateful.NewSingleState(testcase.inputState))

			assert.Equal(testcase.outputState, outputState.Content)

			if len(testcase.err) > 0 {
				assert.Equal(0, len(output))
				assert.EqualError(err, testcase.err)
				return
			}

			if testcase.empty {
				assert.Equal(0, len(output))
				return
			}

			assert.Equal(1, len(output))
			assert.Equal(testcase.output, output[0])
		})
	}
}
