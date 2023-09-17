package stateful_batch_state

import (
	"context"
	"errors"
	"time"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/structure"
)

func New(
	next stateful.SingleFunction,
	keyfn stateful.PersistenceIdFunction[structure.Bytes, structure.Bytes],
	repo stateful.Repository,
) stateless.BatchFunction {
	f := fn{next: next, persistenceIdFunc: keyfn, repository: repo}
	return f.apply
}

// implementation
type fn struct {
	next              stateful.SingleFunction
	persistenceIdFunc stateful.PersistenceIdFunction[structure.Bytes, structure.Bytes]
	repository        stateful.Repository
}

func (r *fn) apply(c context.Context, ms []flow.Message[structure.Bytes, structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], error) {

	// prepare persistence id and message map
	persistenceIds := make([]string, len(ms))
	idMessageMap := map[string]flow.Message[structure.Bytes, structure.Bytes]{}
	for i, m := range ms {
		persistenceId, persistenceIdErr := r.persistenceIdFunc(c, m)
		if persistenceIdErr != nil {
			return flow.EmptySlice(), errors.Join(persistenceIdErr, ErrBatchReadWrite, ErrPersistenceId)
		}
		persistenceIds[i] = persistenceId

		idMessageMap[persistenceId] = m
	}

	// read states
	currentStates, readErr := r.repository.GetAll(c, persistenceIds)
	if readErr != nil {
		return flow.EmptySlice(), errors.Join(readErr, ErrBatchReadWrite, ErrStateGet)
	}

	// prepare result holder
	nextStateMap := map[string]stateful.State[[]byte]{}
	resultMessages := []flow.Message[structure.Bytes, structure.Bytes]{}

	// execute next function, create next state map
	for persistenceId, currentState := range currentStates {
		m := idMessageMap[persistenceId]

		// set default value for new persistence id
		currentState.Id = persistenceId
		if currentState.CreatedTimestampMs == 0 {
			currentState.CreatedTimestampMs = time.Now().UTC().UnixMilli()
		}

		nextMessages, nextState, nextApplyErr := r.next(c, m, currentState)
		if nextApplyErr != nil {
			logger.ErrorErr("bun state last result mapping error", nextApplyErr)
			return flow.EmptySlice(), nextApplyErr
		}

		// set next updated value
		nextState.UpdatedTimestampMs = time.Now().UTC().UnixMilli()

		nextStateMap[persistenceId] = nextState
		if len(nextMessages) > 0 {
			resultMessages = append(resultMessages, nextMessages...)
		}
	}

	// upsert state
	upsertErr := r.repository.UpsertAll(c, nextStateMap)
	if upsertErr != nil {
		logger.ErrorErr("bun state upsert error", upsertErr)
		return flow.EmptySlice(), upsertErr
	}

	return resultMessages, nil
}

var (
	ErrBatchReadWrite = errors.New("batch read write")
	ErrPersistenceId  = errors.New("persistence id")
	ErrStateGet       = errors.New("state get")
)
