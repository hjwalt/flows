package stateless

import (
	"context"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/structure"
)

var (
	base64Format = format.Base64()
	bytesFormat  = format.Bytes()
)

func Base64PersistenceId(ctx context.Context, m flow.Message[structure.Bytes, structure.Bytes]) (string, error) {
	base64Message, conversionErr := flow.Convert(
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
