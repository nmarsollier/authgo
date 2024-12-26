// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/rabbit/send_logout.go

// Package mockgen is a generated GoMock package.
package mockgen

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSendLogoutService is a mock of SendLogoutService interface.
type MockSendLogoutService struct {
	ctrl     *gomock.Controller
	recorder *MockSendLogoutServiceMockRecorder
}

// MockSendLogoutServiceMockRecorder is the mock recorder for MockSendLogoutService.
type MockSendLogoutServiceMockRecorder struct {
	mock *MockSendLogoutService
}

// NewMockSendLogoutService creates a new mock instance.
func NewMockSendLogoutService(ctrl *gomock.Controller) *MockSendLogoutService {
	mock := &MockSendLogoutService{ctrl: ctrl}
	mock.recorder = &MockSendLogoutServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSendLogoutService) EXPECT() *MockSendLogoutServiceMockRecorder {
	return m.recorder
}

// SendLogout mocks base method.
func (m *MockSendLogoutService) SendLogout(token string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendLogout", token)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendLogout indicates an expected call of SendLogout.
func (mr *MockSendLogoutServiceMockRecorder) SendLogout(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendLogout", reflect.TypeOf((*MockSendLogoutService)(nil).SendLogout), token)
}
