package stateful_test

import (
	"context"
	"testing"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/runway/structure"
	"github.com/stretchr/testify/assert"
)

type InMemoryRepository struct {
	State map[string]stateful.State[structure.Bytes]
}

func (r *InMemoryRepository) Get(ctx context.Context, persistenceId string) (stateful.State[structure.Bytes], error) {
	if state, statePresent := r.State[persistenceId]; statePresent {
		return state, nil
	} else {
		return stateful.NewState[[]byte](persistenceId, []byte{}), nil
	}
}

func (r *InMemoryRepository) GetAll(ctx context.Context, persistenceIds []string) (map[string]stateful.State[structure.Bytes], error) {
	stateMap := map[string]stateful.State[structure.Bytes]{}
	for _, persistenceId := range persistenceIds {
		if state, statePresent := r.State[persistenceId]; statePresent {
			stateMap[persistenceId] = state
		} else {
			stateMap[persistenceId] = stateful.NewState[[]byte](persistenceId, []byte{})
		}
	}
	return stateMap, nil
}

func (r *InMemoryRepository) Upsert(ctx context.Context, persistenceId string, dbState stateful.State[structure.Bytes]) error {
	r.State[persistenceId] = dbState
	return nil
}

func (r *InMemoryRepository) UpsertAll(ctx context.Context, stateMap map[string]stateful.State[structure.Bytes]) error {
	for k, v := range stateMap {
		r.State[k] = v
	}
	return nil
}

func TestBatchSimpleApplySuccessful(t *testing.T) {
	assert := assert.New(t)

	applyCount := 0

	repo := &InMemoryRepository{
		State: map[string]stateful.State[[]byte]{},
	}

	txFn := stateful.NewReadWrite(
		stateful.WithReadWritePersistenceIdFunc(func(ctx context.Context, m flow.Message[structure.Bytes, structure.Bytes]) (string, error) {
			return string(m.Key), nil
		}),
		stateful.WithReadWriteRepository(repo),
		stateful.WithReadWriteFunction(func(c context.Context, m flow.Message[structure.Bytes, structure.Bytes], inState stateful.State[structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], stateful.State[structure.Bytes], error) {
			applyCount += 1
			inState.Content = []byte(string(inState.Content) + string(m.Value))
			return []flow.Message[structure.Bytes, structure.Bytes]{
					{
						Key:   m.Key,
						Value: inState.Content,
					},
				},
				inState,
				nil
		}),
	)

	txFn(context.Background(), []flow.Message[[]byte, []byte]{
		{
			Key:   []byte("test"),
			Value: []byte("test"),
		},
		{
			Key:   []byte("test-2"),
			Value: []byte("test"),
		},
	})

	assert.Equal(2, applyCount)

	testState, stateErr := repo.Get(context.Background(), "test")
	assert.NoError(stateErr)
	assert.Equal("test", string(testState.Content))

	test2State, stateErr := repo.Get(context.Background(), "test-2")
	assert.NoError(stateErr)
	assert.Equal("test", string(test2State.Content))

	txFn(context.Background(), []flow.Message[[]byte, []byte]{
		{
			Key:   []byte("test"),
			Value: []byte("2"),
		},
		{
			Key:   []byte("test-2"),
			Value: []byte("2"),
		},
	})

	assert.Equal(4, applyCount)

	testState, stateErr = repo.Get(context.Background(), "test")
	assert.NoError(stateErr)
	assert.Equal("test2", string(testState.Content))

	test2State, stateErr = repo.Get(context.Background(), "test-2")
	assert.NoError(stateErr)
	assert.Equal("test2", string(test2State.Content))

	txFn(context.Background(), []flow.Message[[]byte, []byte]{
		{
			Key:   []byte("test"),
			Value: []byte("3"),
		},
	})

	assert.Equal(5, applyCount)

	testState, stateErr = repo.Get(context.Background(), "test")
	assert.NoError(stateErr)
	assert.Equal("test23", string(testState.Content))

	test2State, stateErr = repo.Get(context.Background(), "test-2")
	assert.NoError(stateErr)
	assert.Equal("test2", string(test2State.Content))

	txFn(context.Background(), []flow.Message[[]byte, []byte]{
		{
			Key:   []byte("test-2"),
			Value: []byte("3"),
		},
	})

	assert.Equal(6, applyCount)

	testState, stateErr = repo.Get(context.Background(), "test")
	assert.NoError(stateErr)
	assert.Equal("test23", string(testState.Content))

	test2State, stateErr = repo.Get(context.Background(), "test-2")
	assert.NoError(stateErr)
	assert.Equal("test23", string(test2State.Content))
}
