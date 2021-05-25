// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package eventmocks

import (
	blockchain "github.com/kaleido-io/firefly/pkg/blockchain"

	fftypes "github.com/kaleido-io/firefly/pkg/fftypes"

	mock "github.com/stretchr/testify/mock"
)

// EventManager is an autogenerated mock type for the EventManager type
type EventManager struct {
	mock.Mock
}

// DeletedSubscriptions provides a mock function with given fields:
func (_m *EventManager) DeletedSubscriptions() chan<- *fftypes.UUID {
	ret := _m.Called()

	var r0 chan<- *fftypes.UUID
	if rf, ok := ret.Get(0).(func() chan<- *fftypes.UUID); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan<- *fftypes.UUID)
		}
	}

	return r0
}

// NewEvents provides a mock function with given fields:
func (_m *EventManager) NewEvents() chan<- int64 {
	ret := _m.Called()

	var r0 chan<- int64
	if rf, ok := ret.Get(0).(func() chan<- int64); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan<- int64)
		}
	}

	return r0
}

// NewSubscriptions provides a mock function with given fields:
func (_m *EventManager) NewSubscriptions() chan<- *fftypes.UUID {
	ret := _m.Called()

	var r0 chan<- *fftypes.UUID
	if rf, ok := ret.Get(0).(func() chan<- *fftypes.UUID); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan<- *fftypes.UUID)
		}
	}

	return r0
}

// SequencedBroadcastBatch provides a mock function with given fields: batch, author, protocolTxID, additionalInfo
func (_m *EventManager) SequencedBroadcastBatch(batch *blockchain.BroadcastBatch, author string, protocolTxID string, additionalInfo map[string]interface{}) error {
	ret := _m.Called(batch, author, protocolTxID, additionalInfo)

	var r0 error
	if rf, ok := ret.Get(0).(func(*blockchain.BroadcastBatch, string, string, map[string]interface{}) error); ok {
		r0 = rf(batch, author, protocolTxID, additionalInfo)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Start provides a mock function with given fields:
func (_m *EventManager) Start() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TransactionUpdate provides a mock function with given fields: txTrackingID, txState, protocolTxID, errorMessage, additionalInfo
func (_m *EventManager) TransactionUpdate(txTrackingID string, txState fftypes.TransactionStatus, protocolTxID string, errorMessage string, additionalInfo map[string]interface{}) error {
	ret := _m.Called(txTrackingID, txState, protocolTxID, errorMessage, additionalInfo)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, fftypes.TransactionStatus, string, string, map[string]interface{}) error); ok {
		r0 = rf(txTrackingID, txState, protocolTxID, errorMessage, additionalInfo)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WaitStop provides a mock function with given fields:
func (_m *EventManager) WaitStop() {
	_m.Called()
}
