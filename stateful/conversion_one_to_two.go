package stateful

import (
	"context"

	"github.com/hjwalt/flows/format"
	"github.com/hjwalt/flows/message"
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
	return func(ctx context.Context, m message.Message[message.Bytes, message.Bytes], ss SingleState[message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], SingleState[message.Bytes], error) {

		formattedMessage, unmarshalError := message.Convert(m, format.Bytes(), format.Bytes(), ik, iv)
		if unmarshalError != nil {
			return make([]message.Message[[]byte, []byte], 0), ss, unmarshalError
		}

		formattedState, stateUnmarshalError := ConvertSingleState(ss, format.Bytes(), s)
		if stateUnmarshalError != nil {
			return make([]message.Message[[]byte, []byte], 0), ss, stateUnmarshalError
		}

		res1, res2, nextState, fnError := source(ctx, formattedMessage, formattedState)
		if fnError != nil {
			return make([]message.Message[[]byte, []byte], 0), ss, fnError
		}

		byteResultMessages := make([]message.Message[[]byte, []byte], 0)

		if res1 != nil {
			bytesResMessage, marshalError := message.Convert(*res1, ok1, ov1, format.Bytes(), format.Bytes())
			if marshalError != nil {
				return make([]message.Message[[]byte, []byte], 0), ss, marshalError
			}
			byteResultMessages = append(byteResultMessages, bytesResMessage)
		}

		if res2 != nil {
			bytesResMessage, marshalError := message.Convert(*res2, ok2, ov2, format.Bytes(), format.Bytes())
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
