package message

import (
	"github.com/hjwalt/flows/format"
	"github.com/hjwalt/runway/logger"
)

// assuming byte compatibility, i.e. bytes <-> proto, string <-> json
func Convert[K1 any, V1 any, K2 any, V2 any](
	source Message[K1, V1],
	k1 format.Format[K1],
	v1 format.Format[V1],
	k2 format.Format[K2],
	v2 format.Format[V2],
) (Message[K2, V2], error) {

	key, err := format.Convert(source.Key, k1, k2)
	if err != nil {
		logger.ErrorErr("message key conversion failure", err)
		return Message[K2, V2]{}, err
	}

	value, err := format.Convert(source.Value, v1, v2)
	if err != nil {
		logger.ErrorErr("message value conversion failure", err)
		return Message[K2, V2]{}, err
	}

	return Message[K2, V2]{
		Topic:     source.Topic,
		Partition: source.Partition,
		Offset:    source.Offset,
		Timestamp: source.Timestamp,
		Key:       key,
		Value:     value,
		Headers:   source.Headers,
	}, nil
}
