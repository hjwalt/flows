package stateful_test

import (
	"context"
	"testing"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/protobuf"
	"github.com/hjwalt/flows/stateful"
	"github.com/stretchr/testify/assert"
)

type SingleStateRepositoryForTest struct {
	ProxiedGet    func(ctx context.Context, persistenceId string) (stateful.SingleState[message.Bytes], error)
	ProxiedUpsert func(ctx context.Context, persistenceId string, dbState stateful.SingleState[message.Bytes]) error
}

func (r SingleStateRepositoryForTest) Get(ctx context.Context, persistenceId string) (stateful.SingleState[message.Bytes], error) {
	return r.ProxiedGet(ctx, persistenceId)
}

func (r SingleStateRepositoryForTest) Upsert(ctx context.Context, persistenceId string, dbState stateful.SingleState[message.Bytes]) error {
	return r.ProxiedUpsert(ctx, persistenceId, dbState)
}

func TestSingleStateDefault(t *testing.T) {
	// assert := assert.New(t)

	expected := stateful.SingleState[message.Bytes]{
		Internal: &protobuf.State{
			State: &protobuf.State_V1{
				V1: &protobuf.StateV1{
					OffsetProgress: map[int32]int64{},
				},
			},
		},
		Results: &protobuf.Results{
			Result: &protobuf.Results_V1{
				V1: &protobuf.ResultV1{
					Messages: make([][]byte, 0),
				},
			},
		},
	}

	cases := []struct {
		name       string
		inputState stateful.SingleState[message.Bytes]
	}{
		{
			name:       "TestSingleStateDefaultEmpty",
			inputState: stateful.SingleState[message.Bytes]{},
		},
		{
			name: "TestSingleStateDefaultOnlyInternal",
			inputState: stateful.SingleState[message.Bytes]{
				Internal: &protobuf.State{},
			},
		},
		{
			name: "TestSingleStateDefaultOnlyInternalWithV1",
			inputState: stateful.SingleState[message.Bytes]{
				Internal: &protobuf.State{
					State: &protobuf.State_V1{},
				},
			},
		},
		{
			name: "TestSingleStateDefaultOnlyInternalWithV1Content",
			inputState: stateful.SingleState[message.Bytes]{
				Internal: &protobuf.State{
					State: &protobuf.State_V1{
						V1: &protobuf.StateV1{},
					},
				},
			},
		},
		{
			name: "TestSingleStateDefaultOnlyResults",
			inputState: stateful.SingleState[message.Bytes]{
				Results: &protobuf.Results{},
			},
		},
		{
			name: "TestSingleStateDefaultOnlyResultsWithV1",
			inputState: stateful.SingleState[message.Bytes]{
				Results: &protobuf.Results{
					Result: &protobuf.Results_V1{},
				},
			},
		},
		{
			name: "TestSingleStateDefaultOnlyResultsWithV1Content",
			inputState: stateful.SingleState[message.Bytes]{
				Results: &protobuf.Results{
					Result: &protobuf.Results_V1{
						V1: &protobuf.ResultV1{},
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert := assert.New(t)
			defaultSet := stateful.SetDefault(c.inputState)
			assert.Equal(expected, defaultSet)
		})
	}

}
