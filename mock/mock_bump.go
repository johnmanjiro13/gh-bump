// Code generated by MockGen. DO NOT EDIT.
// Source: bumper.go

// Package mock is a generated GoMock package.
package mock

import (
	bytes "bytes"
	reflect "reflect"

	survey "github.com/AlecAivazis/survey/v2"
	gomock "github.com/golang/mock/gomock"
	bump "github.com/johnmanjiro13/gh-bump"
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
func (m *MockGh) CreateRelease(version, repo string, isCurrent bool, option *bump.ReleaseOption) (*bytes.Buffer, *bytes.Buffer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRelease", version, repo, isCurrent, option)
	ret0, _ := ret[0].(*bytes.Buffer)
	ret1, _ := ret[1].(*bytes.Buffer)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateRelease indicates an expected call of CreateRelease.
func (mr *MockGhMockRecorder) CreateRelease(version, repo, isCurrent, option interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRelease", reflect.TypeOf((*MockGh)(nil).CreateRelease), version, repo, isCurrent, option)
}

// ListRelease mocks base method.
func (m *MockGh) ListRelease(repo string, isCurrent bool) (*bytes.Buffer, *bytes.Buffer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRelease", repo, isCurrent)
	ret0, _ := ret[0].(*bytes.Buffer)
	ret1, _ := ret[1].(*bytes.Buffer)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListRelease indicates an expected call of ListRelease.
func (mr *MockGhMockRecorder) ListRelease(repo, isCurrent interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRelease", reflect.TypeOf((*MockGh)(nil).ListRelease), repo, isCurrent)
}

// ViewRelease mocks base method.
func (m *MockGh) ViewRelease(repo string, isCurrent bool) (*bytes.Buffer, *bytes.Buffer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ViewRelease", repo, isCurrent)
	ret0, _ := ret[0].(*bytes.Buffer)
	ret1, _ := ret[1].(*bytes.Buffer)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ViewRelease indicates an expected call of ViewRelease.
func (mr *MockGhMockRecorder) ViewRelease(repo, isCurrent interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ViewRelease", reflect.TypeOf((*MockGh)(nil).ViewRelease), repo, isCurrent)
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

// MockPrompter is a mock of Prompter interface.
type MockPrompter struct {
	ctrl     *gomock.Controller
	recorder *MockPrompterMockRecorder
}

// MockPrompterMockRecorder is the mock recorder for MockPrompter.
type MockPrompterMockRecorder struct {
	mock *MockPrompter
}

// NewMockPrompter creates a new mock instance.
func NewMockPrompter(ctrl *gomock.Controller) *MockPrompter {
	mock := &MockPrompter{ctrl: ctrl}
	mock.recorder = &MockPrompterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPrompter) EXPECT() *MockPrompterMockRecorder {
	return m.recorder
}

// Confirm mocks base method.
func (m *MockPrompter) Confirm(question string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Confirm", question)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Confirm indicates an expected call of Confirm.
func (mr *MockPrompterMockRecorder) Confirm(question interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Confirm", reflect.TypeOf((*MockPrompter)(nil).Confirm), question)
}

// Input mocks base method.
func (m *MockPrompter) Input(question string, validator survey.Validator) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Input", question, validator)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Input indicates an expected call of Input.
func (mr *MockPrompterMockRecorder) Input(question, validator interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Input", reflect.TypeOf((*MockPrompter)(nil).Input), question, validator)
}

// Select mocks base method.
func (m *MockPrompter) Select(question string, options []string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Select", question, options)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Select indicates an expected call of Select.
func (mr *MockPrompterMockRecorder) Select(question, options interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Select", reflect.TypeOf((*MockPrompter)(nil).Select), question, options)
}
