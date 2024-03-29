// Code generated by MockGen. DO NOT EDIT.
// Source: ./usecase/receive_event.go
//
// Generated by this command:
//
//	mockgen -source=./usecase/receive_event.go -destination=./mock/usecase/receive_event_mock.go
//

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	reflect "reflect"

	usecase "github.com/wahyusumartha/analytics-open-search/usecase"
	gomock "go.uber.org/mock/gomock"
)

// MockReceiveEvent is a mock of ReceiveEvent interface.
type MockReceiveEvent struct {
	ctrl     *gomock.Controller
	recorder *MockReceiveEventMockRecorder
}

// MockReceiveEventMockRecorder is the mock recorder for MockReceiveEvent.
type MockReceiveEventMockRecorder struct {
	mock *MockReceiveEvent
}

// NewMockReceiveEvent creates a new mock instance.
func NewMockReceiveEvent(ctrl *gomock.Controller) *MockReceiveEvent {
	mock := &MockReceiveEvent{ctrl: ctrl}
	mock.recorder = &MockReceiveEventMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReceiveEvent) EXPECT() *MockReceiveEventMockRecorder {
	return m.recorder
}

// Execute mocks base method.
func (m *MockReceiveEvent) Execute(ctx context.Context, input usecase.EventInput) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", ctx, input)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute.
func (mr *MockReceiveEventMockRecorder) Execute(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockReceiveEvent)(nil).Execute), ctx, input)
}
