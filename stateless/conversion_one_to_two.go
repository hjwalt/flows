package stateless

import (
	"context"

	"github.com/hjwalt/flows/format"
	"github.com/hjwalt/flows/message"
)

func ConvertOneToTwo[IK any, IV any, OK1 any, OV1 any, OK2 any, OV2 any](
	source OneToTwoFunction[IK, IV, OK1, OV1, OK2, OV2],
	ik format.Format[IK],
	iv format.Format[IV],
	ok1 format.Format[OK1],
	ov1 format.Format[OV1],
	ok2 format.Format[OK2],
	ov2 format.Format[OV2],
) SingleFunction {
	return func(ctx context.Context, m message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {

		formattedMessage, unmarshalError := message.Convert(m, format.Bytes(), format.Bytes(), ik, iv)
		if unmarshalError != nil {
			return make([]message.Message[[]byte, []byte], 0), unmarshalError
		}

		res1, res2, fnError := source(ctx, formattedMessage)
		if fnError != nil {
			return make([]message.Message[[]byte, []byte], 0), fnError
		}

		byteResultMessages := make([]message.Message[[]byte, []byte], 0)

		if res1 != nil {
			bytesResMessage, marshalError := message.Convert(*res1, ok1, ov1, format.Bytes(), format.Bytes())
			if marshalError != nil {
				return make([]message.Message[[]byte, []byte], 0), marshalError
			}
			byteResultMessages = append(byteResultMessages, bytesResMessage)
		}

		if res2 != nil {
			bytesResMessage, marshalError := message.Convert(*res2, ok2, ov2, format.Bytes(), format.Bytes())
			if marshalError != nil {
				return make([]message.Message[[]byte, []byte], 0), marshalError
			}
			byteResultMessages = append(byteResultMessages, bytesResMessage)
		}

		return byteResultMessages, nil
	}
}
