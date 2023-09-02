package flows

import (
	"context"

	"github.com/hjwalt/flows/runtime_bun"
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

func RegisterPostgresql(
	name string,
	connectionString string,
	config []runtime.Configuration[*runtime_bun.PostgresqlConnection],
) {
	RegisterPostgresqlConfig(
		runtime_bun.WithApplicationName(name),
		runtime_bun.WithConnectionString(connectionString),
	)
	RegisterPostgresqlConfig(config...)
	inverse.RegisterWithConfigurationRequired[*runtime_bun.PostgresqlConnection](
		QualifierPostgresqlConnection,
		QualifierPostgresqlConnectionConfiguration,
		runtime_bun.NewPostgresqlConnection,
	)
	inverse.Register(QualifierRuntime, InjectorRuntime(QualifierPostgresqlConnection))
}

// Postgresql connection
func RegisterPostgresqlConfig(config ...runtime.Configuration[*runtime_bun.PostgresqlConnection]) {
	inverse.RegisterInstances(QualifierPostgresqlConnectionConfiguration, config)
}

func GetPostgresqlConnection(ctx context.Context) (runtime_bun.BunConnection, error) {
	return inverse.GetLast[runtime_bun.BunConnection](ctx, QualifierPostgresqlConnection)
}
