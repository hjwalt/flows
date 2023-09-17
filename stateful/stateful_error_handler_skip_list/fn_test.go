package stateful_error_handler_skip_list_test

import (
	"context"
	"testing"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateful/stateful_error_handler"
	"github.com/hjwalt/flows/stateful/stateful_error_handler_skip_list"
	"github.com/hjwalt/flows/stateful/stateful_mock"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/structure"
	"github.com/stretchr/testify/assert"
)

func TestFn(t *testing.T) {
	next := stateful_error_handler.New(
		stateful.ConvertTopicOneToOne(stateful_mock.MockOneToOne, format.String(), stateful_mock.InputTopic, stateful_mock.OutputTopic),
		stateful_error_handler_skip_list.Default(),
	)

	cases := []struct {
		name          string
		key           string
		value         string
		state         string
		expectedKey   string
		expectedValue string
		expectedState string
		expectedError error
		expectedEmpty bool
	}{
		{
			name:          "basic success case",
			key:           "test",
			value:         "test",
			state:         "state",
			expectedKey:   "test",
			expectedValue: "test-updated",
			expectedState: "statetest",
			expectedError: nil,
			expectedEmpty: false,
		},
		{
			name:          "ghastly recovered",
			key:           "ghastly",
			value:         "ghastly",
			state:         "state",
			expectedState: "state",
			expectedError: nil,
			expectedEmpty: true,
		},
		{
			name:          "haunter recovered",
			key:           "haunter",
			value:         "haunter",
			state:         "state",
			expectedState: "state",
			expectedError: nil,
			expectedEmpty: true,
		},
		{
			name:          "mock_error exploded",
			key:           "mock_error",
			value:         "mock_error",
			state:         "state",
			expectedState: "state",
			expectedError: stateful_mock.ErrMock,
			expectedEmpty: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)

			k := []byte(tc.key)
			v := []byte(tc.value)

			rm, rss, rerr := next(context.Background(), flow.Message[structure.Bytes, structure.Bytes]{Key: k, Value: v}, stateful.State[[]byte]{Content: []byte(tc.state)})
			if tc.expectedError == nil {
				assert.NoError(rerr)

				rs := string(rss.Content)

				assert.Equal(tc.expectedState, rs)
			} else {
				assert.ErrorIs(rerr, tc.expectedError)
			}

			if !tc.expectedEmpty {
				assert.Equal(1, len(rm))

				rk := string(rm[0].Key)
				rv := string(rm[0].Value)

				assert.Equal(tc.expectedKey, rk)
				assert.Equal(tc.expectedValue, rv)
			}
		})
	}
}
