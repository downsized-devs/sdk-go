// Code generated by MockGen. DO NOT EDIT.
// Source: ./sql/sql.go
//
// Generated by this command:
//
//	mockgen -source ./sql/sql.go -destination ./tests/mock/sql/sql.go
//

// Package mock_sql is a generated GoMock package.
package mock_sql

import (
	reflect "reflect"

	sql "github.com/downsized-devs/sdk-go/sql"
	gomock "go.uber.org/mock/gomock"
)

// MockInterface is a mock of Interface interface.
type MockInterface struct {
	ctrl     *gomock.Controller
	recorder *MockInterfaceMockRecorder
}

// MockInterfaceMockRecorder is the mock recorder for MockInterface.
type MockInterfaceMockRecorder struct {
	mock *MockInterface
}

// NewMockInterface creates a new mock instance.
func NewMockInterface(ctrl *gomock.Controller) *MockInterface {
	mock := &MockInterface{ctrl: ctrl}
	mock.recorder = &MockInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInterface) EXPECT() *MockInterfaceMockRecorder {
	return m.recorder
}

// Follower mocks base method.
func (m *MockInterface) Follower() sql.Command {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Follower")
	ret0, _ := ret[0].(sql.Command)
	return ret0
}

// Follower indicates an expected call of Follower.
func (mr *MockInterfaceMockRecorder) Follower() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Follower", reflect.TypeOf((*MockInterface)(nil).Follower))
}

// Leader mocks base method.
func (m *MockInterface) Leader() sql.Command {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Leader")
	ret0, _ := ret[0].(sql.Command)
	return ret0
}

// Leader indicates an expected call of Leader.
func (mr *MockInterfaceMockRecorder) Leader() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Leader", reflect.TypeOf((*MockInterface)(nil).Leader))
}

// Stop mocks base method.
func (m *MockInterface) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop.
func (mr *MockInterfaceMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockInterface)(nil).Stop))
}
