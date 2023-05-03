package materialise_bun

import (
	"context"
	"errors"
	"strings"

	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/flows/runtime_bun"
	"github.com/hjwalt/runway/reflect"
)

// constructor
func NewBunUpsertRepository[T any](configurations ...runtime.Configuration[*BunUpsertRepository[T]]) *BunUpsertRepository[T] {
	r := &BunUpsertRepository[T]{}

	for _, configuration := range configurations {
		r = configuration(r)
	}
	return r
}

// configuration
func WithBunUpsertRepositoryConnection[T any](connection runtime_bun.BunConnection) runtime.Configuration[*BunUpsertRepository[T]] {
	return func(r *BunUpsertRepository[T]) *BunUpsertRepository[T] {
		r.connection = connection
		return r
	}
}

// implementation
type BunUpsertRepository[T any] struct {
	connection runtime_bun.BunConnection
}

func (r BunUpsertRepository[T]) Upsert(c context.Context, m []T) error {
	if len(m) == 0 {
		return nil
	}

	instanceExample := reflect.Construct[T]()
	pks, cols := Columns(instanceExample)

	if len(pks) == 0 {
		return errors.New("bun upsert: missing primary key, see https://bun.uptrace.dev/guide/models.html#mapping-tables-to-structs")
	}
	if len(cols) == 0 {
		return errors.New("bun upsert: no columns other than primary key")
	}

	pkJoined := strings.Join(pks, ",")
	conflictSet := ConflictSet(cols)

	query := r.connection.Db().
		NewInsert().
		Model(&m).
		On("CONFLICT (" + pkJoined + ") DO UPDATE").
		Set(conflictSet)

	_, err := query.Exec(c)

	if err != nil {
		return errors.Join(errors.New("bun upsert"), err)
	}

	return nil
}
