package flows

import (
	"context"

	"github.com/hjwalt/flows/runtime_bun"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
)

const (
	QualifierPostgresqlConnection = "QualifierPostgresqlConnection"
)

func RegisterPostgresql(
	container inverse.Container,
	name string,
	connectionString string,
	configs []runtime.Configuration[*runtime_bun.PostgresqlConnection],
) {

	resolver := runtime.NewResolver[*runtime_bun.PostgresqlConnection, runtime_bun.BunConnection](
		QualifierPostgresqlConnection,
		container,
		true,
		runtime_bun.NewPostgresqlConnection,
	)

	resolver.AddConfigVal(runtime_bun.WithApplicationName(name))
	resolver.AddConfigVal(runtime_bun.WithConnectionString(connectionString))

	for _, config := range configs {
		resolver.AddConfigVal(config)
	}

	resolver.Register()

	RegisterRuntime(QualifierPostgresqlConnection, container)
}

func GetPostgresqlConnection(ctx context.Context, ci inverse.Container) (runtime_bun.BunConnection, error) {
	return inverse.GenericGetLast[runtime_bun.BunConnection](ci, ctx, QualifierPostgresqlConnection)
}

// ===================================

func RegisterPostgresqlConfig(ci inverse.Container, configs ...runtime.Configuration[*runtime_bun.PostgresqlConnection]) {
	for _, config := range configs {
		ci.AddVal(runtime.QualifierConfig(QualifierPostgresqlConnection), config)
	}
}
