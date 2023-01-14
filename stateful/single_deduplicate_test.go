package stateful_test

import (
	"context"
	"testing"

	"github.com/hjwalt/flows/format"
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/protobuf"
	"github.com/hjwalt/flows/stateful"
	"github.com/stretchr/testify/assert"
)

func TestSingleDeduplicate(t *testing.T) {

	// common data config
	inputMessage := message.Message[string, string]{
		Offset:    100,
		Partition: 1,
		Key:       "input",
		Value:     "input",
		Headers:   map[string][][]byte{},
	}
	byteInputMessage, _ := message.Convert(
		inputMessage,
		format.String(),
		format.String(),
		format.Bytes(),
		format.Bytes(),
	)
	inputBytes, _ := format.Convert(
		byteInputMessage,
		message.Format(),
		format.Bytes(),
	)

	outputMessage := message.Message[string, string]{
		Offset:    100,
		Partition: 1,
		Key:       "test",
		Value:     "test",
		Headers:   map[string][][]byte{},
	}
	byteOutputMessage, _ := message.Convert(
		outputMessage,
		format.String(),
		format.String(),
		format.Bytes(),
		format.Bytes(),
	)
	outputBytes, _ := format.Convert(
		byteOutputMessage,
		message.Format(),
		format.Bytes(),
	)

	cases := []struct {
		name         string
		inputMessage message.Message[message.Bytes, message.Bytes]
		inputState   stateful.SingleState[message.Bytes]
		assertions   func(
			assert *assert.Assertions,
			inputMessage message.Message[message.Bytes, message.Bytes],
			inputState stateful.SingleState[message.Bytes],
			resultMessages []message.Message[message.Bytes, message.Bytes],
			resultState stateful.SingleState[message.Bytes],
			resultError error)
	}{
		{
			name:         "TestSingleStatefulDeduplicateNilInternalStateShouldApply",
			inputMessage: byteInputMessage,
			inputState:   stateful.SingleState[message.Bytes]{},
			assertions: func(
				assert *assert.Assertions,
				inputMessage message.Message[message.Bytes, message.Bytes],
				inputState stateful.SingleState[message.Bytes],
				resultMessages []message.Message[message.Bytes, message.Bytes],
				resultState stateful.SingleState[message.Bytes],
				resultError error,
			) {
				assert.NotNil(resultMessages)
				assert.NotNil(resultState)
				assert.NoError(resultError)

				assert.NotEqual(inputState, resultState)

				assert.Equal(1, len(resultMessages))
				assert.Equal(byteInputMessage, resultMessages[0])

				assert.Equal(1, len(resultState.Results.GetV1().GetMessages()))
				assert.Equal(inputBytes, resultState.Results.GetV1().GetMessages()[0])

				assert.Equal(int64(100), resultState.Internal.GetV1().GetOffsetProgress()[1])
			},
		}, // end TestSingleStatefulDeduplicateNilInternalStateShouldApply
		{
			name:         "TestSingleStatefulDeduplicateSkippingWithoutOutput",
			inputMessage: byteInputMessage,
			inputState: stateful.SingleState[message.Bytes]{
				Internal: &protobuf.State{
					State: &protobuf.State_V1{
						V1: &protobuf.StateV1{
							OffsetProgress: map[int32]int64{
								1: 101,
							},
						},
					},
				},
				Results: &protobuf.Results{
					Result: &protobuf.Results_V1{
						V1: &protobuf.ResultV1{
							Messages: make([][]byte, 0),
						},
					},
				},
			},
			assertions: func(
				assert *assert.Assertions,
				inputMessage message.Message[message.Bytes, message.Bytes],
				inputState stateful.SingleState[message.Bytes],
				resultMessages []message.Message[message.Bytes, message.Bytes],
				resultState stateful.SingleState[message.Bytes],
				resultError error,
			) {
				assert.NotNil(resultMessages)
				assert.NotNil(resultState)
				assert.NoError(resultError)

				assert.Equal(inputState, resultState)
			},
		}, // end TestSingleStatefulDeduplicateSkippingWithoutOutput
		{
			name:         "TestSingleStatefulDeduplicateSkippingWithOutput",
			inputMessage: byteInputMessage,
			inputState: stateful.SingleState[message.Bytes]{
				Internal: &protobuf.State{
					State: &protobuf.State_V1{
						V1: &protobuf.StateV1{
							OffsetProgress: map[int32]int64{
								1: 100,
							},
						},
					},
				},
				Results: &protobuf.Results{
					Result: &protobuf.Results_V1{
						V1: &protobuf.ResultV1{
							Messages: []message.Bytes{outputBytes},
						},
					},
				},
			},
			assertions: func(
				assert *assert.Assertions,
				inputMessage message.Message[message.Bytes, message.Bytes],
				inputState stateful.SingleState[message.Bytes],
				resultMessages []message.Message[message.Bytes, message.Bytes],
				resultState stateful.SingleState[message.Bytes],
				resultError error,
			) {
				assert.NotNil(resultMessages)
				assert.NotNil(resultState)
				assert.NoError(resultError)

				assert.Equal(1, len(resultMessages))
				assert.Equal(byteOutputMessage, resultMessages[0])
			},
		}, // end TestSingleStatefulDeduplicateSkippingWithOutput
		{
			name:         "TestSingleStatefulDeduplicateActingWithoutResultAppend",
			inputMessage: byteInputMessage,
			inputState: stateful.SingleState[message.Bytes]{
				Internal: &protobuf.State{
					State: &protobuf.State_V1{
						V1: &protobuf.StateV1{
							OffsetProgress: map[int32]int64{
								1: 99,
							},
						},
					},
				},
				Results: &protobuf.Results{
					Result: &protobuf.Results_V1{
						V1: &protobuf.ResultV1{
							Messages: []message.Bytes{outputBytes},
						},
					},
				},
			},
			assertions: func(
				assert *assert.Assertions,
				inputMessage message.Message[message.Bytes, message.Bytes],
				inputState stateful.SingleState[message.Bytes],
				resultMessages []message.Message[message.Bytes, message.Bytes],
				resultState stateful.SingleState[message.Bytes],
				resultError error,
			) {

				assert.NotNil(resultMessages)
				assert.NotNil(resultState)
				assert.NoError(resultError)

				assert.Equal(1, len(resultMessages))
				assert.Equal(byteInputMessage, resultMessages[0])

				assert.Equal(1, len(resultState.Results.GetV1().GetMessages()))
				assert.Equal(inputBytes, resultState.Results.GetV1().GetMessages()[0])

				assert.Equal(int64(100), resultState.Internal.GetV1().GetOffsetProgress()[1])
			},
		}, // end TestSingleStatefulDeduplicateActingWithoutResultAppend
		{
			name:         "TestSingleStatefulDeduplicateActingWithResultAppend",
			inputMessage: byteInputMessage,
			inputState: stateful.SingleState[message.Bytes]{
				Internal: &protobuf.State{
					State: &protobuf.State_V1{
						V1: &protobuf.StateV1{
							OffsetProgress: map[int32]int64{},
						},
					},
				},
				Results: &protobuf.Results{
					Result: &protobuf.Results_V1{
						V1: &protobuf.ResultV1{
							Messages: []message.Bytes{outputBytes},
						},
					},
				},
			},
			assertions: func(
				assert *assert.Assertions,
				inputMessage message.Message[message.Bytes, message.Bytes],
				inputState stateful.SingleState[message.Bytes],
				resultMessages []message.Message[message.Bytes, message.Bytes],
				resultState stateful.SingleState[message.Bytes],
				resultError error,
			) {

				assert.NotNil(resultMessages)
				assert.NotNil(resultState)
				assert.NoError(resultError)

				assert.Equal(2, len(resultMessages))
				assert.Equal(byteOutputMessage, resultMessages[0])
				assert.Equal(byteInputMessage, resultMessages[1])

				assert.Equal(1, len(resultState.Results.GetV1().GetMessages()))
				assert.Equal(inputBytes, resultState.Results.GetV1().GetMessages()[0])

				assert.Equal(int64(100), resultState.Internal.GetV1().GetOffsetProgress()[1])
			},
		}, // end TestSingleStatefulDeduplicateActingWithResultAppend
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			assert := assert.New(t)

			fn := stateful.NewSingleStatefulDeduplicate(
				stateful.WithSingleStatefulDeduplicateNextFunction(
					func(c context.Context, m message.Message[message.Bytes, message.Bytes], inState stateful.SingleState[message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], stateful.SingleState[message.Bytes], error) {
						return []message.Message[message.Bytes, message.Bytes]{m}, inState, nil
					},
				),
			)

			ctx := context.Background()
			resultMessages, resultState, resultErr := fn(ctx, c.inputMessage, c.inputState)
			c.assertions(
				assert,
				c.inputMessage,
				c.inputState,
				resultMessages,
				resultState,
				resultErr,
			)
		})
	}
}
