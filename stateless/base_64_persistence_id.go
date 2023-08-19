package stateless

import (
	"context"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/runway/format"
)

var (
	base64Format = format.Base64()
	bytesFormat  = format.Bytes()
)

func Base64PersistenceId(ctx context.Context, m message.Message[message.Bytes, message.Bytes]) (string, error) {
	base64Message, conversionErr := message.Convert(
		m,
		bytesFormat,
		bytesFormat,
		base64Format,
		bytesFormat,
	)

	if conversionErr != nil {
		return "", conversionErr
	}

	return base64Message.Key, nil
}
