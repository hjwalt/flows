package stateful

import (
	"context"

	"github.com/hjwalt/flows/format"
	"github.com/hjwalt/flows/message"
)

var BytesFormat = format.Bytes()

func ConvertOneToOne[S any, IK any, IV any, OK any, OV any](
	source OneToOneFunction[S, IK, IV, OK, OV],
	s format.Format[S],
	ik format.Format[IK],
	iv format.Format[IV],
	ok format.Format[OK],
	ov format.Format[OV],
) SingleFunction {
	return func(ctx context.Context, m message.Message[message.Bytes, message.Bytes], ss SingleState[message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], SingleState[message.Bytes], error) {

		formattedMessage, unmarshalError := message.Convert(m, BytesFormat, BytesFormat, ik, iv)
		if unmarshalError != nil {
			return make([]message.Message[[]byte, []byte], 0), ss, unmarshalError
		}

		formattedState, stateUnmarshalError := ConvertSingleState(ss, BytesFormat, s)
		if stateUnmarshalError != nil {
			return make([]message.Message[[]byte, []byte], 0), ss, stateUnmarshalError
		}

		res, nextState, fnError := source(ctx, formattedMessage, formattedState)
		if fnError != nil {
			return make([]message.Message[[]byte, []byte], 0), ss, fnError
		}

		byteResultMessages := make([]message.Message[[]byte, []byte], 0)

		if res != nil {
			bytesResMessage, marshalError := message.Convert(*res, ok, ov, BytesFormat, BytesFormat)
			if marshalError != nil {
				return make([]message.Message[[]byte, []byte], 0), ss, marshalError
			}
			byteResultMessages = append(byteResultMessages, bytesResMessage)
		}

		bytesNextState, stateMarshalError := ConvertSingleState(nextState, s, BytesFormat)
		if stateMarshalError != nil {
			return make([]message.Message[[]byte, []byte], 0), ss, stateMarshalError
		}

		return byteResultMessages, bytesNextState, nil
	}
}
