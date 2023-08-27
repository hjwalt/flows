package stateful

import (
	"context"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/structure"
)

func ConvertOneToTwo[S any, IK any, IV any, OK1 any, OV1 any, OK2 any, OV2 any](
	source OneToTwoFunction[S, IK, IV, OK1, OV1, OK2, OV2],
	s format.Format[S],
	ik format.Format[IK],
	iv format.Format[IV],
	ok1 format.Format[OK1],
	ov1 format.Format[OV1],
	ok2 format.Format[OK2],
	ov2 format.Format[OV2],
) SingleFunction {
	return func(ctx context.Context, m flow.Message[structure.Bytes, structure.Bytes], ss State[structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], State[structure.Bytes], error) {

		formattedMessage, unmarshalError := flow.Convert(m, format.Bytes(), format.Bytes(), ik, iv)
		if unmarshalError != nil {
			return make([]flow.Message[[]byte, []byte], 0), ss, unmarshalError
		}

		formattedState, stateUnmarshalError := ConvertSingleState(ss, format.Bytes(), s)
		if stateUnmarshalError != nil {
			return make([]flow.Message[[]byte, []byte], 0), ss, stateUnmarshalError
		}

		res1, res2, nextState, fnError := source(ctx, formattedMessage, formattedState)
		if fnError != nil {
			return make([]flow.Message[[]byte, []byte], 0), ss, fnError
		}

		byteResultMessages := make([]flow.Message[[]byte, []byte], 0)

		if res1 != nil {
			bytesResMessage, marshalError := flow.Convert(*res1, ok1, ov1, format.Bytes(), format.Bytes())
			if marshalError != nil {
				return make([]flow.Message[[]byte, []byte], 0), ss, marshalError
			}
			byteResultMessages = append(byteResultMessages, bytesResMessage)
		}

		if res2 != nil {
			bytesResMessage, marshalError := flow.Convert(*res2, ok2, ov2, format.Bytes(), format.Bytes())
			if marshalError != nil {
				return make([]flow.Message[[]byte, []byte], 0), ss, marshalError
			}
			byteResultMessages = append(byteResultMessages, bytesResMessage)
		}

		bytesNextState, stateMarshalError := ConvertSingleState(nextState, s, format.Bytes())
		if stateMarshalError != nil {
			return make([]flow.Message[[]byte, []byte], 0), ss, stateMarshalError
		}

		return byteResultMessages, bytesNextState, nil
	}
}

func ConvertTopicOneToTwo[S any, IK any, IV any, OK1 any, OV1 any, OK2 any, OV2 any](
	source OneToTwoFunction[S, IK, IV, OK1, OV1, OK2, OV2],
	s format.Format[S],
	inputTopic flow.Topic[IK, IV],
	outputTopic1 flow.Topic[OK1, OV1],
	outputTopic2 flow.Topic[OK2, OV2],
) SingleFunction {
	return func(ctx context.Context, m flow.Message[structure.Bytes, structure.Bytes], ss State[structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], State[structure.Bytes], error) {

		formattedMessage, unmarshalError := flow.Convert(m, format.Bytes(), format.Bytes(), inputTopic.KeyFormat(), inputTopic.ValueFormat())
		if unmarshalError != nil {
			return make([]flow.Message[[]byte, []byte], 0), ss, unmarshalError
		}

		formattedState, stateUnmarshalError := ConvertSingleState(ss, format.Bytes(), s)
		if stateUnmarshalError != nil {
			return make([]flow.Message[[]byte, []byte], 0), ss, stateUnmarshalError
		}

		res1, res2, nextState, fnError := source(ctx, formattedMessage, formattedState)
		if fnError != nil {
			return make([]flow.Message[[]byte, []byte], 0), ss, fnError
		}

		byteResultMessages := make([]flow.Message[[]byte, []byte], 0)

		if res1 != nil {
			bytesResMessage, marshalError := flow.Convert(*res1, outputTopic1.KeyFormat(), outputTopic1.ValueFormat(), format.Bytes(), format.Bytes())
			if marshalError != nil {
				return make([]flow.Message[[]byte, []byte], 0), ss, marshalError
			}
			bytesResMessage.Topic = outputTopic1.Name()
			byteResultMessages = append(byteResultMessages, bytesResMessage)
		}

		if res2 != nil {
			bytesResMessage, marshalError := flow.Convert(*res2, outputTopic2.KeyFormat(), outputTopic2.ValueFormat(), format.Bytes(), format.Bytes())
			if marshalError != nil {
				return make([]flow.Message[[]byte, []byte], 0), ss, marshalError
			}
			bytesResMessage.Topic = outputTopic2.Name()
			byteResultMessages = append(byteResultMessages, bytesResMessage)
		}

		bytesNextState, stateMarshalError := ConvertSingleState(nextState, s, format.Bytes())
		if stateMarshalError != nil {
			return make([]flow.Message[[]byte, []byte], 0), ss, stateMarshalError
		}

		return byteResultMessages, bytesNextState, nil
	}
}
