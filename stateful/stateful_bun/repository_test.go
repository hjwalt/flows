package stateful_bun_test

import (
	"context"
	"testing"

	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateful/stateful_bun"
	"github.com/hjwalt/flows/test_helper"
	"github.com/stretchr/testify/assert"
)

func TestStateRepositoryGet(t *testing.T) {
	assert := assert.New(t)

	conn := &test_helper.SqliteBunConnection[*stateful_bun.StateTable]{
		FixtureFolder: "fixture",
		FixtureFile:   "state_table.yaml",
	}
	err := conn.Start()
	assert.NoError(err)

	repo := stateful_bun.NewRepository(
		stateful_bun.WithConnection(conn),
		stateful_bun.WithStateTableName("state_tables"),
	)

	state, getErr := repo.Get(context.Background(), "test")
	assert.NoError(getErr)
	assert.Equal("test", state.Id)
	assert.Equal(int64(100), state.CreatedTimestampMs)

	state, getErr = repo.Get(context.Background(), "test2")
	assert.NoError(getErr)
	assert.Equal("test2", state.Id)
	assert.Equal(int64(0), state.CreatedTimestampMs)

	conn.Stop()
}

func TestStateRepositoryGetAll(t *testing.T) {
	assert := assert.New(t)

	conn := &test_helper.SqliteBunConnection[*stateful_bun.StateTable]{
		FixtureFolder: "fixture",
		FixtureFile:   "state_table.yaml",
	}
	err := conn.Start()
	assert.NoError(err)

	repo := stateful_bun.NewRepository(
		stateful_bun.WithConnection(conn),
		stateful_bun.WithStateTableName("state_tables"),
	)

	states, getErr := repo.GetAll(context.Background(), []string{"test", "test-another", "test-2"})
	assert.NoError(getErr)
	assert.Equal(3, len(states))

	state := states["test"]
	assert.Equal("test", state.Id)
	assert.Equal(int64(100), state.CreatedTimestampMs)

	state = states["test-another"]
	assert.Equal("test-another", state.Id)
	assert.Equal(int64(200), state.CreatedTimestampMs)

	state = states["test-2"]
	assert.Equal("test-2", state.Id)
	assert.Equal(int64(0), state.CreatedTimestampMs)

	conn.Stop()
}

func TestStateRepositoryUpsert(t *testing.T) {

	assert := assert.New(t)

	conn := &test_helper.SqliteBunConnection[*stateful_bun.StateTable]{
		FixtureFolder: "fixture",
		FixtureFile:   "state_table.yaml",
	}
	err := conn.Start()
	assert.NoError(err)

	repo := stateful_bun.NewRepository(
		stateful_bun.WithConnection(conn),
		stateful_bun.WithStateTableName("state_tables"),
	)

	upsertErr := repo.Upsert(context.Background(), "test", stateful.State[[]byte]{
		Id:                 "test",
		Content:            []byte("content"),
		CreatedTimestampMs: 200,
	})
	assert.NoError(upsertErr)

	state, getErr := repo.Get(context.Background(), "test")
	assert.NoError(getErr)
	assert.Equal("test", state.Id)
	assert.Equal("content", string(state.Content))
	assert.Equal(int64(200), state.CreatedTimestampMs)

	upsertErr = repo.Upsert(context.Background(), "test", stateful.State[[]byte]{
		Id:                 "test2",
		Content:            []byte("content2"),
		CreatedTimestampMs: 300,
	})
	assert.NoError(upsertErr)

	state, getErr = repo.Get(context.Background(), "test")
	assert.NoError(getErr)
	assert.Equal("test", state.Id)
	assert.Equal("content2", string(state.Content))
	assert.Equal(int64(300), state.CreatedTimestampMs)

	state, getErr = repo.Get(context.Background(), "test2")
	assert.NoError(getErr)
	assert.Equal("test2", state.Id)
	assert.Equal(int64(0), state.CreatedTimestampMs)

	conn.Stop()
}

func TestStateRepositoryUpsertAll(t *testing.T) {

	assert := assert.New(t)

	conn := &test_helper.SqliteBunConnection[*stateful_bun.StateTable]{
		FixtureFolder: "fixture",
		FixtureFile:   "state_table.yaml",
	}
	err := conn.Start()
	assert.NoError(err)

	repo := stateful_bun.NewRepository(
		stateful_bun.WithConnection(conn),
		stateful_bun.WithStateTableName("state_tables"),
	)

	upsertErr := repo.UpsertAll(
		context.Background(),
		map[string]stateful.State[[]byte]{
			"test": {
				Id:                 "test",
				Content:            []byte("content"),
				CreatedTimestampMs: 200,
			},
			"test-another": {
				Id:                 "test-another",
				Content:            []byte("content"),
				CreatedTimestampMs: 400,
			},
		},
	)

	assert.NoError(upsertErr)

	states, getErr := repo.GetAll(context.Background(), []string{"test", "test-another", "test-2"})
	assert.NoError(getErr)
	assert.Equal(3, len(states))

	state := states["test"]
	assert.Equal("test", state.Id)
	assert.Equal("content", string(state.Content))
	assert.Equal(int64(200), state.CreatedTimestampMs)

	state = states["test-another"]
	assert.Equal("test-another", state.Id)
	assert.Equal("content", string(state.Content))
	assert.Equal(int64(400), state.CreatedTimestampMs)

	state = states["test-2"]
	assert.Equal("test-2", state.Id)
	assert.Equal(int64(0), state.CreatedTimestampMs)

	conn.Stop()
}
