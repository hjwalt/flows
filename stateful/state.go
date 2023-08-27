package stateful

import (
	"context"

	"github.com/hjwalt/flows/protobuf"
	"github.com/hjwalt/runway/structure"
)

func NewState[S any](id string, content S) State[S] {
	singleState := State[S]{
		Id:      id,
		Content: content,
	}
	return SetDefault(singleState)
}

type State[S any] struct {
	Id                 string
	Internal           *protobuf.State
	Results            *protobuf.Results
	Content            S
	CreatedTimestampMs int64
	UpdatedTimestampMs int64
}

func SetDefault[S any](s State[S]) State[S] {
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

type Repository interface {
	Get(ctx context.Context, persistenceId string) (State[structure.Bytes], error)
	GetAll(ctx context.Context, persistenceId []string) (map[string]State[structure.Bytes], error)
	Upsert(ctx context.Context, persistenceId string, dbState State[structure.Bytes]) error
	UpsertAll(ctx context.Context, stateMap map[string]State[structure.Bytes]) error
}
