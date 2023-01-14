package stateful

import (
	"github.com/hjwalt/flows/format"
	"github.com/hjwalt/runway/logger"
)

// assuming byte compatibility, i.e. bytes <-> proto, string <-> json
func ConvertSingleState[V1 any, V2 any](
	source SingleState[V1],
	v1 format.Format[V1],
	v2 format.Format[V2],
) (SingleState[V2], error) {

	value, err := format.Convert(source.Content, v1, v2)
	if err != nil {
		logger.ErrorErr("message value conversion failure", err)
		return SingleState[V2]{}, err
	}

	return SingleState[V2]{
		PersistenceId:      source.PersistenceId,
		Internal:           source.Internal,
		Results:            source.Results,
		Content:            value,
		CreatedTimestampMs: source.CreatedTimestampMs,
		UpdatedTimestampMs: source.UpdatedTimestampMs,
	}, nil
}
