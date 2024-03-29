// Code generated by MockGen. DO NOT EDIT.
// Source: mock/interfaces.go

// Package mock is a generated GoMock package.
package test_helper

import (
	context "context"
	sql "database/sql"
	reflect "reflect"

	sarama "github.com/IBM/sarama"
	gomock "github.com/golang/mock/gomock"
	bun "github.com/uptrace/bun"
	schema "github.com/uptrace/bun/schema"
)

// MockConsumerGroupSession is a mock of ConsumerGroupSession interface.
type MockConsumerGroupSession struct {
	ctrl     *gomock.Controller
	recorder *MockConsumerGroupSessionMockRecorder
}

// MockConsumerGroupSessionMockRecorder is the mock recorder for MockConsumerGroupSession.
type MockConsumerGroupSessionMockRecorder struct {
	mock *MockConsumerGroupSession
}

// NewMockConsumerGroupSession creates a new mock instance.
func NewMockConsumerGroupSession(ctrl *gomock.Controller) *MockConsumerGroupSession {
	mock := &MockConsumerGroupSession{ctrl: ctrl}
	mock.recorder = &MockConsumerGroupSessionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConsumerGroupSession) EXPECT() *MockConsumerGroupSessionMockRecorder {
	return m.recorder
}

// Claims mocks base method.
func (m *MockConsumerGroupSession) Claims() map[string][]int32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Claims")
	ret0, _ := ret[0].(map[string][]int32)
	return ret0
}

// Claims indicates an expected call of Claims.
func (mr *MockConsumerGroupSessionMockRecorder) Claims() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Claims", reflect.TypeOf((*MockConsumerGroupSession)(nil).Claims))
}

// Commit mocks base method.
func (m *MockConsumerGroupSession) Commit() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Commit")
}

// Commit indicates an expected call of Commit.
func (mr *MockConsumerGroupSessionMockRecorder) Commit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockConsumerGroupSession)(nil).Commit))
}

// Context mocks base method.
func (m *MockConsumerGroupSession) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockConsumerGroupSessionMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockConsumerGroupSession)(nil).Context))
}

// GenerationID mocks base method.
func (m *MockConsumerGroupSession) GenerationID() int32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerationID")
	ret0, _ := ret[0].(int32)
	return ret0
}

// GenerationID indicates an expected call of GenerationID.
func (mr *MockConsumerGroupSessionMockRecorder) GenerationID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerationID", reflect.TypeOf((*MockConsumerGroupSession)(nil).GenerationID))
}

// MarkMessage mocks base method.
func (m *MockConsumerGroupSession) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "MarkMessage", msg, metadata)
}

// MarkMessage indicates an expected call of MarkMessage.
func (mr *MockConsumerGroupSessionMockRecorder) MarkMessage(msg, metadata interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkMessage", reflect.TypeOf((*MockConsumerGroupSession)(nil).MarkMessage), msg, metadata)
}

// MarkOffset mocks base method.
func (m *MockConsumerGroupSession) MarkOffset(topic string, partition int32, offset int64, metadata string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "MarkOffset", topic, partition, offset, metadata)
}

// MarkOffset indicates an expected call of MarkOffset.
func (mr *MockConsumerGroupSessionMockRecorder) MarkOffset(topic, partition, offset, metadata interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkOffset", reflect.TypeOf((*MockConsumerGroupSession)(nil).MarkOffset), topic, partition, offset, metadata)
}

// MemberID mocks base method.
func (m *MockConsumerGroupSession) MemberID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MemberID")
	ret0, _ := ret[0].(string)
	return ret0
}

// MemberID indicates an expected call of MemberID.
func (mr *MockConsumerGroupSessionMockRecorder) MemberID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MemberID", reflect.TypeOf((*MockConsumerGroupSession)(nil).MemberID))
}

// ResetOffset mocks base method.
func (m *MockConsumerGroupSession) ResetOffset(topic string, partition int32, offset int64, metadata string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ResetOffset", topic, partition, offset, metadata)
}

// ResetOffset indicates an expected call of ResetOffset.
func (mr *MockConsumerGroupSessionMockRecorder) ResetOffset(topic, partition, offset, metadata interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetOffset", reflect.TypeOf((*MockConsumerGroupSession)(nil).ResetOffset), topic, partition, offset, metadata)
}

// MockBunIDB is a mock of BunIDB interface.
type MockBunIDB struct {
	ctrl     *gomock.Controller
	recorder *MockBunIDBMockRecorder
}

