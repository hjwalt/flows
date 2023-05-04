package stateful_bun

import (
	"context"
	"database/sql"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/flows/runtime_bun"
	"github.com/hjwalt/flows/stateful"
)

// constructor
func NewSingleStateRepository(configurations ...runtime.Configuration[SingleStateRepository]) stateful.SingleStateRepository {
	repository := SingleStateRepository{}
	for _, configuration := range configurations {
		repository = configuration(repository)
	}
	return repository
}

// configuration
func WithSingleStateRepositoryConnection(connection runtime_bun.BunConnection) runtime.Configuration[SingleStateRepository] {
	return func(st SingleStateRepository) SingleStateRepository {
		st.connection = connection
		return st
	}
}

func WithSingleStateRepositoryPersistenceTableName(persistenceTableName string) runtime.Configuration[SingleStateRepository] {
	return func(st SingleStateRepository) SingleStateRepository {
		st.persistenceTableName = persistenceTableName
		return st
	}
}

// implementation
type SingleStateRepository struct {
	persistenceTableName string
	connection           runtime_bun.BunConnection
}

func (r SingleStateRepository) Get(ctx context.Context, persistenceId string) (stateful.SingleState[message.Bytes], error) {
	dbState := &SingleStateTable{}

	readErr := r.connection.Db().
		NewSelect().
		Model(dbState).
		ModelTableExpr(r.persistenceTableName+" AS single_state_table").
		Where("persistence_id = ?", persistenceId).
		// For("UPDATE"). -- model no longer requires transaction
		Scan(ctx)

	if readErr != nil && readErr != sql.ErrNoRows {
		return stateful.SingleState[message.Bytes]{}, readErr
	}

	return TableToState(dbState)
}

func (r SingleStateRepository) Upsert(ctx context.Context, persistenceId string, s stateful.SingleState[message.Bytes]) error {

	dbState, err := StateToTable(s)
	if err != nil {
		return err
	}

	_, upsertErr := r.connection.Db().
		NewInsert().
		Model(dbState).
		ModelTableExpr(r.persistenceTableName + " AS sql_state_table").
		On("CONFLICT (persistence_id) DO UPDATE").
		Set("internal = EXCLUDED.internal").
		Set("results = EXCLUDED.results").
		Set("content = EXCLUDED.content").
		Set("created_timestamp_ms = EXCLUDED.created_timestamp_ms").
		Set("updated_timestamp_ms = EXCLUDED.updated_timestamp_ms").
		Exec(ctx)

	return upsertErr
}
