package join_test

import (
	"context"
	"testing"

	"github.com/hjwalt/flows/join"
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/protobuf"
	"github.com/hjwalt/flows/test_helper"
	"github.com/stretchr/testify/assert"
)

func TestIntermediateToJoinMap(t *testing.T) {
	topics := make(chan string, 10)

	fn := join.NewIntermediateToJoinMap(
		join.WithIntermediateToJoinMapTransactionWrappedFunction(func(ctx context.Context, m message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
			topics <- m.Topic
			return []message.Message[message.Bytes, message.Bytes]{m}, nil
		}),
	)

	cases := []struct {
		name   string
		topic  string
		key    []byte
		value  []byte
		err    error
		output message.Message[message.Bytes, message.Bytes]
	}{
		{
			name:  "normal message",
			topic: "test-topic",
			key:   test_helper.MarshalNoError(t, join.IntermediateKeyFormat, &protobuf.JoinKey{PersistenceId: "key"}),
			value: test_helper.MarshalNoError(t, join.IntermediateValueFormat, message.Message[message.Bytes, message.Bytes]{
				Topic: "test-topic",
				Key:   []byte("key"),
				Value: []byte("value"),
			}),
			output: message.Message[message.Bytes, message.Bytes]{
				Topic:   "test-topic",
				Key:     []byte("key"),
				Value:   []byte("value"),
				Headers: make(map[string][][]byte),
			},
		},
		{
			name:  "deserialise error",
			topic: "test-topic",
			key:   test_helper.MarshalNoError(t, join.IntermediateKeyFormat, &protobuf.JoinKey{PersistenceId: "key"}),
			value: []byte("invalid"),
			err:   join.ErrorIntermediateToJoinDeserialiseMessage,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert := assert.New(t)

			res, err := fn(context.Background(), message.Message[message.Bytes, message.Bytes]{
				Topic: "intermediate-topic",
				Key:   c.key,
				Value: c.value,
			})

			if c.err == nil {
				assert.NoError(err)
				assert.Equal(1, len(topics))

				topic := <-topics
				assert.Equal(c.topic, topic)

				assert.Equal(1, len(res))
				assert.Equal(c.output, res[0])
			} else {
				assert.ErrorIs(err, join.ErrorIntermediateToJoinDeserialiseMessage)
			}

		})
	}
}
