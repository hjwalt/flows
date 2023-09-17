package stateful_error_handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateful/stateful_error_handler"
	"github.com/hjwalt/flows/stateful/stateful_mock"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/structure"
	"github.com/stretchr/testify/assert"
)

func TestFn(t *testing.T) {
	next := stateful_error_handler.New(
		stateful.ConvertTopicOneToOne(stateful_mock.MockOneToOne, format.String(), stateful_mock.InputTopic, stateful_mock.OutputTopic),
		func(ctx context.Context, m flow.Message[structure.Bytes, structure.Bytes], ss stateful.State[structure.Bytes], err error) ([]flow.Message[structure.Bytes, structure.Bytes], stateful.State[structure.Bytes], error) {
			if !errors.Is(err, stateful_mock.ErrRecoverable) {
				return flow.EmptySlice(), ss, err
			}

			k := string(m.Key)
			v := string(m.Value)

			rk := []byte(k + "-recovered")
			rv := []byte(v + "-recovered")

			ss.Content = []byte(string(ss.Content) + "recovered")

			return []flow.Message[structure.Bytes, structure.Bytes]{{Key: rk, Value: rv}}, ss, nil
		},
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
		},
		{
			name:          "basic recoverable case",
			key:           "recoverable",
			value:         "recoverable",
			state:         "state",
			expectedKey:   "recoverable-recovered",
			expectedValue: "recoverable-recovered",
			expectedState: "staterecovered",
			expectedError: nil,
		},
		{
			name:          "basic irrecoverable case",
			key:           "irrecoverable",
			value:         "irrecoverable",
			state:         "state",
			expectedKey:   "n/a",
			expectedValue: "n/a",
			expectedState: "state",
			expectedError: stateful_mock.ErrIrrecoverable,
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
				assert.Equal(1, len(rm))

				rk := string(rm[0].Key)
				rv := string(rm[0].Value)
				rs := string(rss.Content)

				assert.Equal(tc.expectedKey, rk)
				assert.Equal(tc.expectedValue, rv)
				assert.Equal(tc.expectedState, rs)
			} else {
				assert.ErrorIs(rerr, tc.expectedError)
			}
		})
	}
}
