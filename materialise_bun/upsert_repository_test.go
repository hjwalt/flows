package materialise_bun_test

import (
	"context"
	"testing"

	"github.com/hjwalt/flows/materialise_bun"
	"github.com/hjwalt/flows/mock"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
)

type TestEntity struct {
	bun.BaseModel `bun:"table:test_entity"`
	Id            string `bun:",pk"`
	KeyContent    string
	ValueContent  string
	TimestampMs   int64
}

func TestUpsertStandard(t *testing.T) {
	assert := assert.New(t)

	conn := &mock.SqliteBunConnection[*TestEntity]{
		FixtureFolder: "fixture",
		FixtureFile:   "test_entity.yaml",
	}
	err := conn.Start()
	assert.NoError(err)

	repo := materialise_bun.NewBunUpsertRepository(
		materialise_bun.WithBunUpsertRepositoryConnection[*TestEntity](conn),
	)

	ctx := context.Background()

	repo.Upsert(ctx, []*TestEntity{
		{
			Id:           "1",
			KeyContent:   "key1",
			ValueContent: "value1",
		},
		{
			Id:           "3",
			KeyContent:   "key3",
			ValueContent: "value3",
		},
	})

	rows := make([]TestEntity, 0)
	conn.Db().NewSelect().Model(&rows).Scan(ctx)

	assert.Equal(3, len(rows))

	assert.Equal("1", rows[0].Id)
	assert.Equal("key1", rows[0].KeyContent)
	assert.Equal("value1", rows[0].ValueContent)

	assert.Equal("2", rows[1].Id)
	assert.Equal("test 2", rows[1].KeyContent)
	assert.Equal("test 2", rows[1].ValueContent)

	assert.Equal("3", rows[2].Id)
	assert.Equal("key3", rows[2].KeyContent)
	assert.Equal("value3", rows[2].ValueContent)

	conn.Stop()
}
