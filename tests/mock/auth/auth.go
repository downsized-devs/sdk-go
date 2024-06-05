// Code generated by MockGen. DO NOT EDIT.
// Source: ./auth/auth.go
//
// Generated by this command:
//
//	mockgen -source ./auth/auth.go -destination ./tests/mock/auth/auth.go
//

// Package mock_auth is a generated GoMock package.
package mock_auth

import (
	context "context"
	reflect "reflect"

	auth "firebase.google.com/go/auth"
	auth0 "github.com/downsized-devs/sdk-go/auth"
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

// DeleteUser mocks base method.
func (m *MockInterface) DeleteUser(ctx context.Context, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockInterfaceMockRecorder) DeleteUser(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockInterface)(nil).DeleteUser), ctx, userID)
}

// GetUser mocks base method.
func (m *MockInterface) GetUser(ctx context.Context, userParam auth0.FirebaseUserParam) ([]auth0.FirebaseUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", ctx, userParam)
	ret0, _ := ret[0].([]auth0.FirebaseUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockInterfaceMockRecorder) GetUser(ctx, userParam any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockInterface)(nil).GetUser), ctx, userParam)
}

// GetUserAuthInfo mocks base method.
func (m *MockInterface) GetUserAuthInfo(ctx context.Context) (auth0.UserAuthInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserAuthInfo", ctx)
	ret0, _ := ret[0].(auth0.UserAuthInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserAuthInfo indicates an expected call of GetUserAuthInfo.
func (mr *MockInterfaceMockRecorder) GetUserAuthInfo(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserAuthInfo", reflect.TypeOf((*MockInterface)(nil).GetUserAuthInfo), ctx)
}

// GetUsers mocks base method.
func (m *MockInterface) GetUsers(ctx context.Context, userParams []auth0.FirebaseUserParam) ([]auth0.FirebaseUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers", ctx, userParams)
	ret0, _ := ret[0].([]auth0.FirebaseUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockInterfaceMockRecorder) GetUsers(ctx, userParams any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockInterface)(nil).GetUsers), ctx, userParams)
}

// RefreshToken mocks base method.
func (m *MockInterface) RefreshToken(ctx context.Context, refreshToken string) (auth0.RefreshTokenResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshToken", ctx, refreshToken)
	ret0, _ := ret[0].(auth0.RefreshTokenResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshToken indicates an expected call of RefreshToken.
func (mr *MockInterfaceMockRecorder) RefreshToken(ctx, refreshToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshToken", reflect.TypeOf((*MockInterface)(nil).RefreshToken), ctx, refreshToken)
}

// RegisterUser mocks base method.
func (m *MockInterface) RegisterUser(ctx context.Context, user auth0.FirebaseUser) (auth0.FirebaseUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUser", ctx, user)
	ret0, _ := ret[0].(auth0.FirebaseUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockInterfaceMockRecorder) RegisterUser(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockInterface)(nil).RegisterUser), ctx, user)
}

// RevokeUserRefreshToken mocks base method.
func (m *MockInterface) RevokeUserRefreshToken(ctx context.Context, uid string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RevokeUserRefreshToken", ctx, uid)
	ret0, _ := ret[0].(error)
	return ret0
}

// RevokeUserRefreshToken indicates an expected call of RevokeUserRefreshToken.
func (mr *MockInterfaceMockRecorder) RevokeUserRefreshToken(ctx, uid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RevokeUserRefreshToken", reflect.TypeOf((*MockInterface)(nil).RevokeUserRefreshToken), ctx, uid)
}

// SetUserAuthInfo mocks base method.
func (m *MockInterface) SetUserAuthInfo(ctx context.Context, param auth0.UserAuthParam) context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetUserAuthInfo", ctx, param)
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// SetUserAuthInfo indicates an expected call of SetUserAuthInfo.
func (mr *MockInterfaceMockRecorder) SetUserAuthInfo(ctx, param any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUserAuthInfo", reflect.TypeOf((*MockInterface)(nil).SetUserAuthInfo), ctx, param)
}

// SignInWithPassword mocks base method.
func (m *MockInterface) SignInWithPassword(ctx context.Context, param auth0.UserLogin) (auth0.UserLoginResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignInWithPassword", ctx, param)
	ret0, _ := ret[0].(auth0.UserLoginResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignInWithPassword indicates an expected call of SignInWithPassword.
func (mr *MockInterfaceMockRecorder) SignInWithPassword(ctx, param any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignInWithPassword", reflect.TypeOf((*MockInterface)(nil).SignInWithPassword), ctx, param)
}

// UpdateUser mocks base method.
func (m *MockInterface) UpdateUser(ctx context.Context, user auth0.FirebaseUser) (auth0.FirebaseUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, user)
	ret0, _ := ret[0].(auth0.FirebaseUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockInterfaceMockRecorder) UpdateUser(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockInterface)(nil).UpdateUser), ctx, user)
}

// VerifyPassword mocks base method.
func (m *MockInterface) VerifyPassword(ctx context.Context, email, password string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyPassword", ctx, email, password)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyPassword indicates an expected call of VerifyPassword.
func (mr *MockInterfaceMockRecorder) VerifyPassword(ctx, email, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyPassword", reflect.TypeOf((*MockInterface)(nil).VerifyPassword), ctx, email, password)
}

// VerifyToken mocks base method.
func (m *MockInterface) VerifyToken(ctx context.Context, bearertoken string) (*auth.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyToken", ctx, bearertoken)
	ret0, _ := ret[0].(*auth.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyToken indicates an expected call of VerifyToken.
func (mr *MockInterfaceMockRecorder) VerifyToken(ctx, bearertoken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyToken", reflect.TypeOf((*MockInterface)(nil).VerifyToken), ctx, bearertoken)
}