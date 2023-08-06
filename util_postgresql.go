package flows

import (
	"context"
	"errors"

	"github.com/hjwalt/flows/materialise"
	"github.com/hjwalt/flows/materialise_bun"
	"github.com/hjwalt/flows/runtime_bun"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateful_bun"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
)

const (
	QualifierPostgresqlConnectionConfiguration        = "QualifierPostgresqlConnectionConfiguration"
	QualifierPostgresqlConnection                     = "QualifierPostgresqlConnection"
	QualifierPostgresqlSingleStateRepository          = "QualifierPostgresqlSingleStateRepository"
	QualifierPostgresqlSingleStateRepositoryTableName = "QualifierPostgresqlSingleStateRepositoryTableName"
	QualifierPostgresqlUpsertRepository               = "QualifierPostgresqlUpsertRepository"
)

// Postgresql connection
func RegisterPostgresql(config []runtime.Configuration[*runtime_bun.PostgresqlConnection]) {
	inverse.RegisterInstances(QualifierPostgresqlConnectionConfiguration, config)
	inverse.Register(QualifierPostgresqlConnection, InjectorPostgresqlConnection)
	inverse.Register(QualifierRuntime, InjectorRuntime(QualifierPostgresqlConnection))
}

func InjectorPostgresqlConnection(ctx context.Context) (runtime_bun.BunConnection, error) {
	configurations, getConfigurationError := inverse.GetAll[runtime.Configuration[*runtime_bun.PostgresqlConnection]](ctx, QualifierPostgresqlConnectionConfiguration)
	if getConfigurationError != nil && !errors.Is(getConfigurationError, inverse.ErrNotInjected) {
		return nil, getConfigurationError
	}
	return runtime_bun.NewPostgresqlConnection(configurations...), nil
}

// Single state repository
func RegisterPostgresqlSingleState(tableName string) {
	inverse.Register(QualifierPostgresqlSingleStateRepository, InjectorPostgresqlSingleStateRepository)
	inverse.RegisterInstance[string](QualifierPostgresqlSingleStateRepositoryTableName, tableName)
}

func InjectorPostgresqlSingleStateRepository(ctx context.Context) (stateful.SingleStateRepository, error) {
	bunConnection, getBunConnectionError := inverse.GetLast[runtime_bun.BunConnection](ctx, QualifierPostgresqlConnection)
	if getBunConnectionError != nil {
		return nil, getBunConnectionError
	}

	tableName, getTableNameError := inverse.GetLast[string](ctx, QualifierPostgresqlSingleStateRepositoryTableName)
	if getTableNameError != nil {
		return nil, getTableNameError
	}

	return stateful_bun.NewSingleStateRepository(
		stateful_bun.WithSingleStateRepositoryConnection(bunConnection),
		stateful_bun.WithSingleStateRepositoryPersistenceTableName(tableName),
	), nil
}

func GetPostgresqlSingleStateRepository(ctx context.Context) (stateful.SingleStateRepository, error) {
	return inverse.GetLast[stateful.SingleStateRepository](ctx, QualifierPostgresqlSingleStateRepository)
}

// Upsert materialiser
func RegisterPostgresqlUpsert[T any]() {
	inverse.Register(QualifierPostgresqlUpsertRepository, InjectorPostgresqlUpsertRepository[T])
}

func InjectorPostgresqlUpsertRepository[T any](ctx context.Context) (materialise.UpsertRepository[T], error) {
	bunConnection, getBunConnectionError := inverse.GetLast[runtime_bun.BunConnection](ctx, QualifierPostgresqlConnection)
	if getBunConnectionError != nil {
		return nil, getBunConnectionError
	}

	return materialise_bun.NewBunUpsertRepository(
		materialise_bun.WithBunUpsertRepositoryConnection[T](bunConnection),
	), nil
}

func GetPostgresqlUpsertRepository[T any](ctx context.Context) (materialise.UpsertRepository[T], error) {
	return inverse.GetLast[materialise.UpsertRepository[T]](ctx, QualifierPostgresqlUpsertRepository)
}
