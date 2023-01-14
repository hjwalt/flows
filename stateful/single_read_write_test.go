package stateful_test

import (
	"context"
	"testing"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/stateful"
	"github.com/stretchr/testify/assert"
)

func TestSimpleApplySuccessful(t *testing.T) {
	assert := assert.New(t)

	getCount := 0
	upsertCount := 0
	applyCount := 0

	repo := SingleStateRepositoryForTest{
		ProxiedGet: func(ctx context.Context, persistenceId string) (stateful.SingleState[message.Bytes], error) {
			getCount += 1
			return stateful.SingleState[[]byte]{}, nil
		},
		ProxiedUpsert: func(ctx context.Context, persistenceId string, dbState stateful.SingleState[message.Bytes]) error {
			upsertCount += 1
			return nil
		},
	}

	txFn := stateful.NewSingleReadWrite(
		stateful.WithSingleReadWriteTransactionPersistenceIdFunc(func(ctx context.Context, m message.Message[message.Bytes, message.Bytes]) (string, error) {
			return string(m.Key), nil
		}),
		stateful.WithSingleReadWriteRepository(repo),
		stateful.WithSingleReadWriteStatefulFunction(func(c context.Context, m message.Message[message.Bytes, message.Bytes], inState stateful.SingleState[message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], stateful.SingleState[message.Bytes], error) {
			applyCount += 1
			return nil, inState, nil
		}),
	)

	txFn(context.Background(), message.Message[[]byte, []byte]{})

	assert.Equal(1, getCount)
	assert.Equal(1, upsertCount)
	assert.Equal(1, applyCount)
}
