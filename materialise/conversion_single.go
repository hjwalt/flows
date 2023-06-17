package materialise

import (
	"context"

	"github.com/hjwalt/flows/format"
	"github.com/hjwalt/flows/message"
)

func ConvertOneToOne[IK any, IV any, T any](
	source MapFunction[IK, IV, T],
	ik format.Format[IK],
	iv format.Format[IV],
) MapFunction[message.Bytes, message.Bytes, T] {
	return func(ctx context.Context, m message.Message[[]byte, []byte]) ([]T, error) {
		formattedMessage, unmarshalError := message.Convert(m, format.Bytes(), format.Bytes(), ik, iv)
		if unmarshalError != nil {
			return make([]T, 0), unmarshalError
		}
		return source(ctx, formattedMessage)
	}
}