// MockBunIDBMockRecorder is the mock recorder for MockBunIDB.
type MockBunIDBMockRecorder struct {
	mock *MockBunIDB
}

// NewMockBunIDB creates a new mock instance.
func NewMockBunIDB(ctrl *gomock.Controller) *MockBunIDB {
	mock := &MockBunIDB{ctrl: ctrl}
	mock.recorder = &MockBunIDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBunIDB) EXPECT() *MockBunIDBMockRecorder {
	return m.recorder
}

// BeginTx mocks base method.
func (m *MockBunIDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (bun.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginTx", ctx, opts)
	ret0, _ := ret[0].(bun.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BeginTx indicates an expected call of BeginTx.
func (mr *MockBunIDBMockRecorder) BeginTx(ctx, opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginTx", reflect.TypeOf((*MockBunIDB)(nil).BeginTx), ctx, opts)
}

// Dialect mocks base method.
func (m *MockBunIDB) Dialect() schema.Dialect {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Dialect")
	ret0, _ := ret[0].(schema.Dialect)
	return ret0
}

// Dialect indicates an expected call of Dialect.
func (mr *MockBunIDBMockRecorder) Dialect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dialect", reflect.TypeOf((*MockBunIDB)(nil).Dialect))
}

// ExecContext mocks base method.
func (m *MockBunIDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ExecContext", varargs...)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecContext indicates an expected call of ExecContext.
func (mr *MockBunIDBMockRecorder) ExecContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecContext", reflect.TypeOf((*MockBunIDB)(nil).ExecContext), varargs...)
}

// NewAddColumn mocks base method.
func (m *MockBunIDB) NewAddColumn() *bun.AddColumnQuery {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewAddColumn")
	ret0, _ := ret[0].(*bun.AddColumnQuery)
	return ret0
}

// NewAddColumn indicates an expected call of NewAddColumn.
func (mr *MockBunIDBMockRecorder) NewAddColumn() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewAddColumn", reflect.TypeOf((*MockBunIDB)(nil).NewAddColumn))
}

// NewCreateIndex mocks base method.
func (m *MockBunIDB) NewCreateIndex() *bun.CreateIndexQuery {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewCreateIndex")
	ret0, _ := ret[0].(*bun.CreateIndexQuery)
	return ret0
}

// NewCreateIndex indicates an expected call of NewCreateIndex.
func (mr *MockBunIDBMockRecorder) NewCreateIndex() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewCreateIndex", reflect.TypeOf((*MockBunIDB)(nil).NewCreateIndex))
}

// NewCreateTable mocks base method.
func (m *MockBunIDB) NewCreateTable() *bun.CreateTableQuery {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewCreateTable")
	ret0, _ := ret[0].(*bun.CreateTableQuery)
	return ret0
}

// NewCreateTable indicates an expected call of NewCreateTable.
func (mr *MockBunIDBMockRecorder) NewCreateTable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewCreateTable", reflect.TypeOf((*MockBunIDB)(nil).NewCreateTable))
}

// NewDelete mocks base method.
func (m *MockBunIDB) NewDelete() *bun.DeleteQuery {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewDelete")
	ret0, _ := ret[0].(*bun.DeleteQuery)
	return ret0
}

// NewDelete indicates an expected call of NewDelete.
func (mr *MockBunIDBMockRecorder) NewDelete() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewDelete", reflect.TypeOf((*MockBunIDB)(nil).NewDelete))
}

// NewDropColumn mocks base method.
func (m *MockBunIDB) NewDropColumn() *bun.DropColumnQuery {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewDropColumn")
	ret0, _ := ret[0].(*bun.DropColumnQuery)
	return ret0
}

// NewDropColumn indicates an expected call of NewDropColumn.
func (mr *MockBunIDBMockRecorder) NewDropColumn() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewDropColumn", reflect.TypeOf((*MockBunIDB)(nil).NewDropColumn))
}

// NewDropIndex mocks base method.
func (m *MockBunIDB) NewDropIndex() *bun.DropIndexQuery {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewDropIndex")
	ret0, _ := ret[0].(*bun.DropIndexQuery)
	return ret0
}

// NewDropIndex indicates an expected call of NewDropIndex.
func (mr *MockBunIDBMockRecorder) NewDropIndex() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewDropIndex", reflect.TypeOf((*MockBunIDB)(nil).NewDropIndex))
}

