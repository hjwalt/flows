package stateful_bun

import (
	"github.com/hjwalt/flows/format"
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/protobuf"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/runway/logger"
)

type SingleStateTable struct {
	PersistenceId      string `bun:",pk"`
	Internal           []byte
	Results            []byte
	Content            []byte
	CreatedTimestampMs int64
	UpdatedTimestampMs int64
}

var BytesFormat = format.Bytes()
var InternalStateProtoFormat = format.Protobuf[*protobuf.State]()
var ResultsStateProtoFormat = format.Protobuf[*protobuf.Results]()

func TableToState(dbState *SingleStateTable) (stateful.SingleState[message.Bytes], error) {

	internalValue, internalMapErr := format.Convert(dbState.Internal, BytesFormat, InternalStateProtoFormat)
	if internalMapErr != nil {
		logger.ErrorErr("single state mapping error", internalMapErr)
		return stateful.SingleState[message.Bytes]{}, internalMapErr
	}

	resultValue, resultMapErr := format.Convert(dbState.Results, BytesFormat, ResultsStateProtoFormat)
	if resultMapErr != nil {
		logger.ErrorErr("single state mapping error", resultMapErr)
		return stateful.SingleState[message.Bytes]{}, resultMapErr
	}

	return stateful.SingleState[message.Bytes]{
		PersistenceId:      dbState.PersistenceId,
		Internal:           internalValue,
		Results:            resultValue,
		Content:            dbState.Content,
		CreatedTimestampMs: dbState.CreatedTimestampMs,
		UpdatedTimestampMs: dbState.UpdatedTimestampMs,
	}, nil
}

func StateToTable(nextState stateful.SingleState[message.Bytes]) (*SingleStateTable, error) {
	internalBytes, internalMapErr := format.Convert(nextState.Internal, InternalStateProtoFormat, BytesFormat)
	if internalMapErr != nil {
		logger.ErrorErr("single state mapping error", internalMapErr)
		return &SingleStateTable{}, internalMapErr
	}

	resultBytes, resultMapErr := format.Convert(nextState.Results, ResultsStateProtoFormat, BytesFormat)
	if resultMapErr != nil {
		logger.ErrorErr("single state mapping error", resultMapErr)
		return &SingleStateTable{}, resultMapErr
	}

	return &SingleStateTable{
		PersistenceId:      nextState.PersistenceId,
		Internal:           internalBytes,
		Results:            resultBytes,
		Content:            nextState.Content,
		CreatedTimestampMs: nextState.CreatedTimestampMs,
		UpdatedTimestampMs: nextState.UpdatedTimestampMs,
	}, nil
}
