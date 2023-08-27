package stateful_test

import (
	"context"
	"testing"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/protobuf"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/structure"
	"github.com/stretchr/testify/assert"
)

func TestSingleDeduplicate(t *testing.T) {

	// common data config
	inputMessage := flow.Message[string, string]{
		Offset:    100,
		Partition: 1,
		Key:       "input",
		Value:     "input",
		Headers:   map[string][][]byte{},
	}
	byteInputMessage, _ := flow.Convert(
		inputMessage,
		format.String(),
		format.String(),
		format.Bytes(),
		format.Bytes(),
	)
	inputBytes, _ := format.Convert(
		byteInputMessage,
		flow.Format(),
		format.Bytes(),
	)

	outputMessage := flow.Message[string, string]{
		Offset:    100,
		Partition: 1,
		Key:       "test",
		Value:     "test",
		Headers:   map[string][][]byte{},
	}
	byteOutputMessage, _ := flow.Convert(
		outputMessage,
		format.String(),
		format.String(),
		format.Bytes(),
		format.Bytes(),
	)
	outputBytes, _ := format.Convert(
		byteOutputMessage,
		flow.Format(),
		format.Bytes(),
	)

	cases := []struct {
		name         string
		inputMessage flow.Message[structure.Bytes, structure.Bytes]
		inputState   stateful.State[structure.Bytes]
		assertions   func(
			assert *assert.Assertions,
			inputMessage flow.Message[structure.Bytes, structure.Bytes],
			inputState stateful.State[structure.Bytes],
			resultMessages []flow.Message[structure.Bytes, structure.Bytes],
			resultState stateful.State[structure.Bytes],
			resultError error)
	}{
		{
			name:         "TestSingleStatefulDeduplicateNilInternalStateShouldApply",
			inputMessage: byteInputMessage,
			inputState:   stateful.State[structure.Bytes]{},
			assertions: func(
				assert *assert.Assertions,
				inputMessage flow.Message[structure.Bytes, structure.Bytes],
				inputState stateful.State[structure.Bytes],
				resultMessages []flow.Message[structure.Bytes, structure.Bytes],
				resultState stateful.State[structure.Bytes],
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
			inputState: stateful.State[structure.Bytes]{
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
				inputMessage flow.Message[structure.Bytes, structure.Bytes],
				inputState stateful.State[structure.Bytes],
				resultMessages []flow.Message[structure.Bytes, structure.Bytes],
				resultState stateful.State[structure.Bytes],
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
			inputState: stateful.State[structure.Bytes]{
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
							Messages: []structure.Bytes{outputBytes},
						},
					},
				},
			},
			assertions: func(
				assert *assert.Assertions,
				inputMessage flow.Message[structure.Bytes, structure.Bytes],
				inputState stateful.State[structure.Bytes],
				resultMessages []flow.Message[structure.Bytes, structure.Bytes],
				resultState stateful.State[structure.Bytes],
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
			inputState: stateful.State[structure.Bytes]{
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
							Messages: []structure.Bytes{outputBytes},
						},
					},
				},
			},
			assertions: func(
				assert *assert.Assertions,
				inputMessage flow.Message[structure.Bytes, structure.Bytes],
				inputState stateful.State[structure.Bytes],
				resultMessages []flow.Message[structure.Bytes, structure.Bytes],
				resultState stateful.State[structure.Bytes],
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
			inputState: stateful.State[structure.Bytes]{
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
							Messages: []structure.Bytes{outputBytes},
						},
					},
				},
			},
			assertions: func(
				assert *assert.Assertions,
				inputMessage flow.Message[structure.Bytes, structure.Bytes],
				inputState stateful.State[structure.Bytes],
				resultMessages []flow.Message[structure.Bytes, structure.Bytes],
				resultState stateful.State[structure.Bytes],
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

			fn := stateful.NewDeduplicate(
				stateful.WithDeduplicateNextFunction(
					func(c context.Context, m flow.Message[structure.Bytes, structure.Bytes], inState stateful.State[structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], stateful.State[structure.Bytes], error) {
						return []flow.Message[structure.Bytes, structure.Bytes]{m}, inState, nil
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
