package stateful

import (
	"context"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/topic"
	"github.com/hjwalt/runway/format"
)

func ConvertOneToOne[S any, IK any, IV any, OK any, OV any](
	source OneToOneFunction[S, IK, IV, OK, OV],
	s format.Format[S],
	ik format.Format[IK],
	iv format.Format[IV],
	ok format.Format[OK],
	ov format.Format[OV],
) SingleFunction {
	return func(ctx context.Context, m message.Message[message.Bytes, message.Bytes], ss State[message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], State[message.Bytes], error) {

		formattedMessage, unmarshalError := message.Convert(m, format.Bytes(), format.Bytes(), ik, iv)
		if unmarshalError != nil {
			return make([]message.Message[[]byte, []byte], 0), ss, unmarshalError
		}

		formattedState, stateUnmarshalError := ConvertSingleState(ss, format.Bytes(), s)
		if stateUnmarshalError != nil {
			return make([]message.Message[[]byte, []byte], 0), ss, stateUnmarshalError
		}

		res, nextState, fnError := source(ctx, formattedMessage, formattedState)
		if fnError != nil {
			return make([]message.Message[[]byte, []byte], 0), ss, fnError
		}

		byteResultMessages := make([]message.Message[[]byte, []byte], 0)

		if res != nil {
			bytesResMessage, marshalError := message.Convert(*res, ok, ov, format.Bytes(), format.Bytes())
			if marshalError != nil {
				return make([]message.Message[[]byte, []byte], 0), ss, marshalError
			}
			byteResultMessages = append(byteResultMessages, bytesResMessage)
		}

		bytesNextState, stateMarshalError := ConvertSingleState(nextState, s, format.Bytes())
		if stateMarshalError != nil {
			return make([]message.Message[[]byte, []byte], 0), ss, stateMarshalError
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
	return func(ctx context.Context, m message.Message[message.Bytes, message.Bytes], ss State[message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], State[message.Bytes], error) {

		formattedMessage, unmarshalError := message.Convert(m, format.Bytes(), format.Bytes(), inputTopic.KeyFormat(), inputTopic.ValueFormat())
		if unmarshalError != nil {
			return make([]message.Message[[]byte, []byte], 0), ss, unmarshalError
		}

		formattedState, stateUnmarshalError := ConvertSingleState(ss, format.Bytes(), s)
		if stateUnmarshalError != nil {
			return make([]message.Message[[]byte, []byte], 0), ss, stateUnmarshalError
		}

		res, nextState, fnError := source(ctx, formattedMessage, formattedState)
		if fnError != nil {
			return make([]message.Message[[]byte, []byte], 0), ss, fnError
		}

		byteResultMessages := make([]message.Message[[]byte, []byte], 0)

		if res != nil {
			bytesResMessage, marshalError := message.Convert(*res, outputTopic.KeyFormat(), outputTopic.ValueFormat(), format.Bytes(), format.Bytes())
			if marshalError != nil {
				return make([]message.Message[[]byte, []byte], 0), ss, marshalError
			}
			bytesResMessage.Topic = outputTopic.Topic()
			byteResultMessages = append(byteResultMessages, bytesResMessage)
		}

		bytesNextState, stateMarshalError := ConvertSingleState(nextState, s, format.Bytes())
		if stateMarshalError != nil {
			return make([]message.Message[[]byte, []byte], 0), ss, stateMarshalError
		}

		return byteResultMessages, bytesNextState, nil
	}
}
