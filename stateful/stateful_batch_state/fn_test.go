package stateful_batch_state_test

import (
	"context"
	"testing"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateful/stateful_batch_state"
	"github.com/hjwalt/flows/stateful/stateful_mock"
	"github.com/hjwalt/runway/structure"
	"github.com/stretchr/testify/assert"
)

func TestFn(t *testing.T) {
	assert := assert.New(t)

	applyCount := 0

	repo := stateful_mock.Repository()

	txFn := stateful_batch_state.New(
		func(c context.Context, m flow.Message[structure.Bytes, structure.Bytes], inState stateful.State[structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], stateful.State[structure.Bytes], error) {
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
		},
		func(ctx context.Context, m flow.Message[structure.Bytes, structure.Bytes]) (string, error) {
			return string(m.Key), nil
		},
		repo,
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
