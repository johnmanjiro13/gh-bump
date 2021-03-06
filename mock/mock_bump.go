// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/johnmanjiro13/gh-bump/bump (interfaces: Gh)

// Package mock is a generated GoMock package.
package mock

import (
	bytes "bytes"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	bump "github.com/johnmanjiro13/gh-bump/bump"
)

// MockGh is a mock of Gh interface.
type MockGh struct {
	ctrl     *gomock.Controller
	recorder *MockGhMockRecorder
}

// MockGhMockRecorder is the mock recorder for MockGh.
type MockGhMockRecorder struct {
	mock *MockGh
}

// NewMockGh creates a new mock instance.
func NewMockGh(ctrl *gomock.Controller) *MockGh {
	mock := &MockGh{ctrl: ctrl}
	mock.recorder = &MockGhMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGh) EXPECT() *MockGhMockRecorder {
	return m.recorder
}

// CreateRelease mocks base method.
func (m *MockGh) CreateRelease(arg0, arg1 string, arg2 bool, arg3 *bump.ReleaseOption) (*bytes.Buffer, *bytes.Buffer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRelease", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*bytes.Buffer)
	ret1, _ := ret[1].(*bytes.Buffer)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateRelease indicates an expected call of CreateRelease.
func (mr *MockGhMockRecorder) CreateRelease(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRelease", reflect.TypeOf((*MockGh)(nil).CreateRelease), arg0, arg1, arg2, arg3)
}

// ListRelease mocks base method.
func (m *MockGh) ListRelease(arg0 string, arg1 bool) (*bytes.Buffer, *bytes.Buffer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRelease", arg0, arg1)
	ret0, _ := ret[0].(*bytes.Buffer)
	ret1, _ := ret[1].(*bytes.Buffer)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListRelease indicates an expected call of ListRelease.
func (mr *MockGhMockRecorder) ListRelease(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRelease", reflect.TypeOf((*MockGh)(nil).ListRelease), arg0, arg1)
}

// ViewRelease mocks base method.
func (m *MockGh) ViewRelease(arg0 string, arg1 bool) (*bytes.Buffer, *bytes.Buffer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ViewRelease", arg0, arg1)
	ret0, _ := ret[0].(*bytes.Buffer)
	ret1, _ := ret[1].(*bytes.Buffer)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ViewRelease indicates an expected call of ViewRelease.
func (mr *MockGhMockRecorder) ViewRelease(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ViewRelease", reflect.TypeOf((*MockGh)(nil).ViewRelease), arg0, arg1)
}

// ViewRepository mocks base method.
func (m *MockGh) ViewRepository() (*bytes.Buffer, *bytes.Buffer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ViewRepository")
	ret0, _ := ret[0].(*bytes.Buffer)
	ret1, _ := ret[1].(*bytes.Buffer)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ViewRepository indicates an expected call of ViewRepository.
func (mr *MockGhMockRecorder) ViewRepository() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ViewRepository", reflect.TypeOf((*MockGh)(nil).ViewRepository))
}
