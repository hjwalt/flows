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
func NewReadWrite(configurations ...runtime.Configuration[*ReadWrite]) stateless.BatchFunction {
	singleFunction := &ReadWrite{}
	for _, configuration := range configurations {
		singleFunction = configuration(singleFunction)
	}
	return singleFunction.Apply
}

// configuration
func WithReadWritePersistenceIdFunc(persistenceIdFunc func(context.Context, message.Message[message.Bytes, message.Bytes]) (string, error)) runtime.Configuration[*ReadWrite] {
	return func(st *ReadWrite) *ReadWrite {
		st.persistenceIdFunc = persistenceIdFunc
		return st
	}
}

func WithReadWriteRepository(repository Repository) runtime.Configuration[*ReadWrite] {
	return func(st *ReadWrite) *ReadWrite {
		st.repository = repository
		return st
	}
}

func WithReadWriteFunction(next SingleFunction) runtime.Configuration[*ReadWrite] {
	return func(st *ReadWrite) *ReadWrite {
		st.next = next
		return st
	}
}

// implementation
type ReadWrite struct {
	next              SingleFunction
	persistenceIdFunc PersistenceIdFunction[message.Bytes, message.Bytes]
	repository        Repository
}

func (r *ReadWrite) Apply(c context.Context, ms []message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {

	// prepare persistence id and message map
	persistenceIds := make([]string, len(ms))
	idMessageMap := map[string]message.Message[message.Bytes, message.Bytes]{}
	for i, m := range ms {
		persistenceId, persistenceIdErr := r.persistenceIdFunc(c, m)
		if persistenceIdErr != nil {
			return message.EmptySlice(), errors.Join(persistenceIdErr, ErrBatchReadWrite, ErrPersistenceId)
		}
		persistenceIds[i] = persistenceId

		idMessageMap[persistenceId] = m
	}

	// read states
	currentStates, readErr := r.repository.GetAll(c, persistenceIds)
	if readErr != nil {
		return message.EmptySlice(), errors.Join(readErr, ErrBatchReadWrite, ErrStateGet)
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
			return message.EmptySlice(), nextApplyErr
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
		return message.EmptySlice(), upsertErr
	}

	return resultMessages, nil
}

var (
	ErrBatchReadWrite = errors.New("batch read write")
	ErrPersistenceId  = errors.New("persistence id")
	ErrStateGet       = errors.New("state get")
)
