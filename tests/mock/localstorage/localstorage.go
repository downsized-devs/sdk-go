// Code generated by MockGen. DO NOT EDIT.
// Source: ./localstorage/localstorage.go
//
// Generated by this command:
//
//	mockgen -source ./localstorage/localstorage.go -destination ./tests/mock/localstorage/localstorage.go
//

// Package mock_localstorage is a generated GoMock package.
package mock_localstorage

import (
	context "context"
	reflect "reflect"

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

// DeleteIndex mocks base method.
func (m *MockInterface) DeleteIndex(ctx context.Context, indexDir string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteIndex", ctx, indexDir)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteIndex indicates an expected call of DeleteIndex.
func (mr *MockInterfaceMockRecorder) DeleteIndex(ctx, indexDir any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteIndex", reflect.TypeOf((*MockInterface)(nil).DeleteIndex), ctx, indexDir)
}

// Index mocks base method.
func (m *MockInterface) Index(ctx context.Context, key string, data any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Index", ctx, key, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// Index indicates an expected call of Index.
func (mr *MockInterfaceMockRecorder) Index(ctx, key, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Index", reflect.TypeOf((*MockInterface)(nil).Index), ctx, key, data)
}

// NewIndex mocks base method.
func (m *MockInterface) NewIndex(ctx context.Context, indexPath string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewIndex", ctx, indexPath)
	ret0, _ := ret[0].(error)
	return ret0
}

// NewIndex indicates an expected call of NewIndex.
func (mr *MockInterfaceMockRecorder) NewIndex(ctx, indexPath any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewIndex", reflect.TypeOf((*MockInterface)(nil).NewIndex), ctx, indexPath)
}

// Search mocks base method.
func (m *MockInterface) Search(ctx context.Context, query string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, query)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockInterfaceMockRecorder) Search(ctx, query any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockInterface)(nil).Search), ctx, query)
}
