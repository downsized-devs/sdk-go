// Code generated by MockGen. DO NOT EDIT.
// Source: ./messaging/messaging.go
//
// Generated by this command:
//
//	mockgen -source ./messaging/messaging.go -destination ./tests/mock/messaging/messaging.go
//

// Package mock_messaging is a generated GoMock package.
package mock_messaging

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

// BatchSendDryRun mocks base method.
func (m *MockInterface) BatchSendDryRun(ctx context.Context, tokens []string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchSendDryRun", ctx, tokens)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BatchSendDryRun indicates an expected call of BatchSendDryRun.
func (mr *MockInterfaceMockRecorder) BatchSendDryRun(ctx, tokens any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchSendDryRun", reflect.TypeOf((*MockInterface)(nil).BatchSendDryRun), ctx, tokens)
}

// BroadCastToTopic mocks base method.
func (m *MockInterface) BroadCastToTopic(ctx context.Context, topic string, payload map[string]string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BroadCastToTopic", ctx, topic, payload)
	ret0, _ := ret[0].(error)
	return ret0
}

// BroadCastToTopic indicates an expected call of BroadCastToTopic.
func (mr *MockInterfaceMockRecorder) BroadCastToTopic(ctx, topic, payload any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BroadCastToTopic", reflect.TypeOf((*MockInterface)(nil).BroadCastToTopic), ctx, topic, payload)
}

// SubstribeToTpic mocks base method.
func (m *MockInterface) SubstribeToTpic(ctx context.Context, deviceToken, topic string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubstribeToTpic", ctx, deviceToken, topic)
	ret0, _ := ret[0].(error)
	return ret0
}

// SubstribeToTpic indicates an expected call of SubstribeToTpic.
func (mr *MockInterfaceMockRecorder) SubstribeToTpic(ctx, deviceToken, topic any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubstribeToTpic", reflect.TypeOf((*MockInterface)(nil).SubstribeToTpic), ctx, deviceToken, topic)
}

// UnsubscribeFromTopic mocks base method.
func (m *MockInterface) UnsubscribeFromTopic(ctx context.Context, deviceToken, topic string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnsubscribeFromTopic", ctx, deviceToken, topic)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnsubscribeFromTopic indicates an expected call of UnsubscribeFromTopic.
func (mr *MockInterfaceMockRecorder) UnsubscribeFromTopic(ctx, deviceToken, topic any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnsubscribeFromTopic", reflect.TypeOf((*MockInterface)(nil).UnsubscribeFromTopic), ctx, deviceToken, topic)
}
