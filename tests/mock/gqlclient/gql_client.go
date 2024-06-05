// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package mock_gqlclient is a generated GoMock package.
package mock_gqlclient

import (
        context "context"
        reflect "reflect"

        gqlclient "github.com/downsized-devs/sdk-go/gqlclient"
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

// Run mocks base method.
func (m *MockInterface) Run(ctx context.Context, req *gqlclient.Request, resp interface{}) error {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "Run", ctx, req, resp)
        ret0, _ := ret[0].(error)
        return ret0
}

// Run indicates an expected call of Run.
func (mr *MockInterfaceMockRecorder) Run(ctx, req, resp interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockInterface)(nil).Run), ctx, req, resp)
}