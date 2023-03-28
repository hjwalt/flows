package stateful

import (
	"context"
	"time"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/logger"
)

const (
	ContextTransaction = "ContextTransaction"
)

// constructor
func NewSingleReadWrite(configurations ...runtime.Configuration[*SingleStateReadWrite]) stateless.SingleFunction {
	singleFunction := &SingleStateReadWrite{}
	for _, configuration := range configurations {
		singleFunction = configuration(singleFunction)
	}
	return singleFunction.Apply
}

// configuration
func WithSingleReadWriteTransactionPersistenceIdFunc(persistenceIdFunc func(context.Context, message.Message[message.Bytes, message.Bytes]) (string, error)) runtime.Configuration[*SingleStateReadWrite] {
	return func(st *SingleStateReadWrite) *SingleStateReadWrite {
		st.persistenceIdFunc = persistenceIdFunc
		return st
	}
}

func WithSingleReadWriteRepository(repository SingleStateRepository) runtime.Configuration[*SingleStateReadWrite] {
	return func(st *SingleStateReadWrite) *SingleStateReadWrite {
		st.repository = repository
		return st
	}
}

func WithSingleReadWriteStatefulFunction(next SingleFunction) runtime.Configuration[*SingleStateReadWrite] {
	return func(st *SingleStateReadWrite) *SingleStateReadWrite {
		st.next = next
		return st
	}
}

// set

// implementation
type SingleStateReadWrite struct {
	next              SingleFunction
	persistenceIdFunc func(context.Context, message.Message[message.Bytes, message.Bytes]) (string, error)
	repository        SingleStateRepository
}

func (r *SingleStateReadWrite) Apply(c context.Context, m message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {

	persistenceId, persistenceIdErr := r.persistenceIdFunc(c, m)
	if persistenceIdErr != nil {
		logger.ErrorErr("bun get persistence id error", persistenceIdErr)
		return make([]message.Message[[]byte, []byte], 0), persistenceIdErr
	}

	// read states
	currentState, readErr := r.repository.Get(c, persistenceId)
	if readErr != nil {
		logger.ErrorErr("bun state reading error", readErr)
		return make([]message.Message[[]byte, []byte], 0), readErr
	}

	// set default value for new persistence id
	if len(currentState.PersistenceId) == 0 {
		currentState.PersistenceId = persistenceId
		currentState.CreatedTimestampMs = time.Now().UTC().UnixMilli()
	}

	// execute next function
	nextMessages, nextState, nextApplyErr := r.next(c, m, currentState)
	if nextApplyErr != nil {
		logger.ErrorErr("bun state last result mapping error", nextApplyErr)
		return make([]message.Message[[]byte, []byte], 0), nextApplyErr
	}

	// set next updated value
	nextState.UpdatedTimestampMs = time.Now().UTC().UnixMilli()

	// upsert state
	upsertErr := r.repository.Upsert(c, persistenceId, nextState)
	if upsertErr != nil {
		logger.ErrorErr("bun state upsert error", upsertErr)
		return make([]message.Message[[]byte, []byte], 0), upsertErr
	}

	return nextMessages, nil
}
