// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/repository/session.go
//
// Generated by this command:
//
//	mockgen -source=./internal/repository/session.go -destination=./tests/mock/mock_session_repository.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockISessionRepo is a mock of ISessionRepo interface.
type MockISessionRepo struct {
	ctrl     *gomock.Controller
	recorder *MockISessionRepoMockRecorder
	isgomock struct{}
}

// MockISessionRepoMockRecorder is the mock recorder for MockISessionRepo.
type MockISessionRepoMockRecorder struct {
	mock *MockISessionRepo
}

// NewMockISessionRepo creates a new mock instance.
func NewMockISessionRepo(ctrl *gomock.Controller) *MockISessionRepo {
	mock := &MockISessionRepo{ctrl: ctrl}
	mock.recorder = &MockISessionRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockISessionRepo) EXPECT() *MockISessionRepoMockRecorder {
	return m.recorder
}

// CreateSession mocks base method.
func (m *MockISessionRepo) CreateSession(ctx context.Context, userID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", ctx, userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockISessionRepoMockRecorder) CreateSession(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockISessionRepo)(nil).CreateSession), ctx, userID)
}

// DeleteSession mocks base method.
func (m *MockISessionRepo) DeleteSession(ctx context.Context, sessionId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSession", ctx, sessionId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSession indicates an expected call of DeleteSession.
func (mr *MockISessionRepoMockRecorder) DeleteSession(ctx, sessionId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSession", reflect.TypeOf((*MockISessionRepo)(nil).DeleteSession), ctx, sessionId)
}

// GetUserIDByToken mocks base method.
func (m *MockISessionRepo) GetUserIDByToken(ctx context.Context, sessId string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserIDByToken", ctx, sessId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserIDByToken indicates an expected call of GetUserIDByToken.
func (mr *MockISessionRepoMockRecorder) GetUserIDByToken(ctx, sessId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserIDByToken", reflect.TypeOf((*MockISessionRepo)(nil).GetUserIDByToken), ctx, sessId)
}
