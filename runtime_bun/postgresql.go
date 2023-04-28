package runtime_bun

import (
	"database/sql"
	"time"

	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/runway/logger"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

// constructor
func NewPostgresqlConnection(configurations ...runtime.Configuration[*PostgresqlConnection]) BunConnection {
	connection := &PostgresqlConnection{}
	for _, configuration := range configurations {
		connection = configuration(connection)
	}
	return connection
}

// configuration
func WithController(controller runtime.Controller) runtime.Configuration[*PostgresqlConnection] {
	return func(c *PostgresqlConnection) *PostgresqlConnection {
		c.Controller = controller
		return c
	}
}

func WithConnectionString(connectionString string) runtime.Configuration[*PostgresqlConnection] {
	return func(c *PostgresqlConnection) *PostgresqlConnection {
		c.ConnectionString = connectionString
		return c
	}
}

func WithApplicationName(applicationName string) runtime.Configuration[*PostgresqlConnection] {
	return func(c *PostgresqlConnection) *PostgresqlConnection {
		c.ApplicationName = applicationName
		return c
	}
}

func WithMaxOpenConnections(maxOpenConns int) runtime.Configuration[*PostgresqlConnection] {
	return func(c *PostgresqlConnection) *PostgresqlConnection {
		c.MaxOpenConns = maxOpenConns
		return c
	}
}

// implementation
type PostgresqlConnection struct {
	ConnectionString string
	ApplicationName  string
	MaxOpenConns     int
	Controller       runtime.Controller

	db *bun.DB
}

func (r *PostgresqlConnection) Start() error {
	logger.Info("starting bun")

	if len(r.ApplicationName) == 0 {
		r.ApplicationName = "flows"
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(r.ConnectionString),
		pgdriver.WithReadTimeout(5*time.Second),
		pgdriver.WithWriteTimeout(5*time.Second),
		pgdriver.WithApplicationName(r.ApplicationName),
		pgdriver.WithConnParams(map[string]interface{}{
			"bytea_output":                  "hex",             // ensures that misconfigured database does not brick state management
			"default_transaction_isolation": "repeatable read", // ensures that multiple threads accessing the same persistence id does not encounter race condition
		}),
	))

	if r.MaxOpenConns > 10 {
		sqldb.SetMaxOpenConns(r.MaxOpenConns)
	} else {
		sqldb.SetMaxOpenConns(10)
	}

	bunDb := bun.NewDB(sqldb, pgdialect.New())
	bunDb.AddQueryHook(bundebug.NewQueryHook())

	r.db = bunDb

	logger.Info("started bun")

	// mark started in controller
	r.Controller.Started()

	return nil
}

func (r *PostgresqlConnection) Stop() {
	logger.Info("stopping bun")
	r.db.Close()
	r.Controller.Stopped()
	logger.Info("stopped bun")
}

func (r *PostgresqlConnection) Db() bun.IDB {
	return r.db
}
