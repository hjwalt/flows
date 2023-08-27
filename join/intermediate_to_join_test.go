package join_test

import (
	"context"
	"testing"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/join"
	"github.com/hjwalt/flows/protobuf"
	"github.com/hjwalt/flows/test_helper"
	"github.com/hjwalt/runway/structure"
	"github.com/stretchr/testify/assert"
)

func TestIntermediateToJoinMap(t *testing.T) {
	topics := make(chan string, 10)

	fn := join.NewIntermediateToJoinMap(
		join.WithIntermediateToJoinMapTransactionWrappedFunction(func(ctx context.Context, ms []flow.Message[structure.Bytes, structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], error) {
			for _, m := range ms {
				topics <- m.Topic
			}
			return ms, nil
		}),
	)

	cases := []struct {
		name   string
		topic  string
		key    []byte
		value  []byte
		err    error
		output flow.Message[structure.Bytes, structure.Bytes]
	}{
		{
			name:  "normal message",
			topic: "test-topic",
			key:   test_helper.MarshalNoError(t, join.IntermediateKeyFormat, &protobuf.JoinKey{PersistenceId: "key"}),
			value: test_helper.MarshalNoError(t, join.IntermediateValueFormat, flow.Message[structure.Bytes, structure.Bytes]{
				Topic: "test-topic",
				Key:   []byte("key"),
				Value: []byte("value"),
			}),
			output: flow.Message[structure.Bytes, structure.Bytes]{
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

			res, err := fn(context.Background(), []flow.Message[structure.Bytes, structure.Bytes]{
				{
					Topic: "intermediate-topic",
					Key:   c.key,
					Value: c.value,
				},
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
