// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/typical-go/typical-go/pkg/typgo (interfaces: Destroyer)

// Package mock_typgo is a generated GoMock package.
package mock_typgo

import (
	gomock "github.com/golang/mock/gomock"
	typgo "github.com/typical-go/typical-go/pkg/typgo"
	reflect "reflect"
)

// MockDestroyer is a mock of Destroyer interface
type MockDestroyer struct {
	ctrl     *gomock.Controller
	recorder *MockDestroyerMockRecorder
}

// MockDestroyerMockRecorder is the mock recorder for MockDestroyer
type MockDestroyerMockRecorder struct {
	mock *MockDestroyer
}

// NewMockDestroyer creates a new mock instance
func NewMockDestroyer(ctrl *gomock.Controller) *MockDestroyer {
	mock := &MockDestroyer{ctrl: ctrl}
	mock.recorder = &MockDestroyerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDestroyer) EXPECT() *MockDestroyerMockRecorder {
	return m.recorder
}

// Destructors mocks base method
func (m *MockDestroyer) Destructors() []*typgo.Destructor {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Destructors")
	ret0, _ := ret[0].([]*typgo.Destructor)
	return ret0
}

// Destructors indicates an expected call of Destructors
func (mr *MockDestroyerMockRecorder) Destructors() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Destructors", reflect.TypeOf((*MockDestroyer)(nil).Destructors))
}
