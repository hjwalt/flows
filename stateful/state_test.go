package stateful_test

import (
	"context"
	"testing"

	"github.com/hjwalt/flows/protobuf"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/runway/structure"
	"github.com/stretchr/testify/assert"
)

type RepositoryForTest struct {
	ProxiedGet       func(ctx context.Context, persistenceId string) (stateful.State[structure.Bytes], error)
	ProxiedGetAll    func(ctx context.Context, persistenceId []string) (map[string]stateful.State[structure.Bytes], error)
	ProxiedUpsert    func(ctx context.Context, persistenceId string, dbState stateful.State[structure.Bytes]) error
	ProxiedUpsertAll func(ctx context.Context, stateMap map[string]stateful.State[structure.Bytes]) error
}

func (r RepositoryForTest) Get(ctx context.Context, persistenceId string) (stateful.State[structure.Bytes], error) {
	return r.ProxiedGet(ctx, persistenceId)
}

func (r RepositoryForTest) GetAll(ctx context.Context, persistenceId []string) (map[string]stateful.State[structure.Bytes], error) {
	return r.ProxiedGetAll(ctx, persistenceId)
}

func (r RepositoryForTest) Upsert(ctx context.Context, persistenceId string, dbState stateful.State[structure.Bytes]) error {
	return r.ProxiedUpsert(ctx, persistenceId, dbState)
}

func (r RepositoryForTest) UpsertAll(ctx context.Context, stateMap map[string]stateful.State[structure.Bytes]) error {
	return r.ProxiedUpsertAll(ctx, stateMap)
}

func TestSingleStateDefault(t *testing.T) {
	// assert := assert.New(t)

	expected := stateful.State[structure.Bytes]{
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
		inputState stateful.State[structure.Bytes]
	}{
		{
			name:       "TestSingleStateDefaultEmpty",
			inputState: stateful.State[structure.Bytes]{},
		},
		{
			name: "TestSingleStateDefaultOnlyInternal",
			inputState: stateful.State[structure.Bytes]{
				Internal: &protobuf.State{},
			},
		},
		{
			name: "TestSingleStateDefaultOnlyInternalWithV1",
			inputState: stateful.State[structure.Bytes]{
				Internal: &protobuf.State{
					State: &protobuf.State_V1{},
				},
			},
		},
		{
			name: "TestSingleStateDefaultOnlyInternalWithV1Content",
			inputState: stateful.State[structure.Bytes]{
				Internal: &protobuf.State{
					State: &protobuf.State_V1{
						V1: &protobuf.StateV1{},
					},
				},
			},
		},
		{
			name: "TestSingleStateDefaultOnlyResults",
			inputState: stateful.State[structure.Bytes]{
				Results: &protobuf.Results{},
			},
		},
		{
			name: "TestSingleStateDefaultOnlyResultsWithV1",
			inputState: stateful.State[structure.Bytes]{
				Results: &protobuf.Results{
					Result: &protobuf.Results_V1{},
				},
			},
		},
		{
			name: "TestSingleStateDefaultOnlyResultsWithV1Content",
			inputState: stateful.State[structure.Bytes]{
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