// NewDropTable mocks base method.
func (m *MockBunIDB) NewDropTable() *bun.DropTableQuery {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewDropTable")
	ret0, _ := ret[0].(*bun.DropTableQuery)
	return ret0
}

// NewDropTable indicates an expected call of NewDropTable.
func (mr *MockBunIDBMockRecorder) NewDropTable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewDropTable", reflect.TypeOf((*MockBunIDB)(nil).NewDropTable))
}

// NewInsert mocks base method.
func (m *MockBunIDB) NewInsert() *bun.InsertQuery {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewInsert")
	ret0, _ := ret[0].(*bun.InsertQuery)
	return ret0
}

// NewInsert indicates an expected call of NewInsert.
func (mr *MockBunIDBMockRecorder) NewInsert() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewInsert", reflect.TypeOf((*MockBunIDB)(nil).NewInsert))
}

// NewRaw mocks base method.
func (m *MockBunIDB) NewRaw(query string, args ...interface{}) *bun.RawQuery {
	m.ctrl.T.Helper()
	varargs := []interface{}{query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "NewRaw", varargs...)
	ret0, _ := ret[0].(*bun.RawQuery)
	return ret0
}

// NewRaw indicates an expected call of NewRaw.
func (mr *MockBunIDBMockRecorder) NewRaw(query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRaw", reflect.TypeOf((*MockBunIDB)(nil).NewRaw), varargs...)
}

// NewSelect mocks base method.
func (m *MockBunIDB) NewSelect() *bun.SelectQuery {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewSelect")
	ret0, _ := ret[0].(*bun.SelectQuery)
	return ret0
}

// NewSelect indicates an expected call of NewSelect.
func (mr *MockBunIDBMockRecorder) NewSelect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewSelect", reflect.TypeOf((*MockBunIDB)(nil).NewSelect))
}

// NewTruncateTable mocks base method.
func (m *MockBunIDB) NewTruncateTable() *bun.TruncateTableQuery {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewTruncateTable")
	ret0, _ := ret[0].(*bun.TruncateTableQuery)
	return ret0
}

// NewTruncateTable indicates an expected call of NewTruncateTable.
func (mr *MockBunIDBMockRecorder) NewTruncateTable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewTruncateTable", reflect.TypeOf((*MockBunIDB)(nil).NewTruncateTable))
}

// NewUpdate mocks base method.
func (m *MockBunIDB) NewUpdate() *bun.UpdateQuery {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewUpdate")
	ret0, _ := ret[0].(*bun.UpdateQuery)
	return ret0
}

// NewUpdate indicates an expected call of NewUpdate.
func (mr *MockBunIDBMockRecorder) NewUpdate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewUpdate", reflect.TypeOf((*MockBunIDB)(nil).NewUpdate))
}

// NewValues mocks base method.
func (m *MockBunIDB) NewValues(model interface{}) *bun.ValuesQuery {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewValues", model)
	ret0, _ := ret[0].(*bun.ValuesQuery)
	return ret0
}

// NewValues indicates an expected call of NewValues.
func (mr *MockBunIDBMockRecorder) NewValues(model interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewValues", reflect.TypeOf((*MockBunIDB)(nil).NewValues), model)
}

// QueryContext mocks base method.
func (m *MockBunIDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryContext", varargs...)
	ret0, _ := ret[0].(*sql.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryContext indicates an expected call of QueryContext.
func (mr *MockBunIDBMockRecorder) QueryContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryContext", reflect.TypeOf((*MockBunIDB)(nil).QueryContext), varargs...)
}

// QueryRowContext mocks base method.
func (m *MockBunIDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryRowContext", varargs...)
	ret0, _ := ret[0].(*sql.Row)
	return ret0
}

// QueryRowContext indicates an expected call of QueryRowContext.
func (mr *MockBunIDBMockRecorder) QueryRowContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRowContext", reflect.TypeOf((*MockBunIDB)(nil).QueryRowContext), varargs...)
}

// RunInTx mocks base method.
func (m *MockBunIDB) RunInTx(ctx context.Context, opts *sql.TxOptions, f func(context.Context, bun.Tx) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunInTx", ctx, opts, f)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunInTx indicates an expected call of RunInTx.
func (mr *MockBunIDBMockRecorder) RunInTx(ctx, opts, f interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunInTx", reflect.TypeOf((*MockBunIDB)(nil).RunInTx), ctx, opts, f)
}
