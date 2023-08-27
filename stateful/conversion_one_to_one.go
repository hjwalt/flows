package stateful

import (
	"context"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/topic"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/structure"
)

func ConvertOneToOne[S any, IK any, IV any, OK any, OV any](
	source OneToOneFunction[S, IK, IV, OK, OV],
	s format.Format[S],
	ik format.Format[IK],
	iv format.Format[IV],
	ok format.Format[OK],
	ov format.Format[OV],
) SingleFunction {
	return func(ctx context.Context, m flow.Message[structure.Bytes, structure.Bytes], ss State[structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], State[structure.Bytes], error) {

		formattedMessage, unmarshalError := flow.Convert(m, format.Bytes(), format.Bytes(), ik, iv)
		if unmarshalError != nil {
			return flow.EmptySlice(), ss, unmarshalError
		}

		formattedState, stateUnmarshalError := ConvertSingleState(ss, format.Bytes(), s)
		if stateUnmarshalError != nil {
			return flow.EmptySlice(), ss, stateUnmarshalError
		}

		res, nextState, fnError := source(ctx, formattedMessage, formattedState)
		if fnError != nil {
			return flow.EmptySlice(), ss, fnError
		}

		byteResultMessages := make([]flow.Message[[]byte, []byte], 0)

		if res != nil {
			bytesResMessage, marshalError := flow.Convert(*res, ok, ov, format.Bytes(), format.Bytes())
			if marshalError != nil {
				return flow.EmptySlice(), ss, marshalError
			}
			byteResultMessages = append(byteResultMessages, bytesResMessage)
		}

		bytesNextState, stateMarshalError := ConvertSingleState(nextState, s, format.Bytes())
		if stateMarshalError != nil {
			return flow.EmptySlice(), ss, stateMarshalError
		}

		return byteResultMessages, bytesNextState, nil
	}
}

func ConvertTopicOneToOne[S any, IK any, IV any, OK any, OV any](
	source OneToOneFunction[S, IK, IV, OK, OV],
	s format.Format[S],
	inputTopic topic.Topic[IK, IV],
	outputTopic topic.Topic[OK, OV],
) SingleFunction {
	return func(ctx context.Context, m flow.Message[structure.Bytes, structure.Bytes], ss State[structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], State[structure.Bytes], error) {

		formattedMessage, unmarshalError := flow.Convert(m, format.Bytes(), format.Bytes(), inputTopic.KeyFormat(), inputTopic.ValueFormat())
		if unmarshalError != nil {
			return flow.EmptySlice(), ss, unmarshalError
		}

		formattedState, stateUnmarshalError := ConvertSingleState(ss, format.Bytes(), s)
		if stateUnmarshalError != nil {
			return flow.EmptySlice(), ss, stateUnmarshalError
		}

		res, nextState, fnError := source(ctx, formattedMessage, formattedState)
		if fnError != nil {
			return flow.EmptySlice(), ss, fnError
		}

		byteResultMessages := make([]flow.Message[[]byte, []byte], 0)

		if res != nil {
			bytesResMessage, marshalError := flow.Convert(*res, outputTopic.KeyFormat(), outputTopic.ValueFormat(), format.Bytes(), format.Bytes())
			if marshalError != nil {
				return flow.EmptySlice(), ss, marshalError
			}
			bytesResMessage.Topic = outputTopic.Topic()
			byteResultMessages = append(byteResultMessages, bytesResMessage)
		}

		bytesNextState, stateMarshalError := ConvertSingleState(nextState, s, format.Bytes())
		if stateMarshalError != nil {
			return flow.EmptySlice(), ss, stateMarshalError
		}

		return byteResultMessages, bytesNextState, nil
	}
}
