// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/typical-go/typical-go/pkg/typapp (interfaces: Preparer)

// Package mock_typapp is a generated GoMock package.
package mock_typapp

import (
	gomock "github.com/golang/mock/gomock"
	typapp "github.com/typical-go/typical-go/pkg/typapp"
	reflect "reflect"
)

// MockPreparer is a mock of Preparer interface
type MockPreparer struct {
	ctrl     *gomock.Controller
	recorder *MockPreparerMockRecorder
}

// MockPreparerMockRecorder is the mock recorder for MockPreparer
type MockPreparerMockRecorder struct {
	mock *MockPreparer
}

// NewMockPreparer creates a new mock instance
func NewMockPreparer(ctrl *gomock.Controller) *MockPreparer {
	mock := &MockPreparer{ctrl: ctrl}
	mock.recorder = &MockPreparerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPreparer) EXPECT() *MockPreparerMockRecorder {
	return m.recorder
}

// Preparations mocks base method
func (m *MockPreparer) Preparations() []*typapp.Preparation {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Preparations")
	ret0, _ := ret[0].([]*typapp.Preparation)
	return ret0
}

// Preparations indicates an expected call of Preparations
func (mr *MockPreparerMockRecorder) Preparations() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Preparations", reflect.TypeOf((*MockPreparer)(nil).Preparations))
}