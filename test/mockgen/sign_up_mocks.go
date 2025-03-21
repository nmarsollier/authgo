// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/usecases/sign_up.go

// Package mockgen is a generated GoMock package.
package mockgen

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	usecases "github.com/nmarsollier/authgo/internal/usecases"
)

// MockSignUpUseCase is a mock of SignUpUseCase interface.
type MockSignUpUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockSignUpUseCaseMockRecorder
}

// MockSignUpUseCaseMockRecorder is the mock recorder for MockSignUpUseCase.
type MockSignUpUseCaseMockRecorder struct {
	mock *MockSignUpUseCase
}

// NewMockSignUpUseCase creates a new mock instance.
func NewMockSignUpUseCase(ctrl *gomock.Controller) *MockSignUpUseCase {
	mock := &MockSignUpUseCase{ctrl: ctrl}
	mock.recorder = &MockSignUpUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSignUpUseCase) EXPECT() *MockSignUpUseCaseMockRecorder {
	return m.recorder
}

// SignUp mocks base method.
func (m *MockSignUpUseCase) SignUp(request *usecases.SignUpRequest) (*usecases.TokenResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", request)
	ret0, _ := ret[0].(*usecases.TokenResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignUp indicates an expected call of SignUp.
func (mr *MockSignUpUseCaseMockRecorder) SignUp(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockSignUpUseCase)(nil).SignUp), request)
}
