package stateful_bun_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hjwalt/flows/mock"
	"github.com/uptrace/bun"
)

func ConnectionForTest(t *testing.T) (*gomock.Controller, *mock.MockBunConnection, *mock.MockBunIDB) {
	ctrl := gomock.NewController(t)

	mockIdb := mock.NewMockBunIDB(ctrl)
	mockConnection := mock.NewMockBunConnection(ctrl)
	mockConnection.EXPECT().Db().AnyTimes().DoAndReturn(func() bun.IDB {
		return mockIdb
	})

	return ctrl, mockConnection, mockIdb
}
