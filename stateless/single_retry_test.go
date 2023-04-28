package stateless_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/avast/retry-go"
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/reflect"
	"github.com/stretchr/testify/assert"
)

func TestSingleRetry(t *testing.T) {
	retryRuntime := runtime_retry.NewRetry(
		runtime_retry.WithRetryOption(
			retry.Attempts(3),
			retry.Delay(10*time.Millisecond),
			retry.MaxDelay(time.Second),
			retry.MaxJitter(time.Second),
			retry.DelayType(retry.BackOffDelay),
		),
	)

	retryFn := stateless.NewSingleRetry(
		stateless.WithSingleRetryRuntime(retryRuntime),
		stateless.WithSingleRetryNextFunction(func(ctx context.Context, m message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
			countStr := string(m.Value)
			countInt := reflect.GetInt64(countStr)

			count := runtime_retry.GetTryCount(ctx)

			if countInt > count {
				return make([]message.Message[message.Bytes, message.Bytes], 0), errors.New("count not done")
			}

			return []message.Message[message.Bytes, message.Bytes]{
				{
					Key:   m.Key,
					Value: []byte(reflect.GetString(count)),
				},
			}, nil
		}),
	)

	cases := []struct {
		name         string
		inputMessage message.Message[message.Bytes, message.Bytes]
		countExpect  string
		err          error
	}{
		{
			name: "try once",
			inputMessage: message.Message[[]byte, []byte]{
				Key:   []byte("1"),
				Value: []byte("1"),
			},
			countExpect: "1",
		},
		{
			name: "try twice",
			inputMessage: message.Message[[]byte, []byte]{
				Key:   []byte("2"),
				Value: []byte("2"),
			},
			countExpect: "2",
		},
		{
			name: "try four times",
			inputMessage: message.Message[[]byte, []byte]{
				Key:   []byte("4"),
				Value: []byte("4"),
			},
			err:         stateless.ErrorRetryAttempt,
			countExpect: "3",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert := assert.New(t)

			r, err := retryFn(context.Background(), c.inputMessage)

			if c.err == nil {
				assert.NoError(err)
				assert.Equal(c.countExpect, string(r[0].Value))
			} else {
				assert.ErrorIs(err, c.err)
			}
		})
	}
}