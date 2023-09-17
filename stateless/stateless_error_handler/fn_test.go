package stateless_error_handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/stateless/stateless_error_handler"
	"github.com/hjwalt/flows/stateless/stateless_mock"
	"github.com/hjwalt/flows/stateless/stateless_one_to_one"
	"github.com/hjwalt/runway/structure"
	"github.com/stretchr/testify/assert"
)

func TestFn(t *testing.T) {
	next := stateless_error_handler.New(
		stateless_one_to_one.New(stateless_mock.MockOneToOne, stateless_mock.InputTopic, stateless_mock.OutputTopic),
		func(ctx context.Context, m flow.Message[structure.Bytes, structure.Bytes], err error) ([]flow.Message[structure.Bytes, structure.Bytes], error) {
			if !errors.Is(err, stateless_mock.ErrRecoverable) {
				return flow.EmptySlice(), err
			}

			k := string(m.Key)
			v := string(m.Value)

			rk := []byte(k + "-recovered")
			rv := []byte(v + "-recovered")

			return []flow.Message[structure.Bytes, structure.Bytes]{{Key: rk, Value: rv}}, nil
		},
	)

	cases := []struct {
		name          string
		key           string
		value         string
		expectedKey   string
		expectedValue string
		expectedError error
	}{
		{
			name:          "basic success case",
			key:           "test",
			value:         "test",
			expectedKey:   "test",
			expectedValue: "test-updated",
			expectedError: nil,
		},
		{
			name:          "basic recoverable case",
			key:           "recoverable",
			value:         "recoverable",
			expectedKey:   "recoverable-recovered",
			expectedValue: "recoverable-recovered",
			expectedError: nil,
		},
		{
			name:          "basic irrecoverable case",
			key:           "irrecoverable",
			value:         "irrecoverable",
			expectedKey:   "n/a",
			expectedValue: "n/a",
			expectedError: stateless_mock.ErrIrrecoverable,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)

			k := []byte(tc.key)
			v := []byte(tc.value)

			rm, rerr := next(context.Background(), flow.Message[structure.Bytes, structure.Bytes]{Key: k, Value: v})
			if tc.expectedError == nil {
				assert.NoError(rerr)
				assert.Equal(1, len(rm))

				rk := string(rm[0].Key)
				rv := string(rm[0].Value)

				assert.Equal(tc.expectedKey, rk)
				assert.Equal(tc.expectedValue, rv)
			} else {
				assert.ErrorIs(rerr, tc.expectedError)
			}
		})
	}
}
