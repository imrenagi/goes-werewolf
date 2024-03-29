// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/imrenagi/goes-werewolf/internal/app/werewolf/services (interfaces: PollDAO)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	models "github.com/imrenagi/goes-werewolf/internal/app/werewolf/models"
	reflect "reflect"
)

// MockPollDAO is a mock of PollDAO interface
type MockPollDAO struct {
	ctrl     *gomock.Controller
	recorder *MockPollDAOMockRecorder
}

// MockPollDAOMockRecorder is the mock recorder for MockPollDAO
type MockPollDAOMockRecorder struct {
	mock *MockPollDAO
}

// NewMockPollDAO creates a new mock instance
func NewMockPollDAO(ctrl *gomock.Controller) *MockPollDAO {
	mock := &MockPollDAO{ctrl: ctrl}
	mock.recorder = &MockPollDAOMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPollDAO) EXPECT() *MockPollDAOMockRecorder {
	return m.recorder
}

// GetPolls mocks base method
func (m *MockPollDAO) GetPolls(arg0, arg1 string) ([]models.Poll, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPolls", arg0, arg1)
	ret0, _ := ret[0].([]models.Poll)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPolls indicates an expected call of GetPolls
func (mr *MockPollDAOMockRecorder) GetPolls(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPolls", reflect.TypeOf((*MockPollDAO)(nil).GetPolls), arg0, arg1)
}
