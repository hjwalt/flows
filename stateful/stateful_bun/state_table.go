package stateful_bun

import (
	"github.com/hjwalt/flows/protobuf"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/structure"
)

type StateTable struct {
	Id                 string `bun:",pk"`
	Internal           []byte
	Results            []byte
	Content            []byte
	CreatedTimestampMs int64
	UpdatedTimestampMs int64
}

var InternalStateProtoFormat = format.Protobuf[*protobuf.State]()
var ResultsStateProtoFormat = format.Protobuf[*protobuf.Results]()

func TableToState(dbState *StateTable) (stateful.State[structure.Bytes], error) {

	internalValue, internalMapErr := format.Convert(dbState.Internal, format.Bytes(), InternalStateProtoFormat)
	if internalMapErr != nil {
		logger.ErrorErr("single state mapping error", internalMapErr)
		return stateful.State[structure.Bytes]{}, internalMapErr
	}

	resultValue, resultMapErr := format.Convert(dbState.Results, format.Bytes(), ResultsStateProtoFormat)
	if resultMapErr != nil {
		logger.ErrorErr("single state mapping error", resultMapErr)
		return stateful.State[structure.Bytes]{}, resultMapErr
	}

	return stateful.State[structure.Bytes]{
		Id:                 dbState.Id,
		Internal:           internalValue,
		Results:            resultValue,
		Content:            dbState.Content,
		CreatedTimestampMs: dbState.CreatedTimestampMs,
		UpdatedTimestampMs: dbState.UpdatedTimestampMs,
	}, nil
}

func StateToTable(nextState stateful.State[structure.Bytes]) (*StateTable, error) {
	internalBytes, internalMapErr := format.Convert(nextState.Internal, InternalStateProtoFormat, format.Bytes())
	if internalMapErr != nil {
		logger.ErrorErr("single state mapping error", internalMapErr)
		return &StateTable{}, internalMapErr
	}

	resultBytes, resultMapErr := format.Convert(nextState.Results, ResultsStateProtoFormat, format.Bytes())
	if resultMapErr != nil {
		logger.ErrorErr("single state mapping error", resultMapErr)
		return &StateTable{}, resultMapErr
	}

	return &StateTable{
		Id:                 nextState.Id,
		Internal:           internalBytes,
		Results:            resultBytes,
		Content:            nextState.Content,
		CreatedTimestampMs: nextState.CreatedTimestampMs,
		UpdatedTimestampMs: nextState.UpdatedTimestampMs,
	}, nil
}
