package stateful

import (
	"context"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/protobuf"
)

func NewSingleState[S any](content S) SingleState[S] {
	singleState := SingleState[S]{
		Content: content,
	}
	return SetDefault(singleState)
}

type SingleState[S any] struct {
	PersistenceId      string
	Internal           *protobuf.State
	Results            *protobuf.Results
	Content            S
	CreatedTimestampMs int64
	UpdatedTimestampMs int64
}

func SetDefault[S any](s SingleState[S]) SingleState[S] {
	stateToUse := s

	if stateToUse.Internal == nil {
		stateToUse.Internal = &protobuf.State{}
	}
	if stateToUse.Internal.State == nil {
		stateToUse.Internal.State = &protobuf.State_V1{
			V1: &protobuf.StateV1{},
		}
	}
	if stateToUse.Internal.GetV1() == nil {
		stateToUse.Internal.State = &protobuf.State_V1{
			V1: &protobuf.StateV1{},
		}
	}
	if stateToUse.Internal.GetV1().OffsetProgress == nil {
		stateToUse.Internal.GetV1().OffsetProgress = map[int32]int64{}
	}

	if stateToUse.Results == nil {
		stateToUse.Results = &protobuf.Results{}
	}
	if stateToUse.Results.Result == nil {
		stateToUse.Results.Result = &protobuf.Results_V1{
			V1: &protobuf.ResultV1{},
		}
	}
	if stateToUse.Results.GetV1() == nil {
		stateToUse.Results.Result = &protobuf.Results_V1{
			V1: &protobuf.ResultV1{},
		}
	}
	if stateToUse.Results.GetV1().Messages == nil {
		stateToUse.Results.GetV1().Messages = make([][]byte, 0)
	}

	return stateToUse
}

type SingleStateRepository interface {
	Get(ctx context.Context, persistenceId string) (SingleState[message.Bytes], error)
	Upsert(ctx context.Context, persistenceId string, dbState SingleState[message.Bytes]) error
}
