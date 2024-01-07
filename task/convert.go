package task

import (
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/logger"
)

// assuming byte compatibility, i.e. bytes <-> proto, string <-> json
func Convert[V1 any, V2 any](
	source Message[V1],
	v1 format.Format[V1],
	v2 format.Format[V2],
) (Message[V2], error) {
	value, err := format.Convert(source.Value, v1, v2)
	if err != nil {
		logger.ErrorErr("message value conversion failure", err)
		return Message[V2]{}, err
	}

	return Message[V2]{
		Value:     value,
		Headers:   source.Headers,
		Timestamp: source.Timestamp,
	}, nil
}
