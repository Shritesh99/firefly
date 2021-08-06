// Code generated by mockery v1.0.0. DO NOT EDIT.

package sysmessagingmocks

import (
	context "context"

	fftypes "github.com/hyperledger-labs/firefly/pkg/fftypes"
	mock "github.com/stretchr/testify/mock"
)

// MessageSender is an autogenerated mock type for the MessageSender type
type MessageSender struct {
	mock.Mock
}

// SendMessageWithID provides a mock function with given fields: ctx, ns, in, waitConfirm
func (_m *MessageSender) SendMessageWithID(ctx context.Context, ns string, in *fftypes.MessageInOut, waitConfirm bool) (*fftypes.Message, error) {
	ret := _m.Called(ctx, ns, in, waitConfirm)

	var r0 *fftypes.Message
	if rf, ok := ret.Get(0).(func(context.Context, string, *fftypes.MessageInOut, bool) *fftypes.Message); ok {
		r0 = rf(ctx, ns, in, waitConfirm)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fftypes.Message)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *fftypes.MessageInOut, bool) error); ok {
		r1 = rf(ctx, ns, in, waitConfirm)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
