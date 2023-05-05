package test_helper

import (
	"context"
	sql "database/sql"
	"os"

	"github.com/hjwalt/runway/reflect"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

type SqliteBunConnection[T any] struct {
	db            *bun.DB
	FixtureFolder string
	FixtureFile   string
}

func (r *SqliteBunConnection[T]) Start() error {
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		return err
	}
	sqldb.SetMaxOpenConns(1)

	db := bun.NewDB(sqldb, sqlitedialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	// Register models before loading fixtures.
	db.RegisterModel(reflect.Construct[T]())

	// Automatically create tables.
	fixture := dbfixture.New(db, dbfixture.WithRecreateTables())

	// Load fixtures.
	if err := fixture.Load(context.Background(), os.DirFS(r.FixtureFolder), r.FixtureFile); err != nil {
		return err
	}

	r.db = db
	return nil
}

func (r *SqliteBunConnection[T]) Stop() {
	r.db.Close()
}

func (r *SqliteBunConnection[T]) Db() bun.IDB {
	return r.db
}
