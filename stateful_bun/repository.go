package stateful_bun

import (
	"context"
	"database/sql"

	"github.com/hjwalt/flows/runtime_bun"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/runway/runtime"
	"github.com/hjwalt/runway/structure"
	"github.com/uptrace/bun"
)

// constructor
func NewRepository(configurations ...runtime.Configuration[Repository]) stateful.Repository {
	repository := Repository{}
	for _, configuration := range configurations {
		repository = configuration(repository)
	}
	return repository
}

// configuration
func WithConnection(connection runtime_bun.BunConnection) runtime.Configuration[Repository] {
	return func(st Repository) Repository {
		st.connection = connection
		return st
	}
}

func WithStateTableName(persistenceTableName string) runtime.Configuration[Repository] {
	return func(st Repository) Repository {
		st.stateTableName = persistenceTableName
		return st
	}
}

// implementation
type Repository struct {
	stateTableName string
	connection     runtime_bun.BunConnection
}

func (r Repository) Get(ctx context.Context, persistenceId string) (stateful.State[structure.Bytes], error) {
	dbState := &StateTable{}

	readErr := r.connection.Db().
		NewSelect().
		Model(dbState).
		ModelTableExpr(r.stateTableName+" AS state_table").
		Where("id = ?", persistenceId).
		Scan(ctx)

	if readErr != nil && readErr != sql.ErrNoRows {
		return stateful.NewState[structure.Bytes](persistenceId, []byte{}), readErr
	}

	state, convertErr := TableToState(dbState)
	if convertErr != nil {
		return stateful.NewState[structure.Bytes](persistenceId, []byte{}), convertErr
	}

	state.Id = persistenceId

	return state, nil
}

func (r Repository) GetAll(ctx context.Context, persistenceId []string) (map[string]stateful.State[structure.Bytes], error) {
	dbState := []StateTable{}

	readErr := r.connection.Db().
		NewSelect().
		Model(&dbState).
		ModelTableExpr(r.stateTableName+" AS state_table").
		Where("id IN (?)", bun.In(persistenceId)).
		Scan(ctx)

	if readErr != nil && readErr != sql.ErrNoRows {
		return map[string]stateful.State[structure.Bytes]{}, readErr
	}

	stateMap := map[string]stateful.State[structure.Bytes]{}
	for _, stateEntry := range dbState {
		if stateMapped, mapErr := TableToState(&stateEntry); mapErr != nil {
			return map[string]stateful.State[structure.Bytes]{}, readErr
		} else {
			stateMap[stateEntry.Id] = stateMapped
		}
	}

	for _, persistenceIdEntry := range persistenceId {
		if _, idPresent := stateMap[persistenceIdEntry]; !idPresent {
			stateMap[persistenceIdEntry] = stateful.NewState[structure.Bytes](persistenceIdEntry, []byte{})
		}
	}

	return stateMap, nil
}

func (r Repository) Upsert(ctx context.Context, persistenceId string, s stateful.State[structure.Bytes]) error {

	dbState, err := StateToTable(s)
	if err != nil {
		return err
	}

	dbState.Id = persistenceId

	_, upsertErr := r.connection.Db().
		NewInsert().
		Model(dbState).
		ModelTableExpr(r.stateTableName + " AS state_table").
		On("CONFLICT (id) DO UPDATE").
		Set("internal = EXCLUDED.internal").
		Set("results = EXCLUDED.results").
		Set("content = EXCLUDED.content").
		Set("created_timestamp_ms = EXCLUDED.created_timestamp_ms").
		Set("updated_timestamp_ms = EXCLUDED.updated_timestamp_ms").
		Exec(ctx)

	return upsertErr
}

func (r Repository) UpsertAll(ctx context.Context, stateMap map[string]stateful.State[structure.Bytes]) error {

	dbStates := []*StateTable{}

	for k, v := range stateMap {
		dbState, err := StateToTable(v)
		if err != nil {
			return err
		}
		dbState.Id = k

		dbStates = append(dbStates, dbState)
	}

	_, upsertErr := r.connection.Db().
		NewInsert().
		Model(&dbStates).
		ModelTableExpr(r.stateTableName + " AS state_table").
		On("CONFLICT (id) DO UPDATE").
		Set("internal = EXCLUDED.internal").
		Set("results = EXCLUDED.results").
		Set("content = EXCLUDED.content").
		Set("created_timestamp_ms = EXCLUDED.created_timestamp_ms").
		Set("updated_timestamp_ms = EXCLUDED.updated_timestamp_ms").
		Exec(ctx)

	return upsertErr
}
