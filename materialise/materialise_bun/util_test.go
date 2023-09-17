package materialise_bun_test

import (
	"testing"

	"github.com/hjwalt/flows/materialise/materialise_bun"
	"github.com/stretchr/testify/assert"
)

type TestTable struct {
	BaseModel          string `bun:"table:notifications,alias:n"`
	Uuid               string `bun:",pk"`
	BusinessId         string
	OutletId           string
	UserId             string
	Title              string
	HtmlContent        string
	IsRead             bool
	IsDeleted          bool
	CreatedTimestampMs int64
	ReadTimestampMs    int64
	UpdatedTimestampMs int64
	DeletedTimestampMs int64
}

func TestColumnsStruct(t *testing.T) {
	assert := assert.New(t)

	pks, cols := materialise_bun.Columns(TestTable{})

	assert.Equal([]string{"uuid"}, pks)
	assert.Equal([]string{"business_id", "outlet_id", "user_id", "title", "html_content", "is_read", "is_deleted", "created_timestamp_ms", "read_timestamp_ms", "updated_timestamp_ms", "deleted_timestamp_ms"}, cols)
}

func TestColumnsPointer(t *testing.T) {
	assert := assert.New(t)

	pks, cols := materialise_bun.Columns(&TestTable{})

	assert.Equal([]string{"uuid"}, pks)
	assert.Equal([]string{"business_id", "outlet_id", "user_id", "title", "html_content", "is_read", "is_deleted", "created_timestamp_ms", "read_timestamp_ms", "updated_timestamp_ms", "deleted_timestamp_ms"}, cols)
}

func TestColumnsGeneric(t *testing.T) {
	assert := assert.New(t)

	pks, cols := materialise_bun.ColumnsGeneric(TestTable{})

	assert.Equal([]string{"uuid"}, pks)
	assert.Equal([]string{"business_id", "outlet_id", "user_id", "title", "html_content", "is_read", "is_deleted", "created_timestamp_ms", "read_timestamp_ms", "updated_timestamp_ms", "deleted_timestamp_ms"}, cols)
}

func TestColumnsGenericPointer(t *testing.T) {
	assert := assert.New(t)

	pks, cols := materialise_bun.ColumnsGeneric[*TestTable](nil)

	assert.Equal([]string{"uuid"}, pks)
	assert.Equal([]string{"business_id", "outlet_id", "user_id", "title", "html_content", "is_read", "is_deleted", "created_timestamp_ms", "read_timestamp_ms", "updated_timestamp_ms", "deleted_timestamp_ms"}, cols)
}
