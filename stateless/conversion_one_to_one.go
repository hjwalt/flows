package stateless

import (
	"context"

	"github.com/hjwalt/flows/format"
	"github.com/hjwalt/flows/message"
)

func ConvertOneToOne[IK any, IV any, OK any, OV any](
	source OneToOneFunction[IK, IV, OK, OV],
	ik format.Format[IK],
	iv format.Format[IV],
	ok format.Format[OK],
	ov format.Format[OV],
) SingleFunction {
	return func(ctx context.Context, m message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
		formattedMessage, unmarshalError := message.Convert(m, format.Bytes(), format.Bytes(), ik, iv)
		if unmarshalError != nil {
			return make([]message.Message[[]byte, []byte], 0), unmarshalError
		}

		res, fnError := source(ctx, formattedMessage)
		if fnError != nil {
			return make([]message.Message[[]byte, []byte], 0), fnError
		}

		byteResultMessages := make([]message.Message[[]byte, []byte], 0)

		if res != nil {
			bytesResMessage, marshalError := message.Convert(*res, ok, ov, format.Bytes(), format.Bytes())
			if marshalError != nil {
				return make([]message.Message[[]byte, []byte], 0), marshalError
			}
			byteResultMessages = append(byteResultMessages, bytesResMessage)
		}

		return byteResultMessages, nil
	}
}
