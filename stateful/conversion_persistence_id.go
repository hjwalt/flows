package stateful

import (
	"context"

	"github.com/hjwalt/flows/format"
	"github.com/hjwalt/flows/message"
)

func ConvertPersistenceId[IK any, IV any](
	source PersistenceIdFunction[IK, IV],
	ik format.Format[IK],
	iv format.Format[IV],
) PersistenceIdFunction[[]byte, []byte] {
	return func(ctx context.Context, m message.Message[[]byte, []byte]) (string, error) {
		formattedMessage, unmarshalError := message.Convert(m, format.Bytes(), format.Bytes(), ik, iv)
		if unmarshalError != nil {
			return "", unmarshalError
		}

		return source(ctx, formattedMessage)
	}
}
