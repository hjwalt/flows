package collect

import (
	"context"
	"errors"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/runtime"
)

// constructor
func NewCollect(configurations ...runtime.Configuration[*Collect]) stateless.BatchFunction {
	singleFunction := &Collect{}
	for _, configuration := range configurations {
		singleFunction = configuration(singleFunction)
	}
	return singleFunction.Apply
}

// configuration
func WithCollectPersistenceIdFunc(persistenceIdFunc func(context.Context, message.Message[message.Bytes, message.Bytes]) (string, error)) runtime.Configuration[*Collect] {
	return func(st *Collect) *Collect {
		st.persistenceIdFunc = persistenceIdFunc
		return st
	}
}

func WithCollectCollector(collect Collector) runtime.Configuration[*Collect] {
	return func(st *Collect) *Collect {
		st.collect = collect
		return st
	}
}

func WithCollectAggregator(next Aggregator[message.Bytes, message.Bytes, message.Bytes]) runtime.Configuration[*Collect] {
	return func(st *Collect) *Collect {
		st.next = next
		return st
	}
}

type Collect struct {
	next              Aggregator[message.Bytes, message.Bytes, message.Bytes]
	persistenceIdFunc stateful.PersistenceIdFunction[message.Bytes, message.Bytes]
	collect           Collector
}

func (r *Collect) Apply(c context.Context, ms []message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
	stateMap := make(map[string]stateful.State[[]byte])
	for _, m := range ms {
		persistenceId, persistenceIdErr := r.persistenceIdFunc(c, m)
		if persistenceIdErr != nil {
			return message.EmptySlice(), errors.Join(persistenceIdErr, ErrBatchCollect, stateful.ErrPersistenceId)
		}

		if state, stateExist := stateMap[persistenceId]; stateExist {
			nextState, nextApplyErr := r.next(c, m, state)
			if nextApplyErr != nil {
				return message.EmptySlice(), errors.Join(nextApplyErr, ErrBatchCollect, ErrCollectStateCreate)
			}
			stateMap[persistenceId] = nextState
		} else {
			nextState, nextApplyErr := r.next(c, m, stateful.NewState[[]byte](persistenceId, []byte{}))
			if nextApplyErr != nil {
				return message.EmptySlice(), errors.Join(nextApplyErr, ErrBatchCollect, ErrCollectStateCreate)
			}
			stateMap[persistenceId] = nextState
		}
	}

	resultMessages := message.EmptySlice()
	for k, s := range stateMap {
		nextMessages, nextMessageErr := r.collect(c, k, s)
		if nextMessageErr != nil {
			return message.EmptySlice(), errors.Join(nextMessageErr, ErrBatchCollect, ErrCollectStateMap)
		}
		resultMessages = append(resultMessages, nextMessages...)
	}

	return resultMessages, nil
}

var (
	ErrBatchCollect       = errors.New("batch collect")
	ErrCollectStateMap    = errors.New("collect map message from state")
	ErrCollectStateCreate = errors.New("collect map state from input message")
)
