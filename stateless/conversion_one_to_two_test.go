package stateless_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/flows/test_helper"
	"github.com/stretchr/testify/assert"
)

func TestConvertOneToTwo(t *testing.T) {
	testcases := []struct {
		name    string
		input   flow.Message[[]byte, []byte]
		output1 flow.Message[[]byte, []byte]
		output2 flow.Message[[]byte, []byte]
		oneNil  bool
		twoNil  bool
		err     string
	}{
		{
			name: "basic conversion",
			input: flow.Message[[]byte, []byte]{
				Key:   []byte("k"),
				Value: []byte("v"),
			},
			output1: flow.Message[[]byte, []byte]{
				Key:   []byte("k-1"),
				Value: []byte("v-1"),
			},
			output2: flow.Message[[]byte, []byte]{
				Key:   []byte("k-2"),
				Value: []byte("v-2"),
			},
			oneNil: false,
			twoNil: false,
			err:    "",
		},
		{
			name: "one nil",
			input: flow.Message[[]byte, []byte]{
				Key:   []byte("empty-1"),
				Value: []byte("v"),
			},
			output2: flow.Message[[]byte, []byte]{
				Key:   []byte("empty-1-2"),
				Value: []byte("v-2"),
			},
			oneNil: true,
			twoNil: false,
			err:    "",
		},
		{
			name: "two nil",
			input: flow.Message[[]byte, []byte]{
				Key:   []byte("empty-2"),
				Value: []byte("v"),
			},
			output1: flow.Message[[]byte, []byte]{
				Key:   []byte("empty-2-1"),
				Value: []byte("v-1"),
			},
			oneNil: false,
			twoNil: true,
			err:    "",
		},
		{
			name: "both nil",
			input: flow.Message[[]byte, []byte]{
				Key:   []byte("empty"),
				Value: []byte("v"),
			},
			oneNil: true,
			twoNil: true,
			err:    "",
		},
		{
			name: "error input conversion",
			input: flow.Message[[]byte, []byte]{
				Key:   []byte("error"),
				Value: []byte("error"),
			},
			oneNil: false,
			twoNil: false,
			err:    "error",
		},
		{
			name: "error execute",
			input: flow.Message[[]byte, []byte]{
				Key:   []byte("stupid error"),
				Value: []byte("v"),
			},
			oneNil: false,
			twoNil: false,
			err:    "stupid error",
		},
		{
			name: "error output 1 conversion",
			input: flow.Message[[]byte, []byte]{
				Key:   []byte("output-1 error"),
				Value: []byte("v"),
			},
			err: "error",
		},
		{
			name: "error output 2 conversion",
			input: flow.Message[[]byte, []byte]{
				Key:   []byte("output-2 error"),
				Value: []byte("v"),
			},
			err: "error",
		},
	}

	crappyStringFormat := test_helper.CrappyStringFormat()
	oneToTwo := func(ctx context.Context, m flow.Message[string, string]) (*flow.Message[string, string], *flow.Message[string, string], error) {
		if strings.ToLower(m.Key) == "stupid error" {
			return &flow.Message[string, string]{}, &flow.Message[string, string]{}, errors.New(m.Key)
		}

		if strings.ToLower(m.Key) == "output-1 error" {
			return &flow.Message[string, string]{
					Key:   "error",
					Value: "error",
				},
				&flow.Message[string, string]{
					Key:   m.Key + "-2",
					Value: m.Value + "-2",
				},
				nil
		}

		if strings.ToLower(m.Key) == "output-2 error" {
			return &flow.Message[string, string]{
					Key:   m.Key + "-1",
					Value: m.Value + "-1",
				},
				&flow.Message[string, string]{
					Key:   "error",
					Value: "error",
				},
				nil
		}

		if strings.ToLower(m.Key) == "empty-1" {
			return nil,
				&flow.Message[string, string]{
					Key:   m.Key + "-2",
					Value: m.Value + "-2",
				},
				nil
		}

		if strings.ToLower(m.Key) == "empty-2" {
			return &flow.Message[string, string]{
					Key:   m.Key + "-1",
					Value: m.Value + "-1",
				},
				nil,
				nil
		}

		if strings.ToLower(m.Key) == "empty" {
			return nil, nil, nil
		}

		return &flow.Message[string, string]{
				Key:   m.Key + "-1",
				Value: m.Value + "-1",
			},
			&flow.Message[string, string]{
				Key:   m.Key + "-2",
				Value: m.Value + "-2",
			},
			nil
	}

	converted := stateless.ConvertOneToTwo(oneToTwo, crappyStringFormat, crappyStringFormat, crappyStringFormat, crappyStringFormat, crappyStringFormat, crappyStringFormat)

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			assert := assert.New(t)

			output, err := converted(context.Background(), testcase.input)
			if len(testcase.err) > 0 {
				assert.Equal(0, len(output))
				assert.Contains(err.Error(), testcase.err)
				return
			}

			expectedCount := 2
			expectedMessages := make([]flow.Message[[]byte, []byte], 0)

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
