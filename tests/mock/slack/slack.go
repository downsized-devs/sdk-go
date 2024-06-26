// Code generated by MockGen. DO NOT EDIT.
// Source: ./slack/slack.go
//
// Generated by this command:
//
//	mockgen -source ./slack/slack.go -destination ./tests/mock/slack/slack.go
//

// Package mock_slack is a generated GoMock package.
package mock_slack

import (
	context "context"
	reflect "reflect"

	slack "github.com/downsized-devs/sdk-go/slack"
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

// SendMessage mocks base method.
func (m *MockInterface) SendMessage(ctx context.Context, channelID string, attachment slack.Attachment, fields []slack.AttachmentField) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", ctx, channelID, attachment, fields)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockInterfaceMockRecorder) SendMessage(ctx, channelID, attachment, fields any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockInterface)(nil).SendMessage), ctx, channelID, attachment, fields)
}
