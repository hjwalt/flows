package stateful

import (
	"context"
	"errors"
	"time"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/runtime"
)

// constructor
func NewBatchReadWrite(configurations ...runtime.Configuration[*BatchReadWrite]) stateless.BatchFunction {
	singleFunction := &BatchReadWrite{}
	for _, configuration := range configurations {
		singleFunction = configuration(singleFunction)
	}
	return singleFunction.Apply
}

// configuration
func WithBatchReadWritePersistenceIdFunc(persistenceIdFunc func(context.Context, message.Message[message.Bytes, message.Bytes]) (string, error)) runtime.Configuration[*BatchReadWrite] {
	return func(st *BatchReadWrite) *BatchReadWrite {
		st.persistenceIdFunc = persistenceIdFunc
		return st
	}
}

func WithBatchReadWriteRepository(repository Repository) runtime.Configuration[*BatchReadWrite] {
	return func(st *BatchReadWrite) *BatchReadWrite {
		st.repository = repository
		return st
	}
}

func WithBatchReadWriteStatefulFunction(next SingleFunction) runtime.Configuration[*BatchReadWrite] {
	return func(st *BatchReadWrite) *BatchReadWrite {
		st.next = next
		return st
	}
}

// implementation
type BatchReadWrite struct {
	next              SingleFunction
	persistenceIdFunc func(context.Context, message.Message[message.Bytes, message.Bytes]) (string, error)
	repository        Repository
}

func (r *BatchReadWrite) Apply(c context.Context, ms []message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {

	// prepare persistence id and message map
	persistenceIds := make([]string, len(ms))
	idMessageMap := map[string]message.Message[message.Bytes, message.Bytes]{}
	for i, m := range ms {
		persistenceId, persistenceIdErr := r.persistenceIdFunc(c, m)
		if persistenceIdErr != nil {
			return make([]message.Message[[]byte, []byte], 0), errors.Join(persistenceIdErr, ErrBatchReadWrite, ErrPersistenceId)
		}
		persistenceIds[i] = persistenceId

		idMessageMap[persistenceId] = m
	}

	// read states
	currentStates, readErr := r.repository.GetAll(c, persistenceIds)
	if readErr != nil {
		return make([]message.Message[[]byte, []byte], 0), errors.Join(readErr, ErrBatchReadWrite, ErrStateGet)
	}

	// prepare result holder
	nextStateMap := map[string]State[[]byte]{}
	resultMessages := []message.Message[message.Bytes, message.Bytes]{}

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
			return make([]message.Message[[]byte, []byte], 0), nextApplyErr
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
		return make([]message.Message[[]byte, []byte], 0), upsertErr
	}

	return resultMessages, nil
}

var (
	ErrBatchReadWrite = errors.New("batch read write")
	ErrPersistenceId  = errors.New("persistence id")
	ErrStateGet       = errors.New("state get")
)
