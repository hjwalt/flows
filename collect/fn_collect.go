package collect

import (
	"context"
	"errors"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/runtime"
	"github.com/hjwalt/runway/structure"
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
func WithCollectPersistenceIdFunc(persistenceIdFunc func(context.Context, flow.Message[structure.Bytes, structure.Bytes]) (string, error)) runtime.Configuration[*Collect] {
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

func WithCollectAggregator(next Aggregator[structure.Bytes, structure.Bytes, structure.Bytes]) runtime.Configuration[*Collect] {
	return func(st *Collect) *Collect {
		st.next = next
		return st
	}
}

type Collect struct {
	next              Aggregator[structure.Bytes, structure.Bytes, structure.Bytes]
	persistenceIdFunc stateful.PersistenceIdFunction[structure.Bytes, structure.Bytes]
	collect           Collector
}

func (r *Collect) Apply(c context.Context, ms []flow.Message[structure.Bytes, structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], error) {
	stateMap := make(map[string]stateful.State[[]byte])
	for _, m := range ms {
		persistenceId, persistenceIdErr := r.persistenceIdFunc(c, m)
		if persistenceIdErr != nil {
			return flow.EmptySlice(), errors.Join(persistenceIdErr, ErrBatchCollect, stateful.ErrPersistenceId)
		}

		if state, stateExist := stateMap[persistenceId]; stateExist {
			nextState, nextApplyErr := r.next(c, m, state)
			if nextApplyErr != nil {
				return flow.EmptySlice(), errors.Join(nextApplyErr, ErrBatchCollect, ErrCollectStateCreate)
			}
			stateMap[persistenceId] = nextState
		} else {
			nextState, nextApplyErr := r.next(c, m, stateful.NewState[[]byte](persistenceId, []byte{}))
			if nextApplyErr != nil {
				return flow.EmptySlice(), errors.Join(nextApplyErr, ErrBatchCollect, ErrCollectStateCreate)
			}
			stateMap[persistenceId] = nextState
		}
	}

	resultMessages := flow.EmptySlice()
	for k, s := range stateMap {
		nextMessages, nextMessageErr := r.collect(c, k, s)
		if nextMessageErr != nil {
			return flow.EmptySlice(), errors.Join(nextMessageErr, ErrBatchCollect, ErrCollectStateMap)
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
