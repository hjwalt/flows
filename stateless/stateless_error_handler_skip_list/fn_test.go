package stateless_error_handler_skip_list_test

import (
	"context"
	"testing"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/stateless/stateless_error_handler"
	"github.com/hjwalt/flows/stateless/stateless_error_handler_skip_list"
	"github.com/hjwalt/flows/stateless/stateless_mock"
	"github.com/hjwalt/flows/stateless/stateless_one_to_one"
	"github.com/hjwalt/runway/format"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	testcases := []struct {
		name        string
		inputTopic  flow.Topic[string, string]
		outputTopic flow.Topic[string, string]
		input       flow.Message[[]byte, []byte]
		output      flow.Message[[]byte, []byte]
		empty       bool
		err         []error
	}{
		{
			name:        "basic conversion",
			inputTopic:  flow.GenericTopic("input", format.Gengar(), format.Gengar()),
			outputTopic: flow.GenericTopic("output", format.Gengar(), format.Gengar()),
			input: flow.Message[[]byte, []byte]{
				Key:   []byte("k"),
				Value: []byte("v"),
			},
			output: flow.Message[[]byte, []byte]{
				Topic: "output",
				Key:   []byte("k"),
				Value: []byte("v-updated"),
			},
			err:   []error{},
			empty: false,
		},
		{
			name:        "empty result",
			inputTopic:  flow.GenericTopic("input", format.Gengar(), format.Gengar()),
			outputTopic: flow.GenericTopic("output", format.Gengar(), format.Gengar()),
			input: flow.Message[[]byte, []byte]{
				Key:   []byte("empty"),
				Value: []byte("v"),
			},
			err:   []error{},
			empty: true,
		},
		{
			name:        "error input conversion",
			inputTopic:  flow.GenericTopic("input", format.Gengar(), format.Gengar()),
			outputTopic: flow.GenericTopic("output", format.Gengar(), format.Gengar()),
			input: flow.Message[[]byte, []byte]{
				Key:   []byte("ghastly"),
				Value: []byte("v"),
			},
			err:   []error{},
			empty: true,
		},
		{
			name:        "error output conversion",
			inputTopic:  flow.GenericTopic("input", format.Gengar(), format.Gengar()),
			outputTopic: flow.GenericTopic("output", format.Gengar(), format.Gengar()),
			input: flow.Message[[]byte, []byte]{
				Key:   []byte("haunter"),
				Value: []byte("v"),
			},
			err:   []error{},
			empty: true,
		},
		{
			name:        "error execute",
			inputTopic:  flow.GenericTopic("input", format.Gengar(), format.Gengar()),
			outputTopic: flow.GenericTopic("output", format.Gengar(), format.Gengar()),
			input: flow.Message[[]byte, []byte]{
				Key:   []byte("mock_error"),
				Value: []byte("v"),
			},
			err:   []error{stateless_mock.ErrMock},
			empty: true,
		},
		{
			name:        "empty topic output",
			inputTopic:  flow.GenericTopic("input", format.Gengar(), format.Gengar()),
			outputTopic: flow.GenericTopic("", format.Gengar(), format.Gengar()),
			input: flow.Message[[]byte, []byte]{
				Key:   []byte("mock_topic"),
				Value: []byte("v"),
			},
			output: flow.Message[[]byte, []byte]{
				Topic: "mock_topic",
				Key:   []byte("mock_topic"),
				Value: []byte("v-mock_topic"),
			},
			err:   []error{},
			empty: false,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			assert := assert.New(t)

			handled := stateless_error_handler.New(
				stateless_one_to_one.New(stateless_mock.MockOneToOne, testcase.inputTopic, testcase.outputTopic),
				stateless_error_handler_skip_list.Default(),
			)

			output, err := handled(context.Background(), testcase.input)
			if len(testcase.err) > 0 {
				for _, terr := range testcase.err {
					assert.ErrorIs(err, terr)
				}
			}

			if testcase.empty {
				assert.Equal(0, len(output))
			} else {
				assert.Equal(1, len(output))
				assert.Equal(testcase.output, output[0])
			}

		})
	}
}
