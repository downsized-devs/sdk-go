// Code generated by MockGen. DO NOT EDIT.
// Source: ./audit/audit.go

// Package mock_audit is a generated GoMock package.
package mock_audit

import (
	context "context"
	reflect "reflect"

	audit "github.com/downsized-devs/sdk-go/audit"
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

// Capture mocks base method.
func (m *MockInterface) Capture(ctx context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Capture", ctx)
}

// Capture indicates an expected call of Capture.
func (mr *MockInterfaceMockRecorder) Capture(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Capture", reflect.TypeOf((*MockInterface)(nil).Capture), ctx)
}

// Record mocks base method.
func (m *MockInterface) Record(ctx context.Context, log audit.Collection) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Record", ctx, log)
}

// Record indicates an expected call of Record.
func (mr *MockInterfaceMockRecorder) Record(ctx, log interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Record", reflect.TypeOf((*MockInterface)(nil).Record), ctx, log)
}
