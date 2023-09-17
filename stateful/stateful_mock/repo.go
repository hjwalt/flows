package stateful_mock

import (
	"context"

	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/runway/structure"
)

func Repository() stateful.Repository {
	return &repo{State: map[string]stateful.State[[]byte]{}}
}

type repo struct {
	State map[string]stateful.State[structure.Bytes]
}

func (r *repo) Get(ctx context.Context, persistenceId string) (stateful.State[structure.Bytes], error) {
	if state, statePresent := r.State[persistenceId]; statePresent {
		return state, nil
	} else {
		return stateful.NewState[[]byte](persistenceId, []byte{}), nil
	}
}

func (r *repo) GetAll(ctx context.Context, persistenceIds []string) (map[string]stateful.State[structure.Bytes], error) {
	stateMap := map[string]stateful.State[structure.Bytes]{}
	for _, persistenceId := range persistenceIds {
		if state, statePresent := r.State[persistenceId]; statePresent {
			stateMap[persistenceId] = state
		} else {
			stateMap[persistenceId] = stateful.NewState[[]byte](persistenceId, []byte{})
		}
	}
	return stateMap, nil
}

func (r *repo) Upsert(ctx context.Context, persistenceId string, dbState stateful.State[structure.Bytes]) error {
	r.State[persistenceId] = dbState
	return nil
}

func (r *repo) UpsertAll(ctx context.Context, stateMap map[string]stateful.State[structure.Bytes]) error {
	for k, v := range stateMap {
		r.State[k] = v
	}
	return nil
}
