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

func TestConvertOneToTwo(t *testing.T) {
	testcases := []struct {
		name        string
		input       message.Message[[]byte, []byte]
		inputState  []byte
		output1     message.Message[[]byte, []byte]
		output2     message.Message[[]byte, []byte]
		outputState []byte
		oneNil      bool
		twoNil      bool
		err         string
	}{
		{
			name: "basic conversion",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("k"),
				Value: []byte("v"),
			},
			output1: message.Message[[]byte, []byte]{
				Key:   []byte("k-1"),
				Value: []byte("v-1"),
			},
			output2: message.Message[[]byte, []byte]{
				Key:   []byte("k-2"),
				Value: []byte("v-2"),
			},
			inputState:  []byte("state"),
			outputState: []byte("state-1-2"),
			err:         "",
			oneNil:      false,
			twoNil:      false,
		},
		{
			name: "one nil",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("empty-1"),
				Value: []byte("v"),
			},
			output2: message.Message[[]byte, []byte]{
				Key:   []byte("empty-1-2"),
				Value: []byte("v-2"),
			},
			inputState:  []byte("state"),
			outputState: []byte("state-2"),
			oneNil:      true,
			twoNil:      false,
			err:         "",
		},
		{
			name: "two nil",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("empty-2"),
				Value: []byte("v"),
			},
			output1: message.Message[[]byte, []byte]{
				Key:   []byte("empty-2-1"),
				Value: []byte("v-1"),
			},
			inputState:  []byte("state"),
			outputState: []byte("state-1"),
			oneNil:      false,
			twoNil:      true,
			err:         "",
		},
		{
			name: "both nil",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("empty"),
				Value: []byte("v"),
			},
			inputState:  []byte("state"),
			outputState: []byte("state"),
			oneNil:      true,
			twoNil:      true,
			err:         "",
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
			oneNil:      false,
			twoNil:      false,
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
			oneNil:      false,
			twoNil:      false,
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
			oneNil:      false,
			twoNil:      false,
		},
		{
			name: "error output 1 conversion",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("output-1 error"),
				Value: []byte("v"),
			},
			inputState:  []byte("state"),
			outputState: []byte("state"),
			err:         "error",
		},
		{
			name: "error output 2 conversion",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("output-2 error"),
				Value: []byte("v"),
			},
			inputState:  []byte("state"),
			outputState: []byte("state"),
			err:         "error",
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
			oneNil:      false,
			twoNil:      false,
		},
	}

	crappyStringFormat := mock.CrappyStringFormat()
	oneToTwo := func(ctx context.Context, m message.Message[string, string], ss stateful.SingleState[string]) (*message.Message[string, string], *message.Message[string, string], stateful.SingleState[string], error) {

		if strings.ToLower(m.Key) == "stupid error" {
			return &message.Message[string, string]{},
				&message.Message[string, string]{},
				ss,
				errors.New(m.Key)
		}

		if strings.ToLower(m.Key) == "output-1 error" {
			return &message.Message[string, string]{
					Key:   "error",
					Value: "error",
				},
				&message.Message[string, string]{
					Key:   m.Key + "-2",
					Value: m.Value + "-2",
				},
				ss,
				nil
		}

		if strings.ToLower(m.Key) == "output-2 error" {
			return &message.Message[string, string]{
					Key:   m.Key + "-1",
					Value: m.Value + "-1",
				},
				&message.Message[string, string]{
					Key:   "error",
					Value: "error",
				},
				ss,
				nil
		}

		if strings.ToLower(m.Key) == "empty-1" {
			ss.Content += "-2"

			return nil,
				&message.Message[string, string]{
					Key:   m.Key + "-2",
					Value: m.Value + "-2",
				},
				ss,
				nil
		}

		if strings.ToLower(m.Key) == "empty-2" {
			ss.Content += "-1"

			return &message.Message[string, string]{
					Key:   m.Key + "-1",
					Value: m.Value + "-1",
				},
				nil,
				ss,
				nil
		}

		if strings.ToLower(m.Key) == "empty" {
			return nil, nil, ss, nil
		}

		if ss.Content == "erro" {
			ss.Content = "error"
		} else {
			ss.Content += "-1-2"
		}

		return &message.Message[string, string]{
				Key:   m.Key + "-1",
				Value: m.Value + "-1",
			},
			&message.Message[string, string]{
				Key:   m.Key + "-2",
				Value: m.Value + "-2",
			},
			ss,
			nil
	}

	converted := stateful.ConvertOneToTwo(oneToTwo, crappyStringFormat, crappyStringFormat, crappyStringFormat, crappyStringFormat, crappyStringFormat, crappyStringFormat, crappyStringFormat)

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

			expectedCount := 2
			expectedMessages := make([]message.Message[[]byte, []byte], 0)

			if testcase.oneNil {
				expectedCount -= 1
			} else {
				expectedMessages = append(expectedMessages, testcase.output1)
			}

			if testcase.twoNil {
				expectedCount -= 1
			} else {
				expectedMessages = append(expectedMessages, testcase.output2)
			}

			assert.Equal(expectedCount, len(output))
			assert.Equal(expectedMessages, output)
		})
	}
}
