// Code generated by MockGen. DO NOT EDIT.
// Source: ./usecases/sign_in.go

// Package usecases is a generated GoMock package.
package usecases

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	usecases "github.com/nmarsollier/authgo/usecases"
)

// MockSignInUseCase is a mock of SignInUseCase interface.
type MockSignInUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockSignInUseCaseMockRecorder
}

// MockSignInUseCaseMockRecorder is the mock recorder for MockSignInUseCase.
type MockSignInUseCaseMockRecorder struct {
	mock *MockSignInUseCase
}

// NewMockSignInUseCase creates a new mock instance.
func NewMockSignInUseCase(ctrl *gomock.Controller) *MockSignInUseCase {
	mock := &MockSignInUseCase{ctrl: ctrl}
	mock.recorder = &MockSignInUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSignInUseCase) EXPECT() *MockSignInUseCaseMockRecorder {
	return m.recorder
}

// SignIn mocks base method.
func (m *MockSignInUseCase) SignIn(request *usecases.SignInRequest) (*usecases.TokenResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIn", request)
	ret0, _ := ret[0].(*usecases.TokenResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignIn indicates an expected call of SignIn.
func (mr *MockSignInUseCaseMockRecorder) SignIn(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIn", reflect.TypeOf((*MockSignInUseCase)(nil).SignIn), request)
}
