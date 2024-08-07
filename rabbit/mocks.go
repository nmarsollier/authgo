// Code generated by MockGen. DO NOT EDIT.
// Source: ./rabbit/rabbit.go

// Package rabbit is a generated GoMock package.
package rabbit

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRabbit is a mock of Rabbit interface.
type MockRabbit struct {
	ctrl     *gomock.Controller
	recorder *MockRabbitMockRecorder
}

// MockRabbitMockRecorder is the mock recorder for MockRabbit.
type MockRabbitMockRecorder struct {
	mock *MockRabbit
}

// NewMockRabbit creates a new mock instance.
func NewMockRabbit(ctrl *gomock.Controller) *MockRabbit {
	mock := &MockRabbit{ctrl: ctrl}
	mock.recorder = &MockRabbitMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRabbit) EXPECT() *MockRabbitMockRecorder {
	return m.recorder
}

// SendLogout mocks base method.
func (m *MockRabbit) SendLogout(token string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendLogout", token)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendLogout indicates an expected call of SendLogout.
func (mr *MockRabbitMockRecorder) SendLogout(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendLogout", reflect.TypeOf((*MockRabbit)(nil).SendLogout), token)
}