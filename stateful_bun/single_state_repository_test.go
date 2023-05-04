package stateful_bun_test

import (
	"context"
	"testing"

	"github.com/hjwalt/flows/mock"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateful_bun"
	"github.com/stretchr/testify/assert"
)

func TestSingleStateRepositoryGet(t *testing.T) {
	assert := assert.New(t)

	conn := &mock.SqliteBunConnection[*stateful_bun.SingleStateTable]{
		FixtureFolder: "fixture",
		FixtureFile:   "single_state_table.yaml",
	}
	err := conn.Start()
	assert.NoError(err)

	repo := stateful_bun.NewSingleStateRepository(
		stateful_bun.WithSingleStateRepositoryConnection(conn),
		stateful_bun.WithSingleStateRepositoryPersistenceTableName("single_state_tables"),
	)

	state, getErr := repo.Get(context.Background(), "test")
	assert.NoError(getErr)
	assert.Equal("test", state.PersistenceId)
	assert.Equal(int64(100), state.CreatedTimestampMs)

	state, getErr = repo.Get(context.Background(), "test2")
	assert.NoError(getErr)
	assert.Equal("", state.PersistenceId)

	conn.Stop()
}

func TestSingleStateRepositoryUpsert(t *testing.T) {

	assert := assert.New(t)

	conn := &mock.SqliteBunConnection[*stateful_bun.SingleStateTable]{
		FixtureFolder: "fixture",
		FixtureFile:   "single_state_table.yaml",
	}
	err := conn.Start()
	assert.NoError(err)

	repo := stateful_bun.NewSingleStateRepository(
		stateful_bun.WithSingleStateRepositoryConnection(conn),
		stateful_bun.WithSingleStateRepositoryPersistenceTableName("single_state_tables"),
	)

	upsertErr := repo.Upsert(context.Background(), "test", stateful.SingleState[[]byte]{
		PersistenceId:      "test",
		Content:            []byte("content"),
		CreatedTimestampMs: 200,
	})
	assert.NoError(upsertErr)

	state, getErr := repo.Get(context.Background(), "test")
	assert.NoError(getErr)
	assert.Equal("test", state.PersistenceId)
	assert.Equal("content", string(state.Content))
	assert.Equal(int64(200), state.CreatedTimestampMs)

	upsertErr = repo.Upsert(context.Background(), "test", stateful.SingleState[[]byte]{
		PersistenceId:      "test2",
		Content:            []byte("content2"),
		CreatedTimestampMs: 300,
	})
	assert.NoError(upsertErr)

	state, getErr = repo.Get(context.Background(), "test2")
	assert.NoError(getErr)
	assert.Equal("test2", state.PersistenceId)
	assert.Equal("content2", string(state.Content))
	assert.Equal(int64(300), state.CreatedTimestampMs)

	conn.Stop()
}
