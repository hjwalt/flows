package flows

import (
	"github.com/hjwalt/flows/runtime_bun"
	"github.com/hjwalt/runway/runtime"
)

func RegisterPostgresqlStateful(
	name string,
	connectionString string,
	tableName string,
	config []runtime.Configuration[*runtime_bun.PostgresqlConnection],
) {
	RegisterPostgresqlConfig(
		runtime_bun.WithApplicationName(name),
		runtime_bun.WithConnectionString(connectionString),
	)
	RegisterPostgresqlConfig(config...)
	RegisterPostgresql()
	RegisterPostgresqlSingleState(tableName)
}

func RegisterPostgresqlMaterialise[S any](
	name string,
	connectionString string,
	config []runtime.Configuration[*runtime_bun.PostgresqlConnection],
) {
	RegisterPostgresqlConfig(
		runtime_bun.WithApplicationName(name),
		runtime_bun.WithConnectionString(connectionString),
	)
	RegisterPostgresqlConfig(config...)
	RegisterPostgresql()
	RegisterPostgresqlUpsert[S]()
}
