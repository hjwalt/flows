package flows

import (
	"github.com/hjwalt/flows/materialise"
	"github.com/hjwalt/flows/materialise_bun"
	"github.com/hjwalt/flows/runtime_bun"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateful_bun"
	"github.com/hjwalt/runway/runtime"
)

func Postgresql(configs []runtime.Configuration[*runtime_bun.PostgresqlConnection]) runtime_bun.BunConnection {
	return runtime_bun.NewPostgresqlConnection(configs...)
}

func PostgresqlSingleStateRepository(conn runtime_bun.BunConnection, tableName string) stateful.SingleStateRepository {
	return stateful_bun.NewSingleStateRepository(
		stateful_bun.WithSingleStateRepositoryConnection(conn),
		stateful_bun.WithSingleStateRepositoryPersistenceTableName(tableName),
	)
}

func PostgresqlUpsertRepository[T any](conn runtime_bun.BunConnection) materialise.UpsertRepository[T] {
	return materialise_bun.NewBunUpsertRepository(
		materialise_bun.WithBunUpsertRepositoryConnection[T](conn),
	)
}
