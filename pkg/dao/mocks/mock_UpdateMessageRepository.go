// Code generated by mockery v2.43.2. DO NOT EDIT.

package daomocks

import (
	context "context"

	dao "github.com/in-rich/uservice-search/pkg/dao"
	entities "github.com/in-rich/uservice-search/pkg/entities"

	mock "github.com/stretchr/testify/mock"
)

// MockUpdateMessageRepository is an autogenerated mock type for the UpdateMessageRepository type
type MockUpdateMessageRepository struct {
	mock.Mock
}

type MockUpdateMessageRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUpdateMessageRepository) EXPECT() *MockUpdateMessageRepository_Expecter {
	return &MockUpdateMessageRepository_Expecter{mock: &_m.Mock}
}

// UpdateMessage provides a mock function with given fields: ctx, teamID, messageID, data
func (_m *MockUpdateMessageRepository) UpdateMessage(ctx context.Context, teamID string, messageID string, data *dao.UpdateMessageData) (*entities.Message, error) {
	ret := _m.Called(ctx, teamID, messageID, data)

	if len(ret) == 0 {
		panic("no return value specified for UpdateMessage")
	}

	var r0 *entities.Message
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *dao.UpdateMessageData) (*entities.Message, error)); ok {
		return rf(ctx, teamID, messageID, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *dao.UpdateMessageData) *entities.Message); ok {
		r0 = rf(ctx, teamID, messageID, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Message)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, *dao.UpdateMessageData) error); ok {
		r1 = rf(ctx, teamID, messageID, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUpdateMessageRepository_UpdateMessage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateMessage'
type MockUpdateMessageRepository_UpdateMessage_Call struct {
	*mock.Call
}

// UpdateMessage is a helper method to define mock.On call
//   - ctx context.Context
//   - teamID string
//   - messageID string
//   - data *dao.UpdateMessageData
func (_e *MockUpdateMessageRepository_Expecter) UpdateMessage(ctx interface{}, teamID interface{}, messageID interface{}, data interface{}) *MockUpdateMessageRepository_UpdateMessage_Call {
	return &MockUpdateMessageRepository_UpdateMessage_Call{Call: _e.mock.On("UpdateMessage", ctx, teamID, messageID, data)}
}

func (_c *MockUpdateMessageRepository_UpdateMessage_Call) Run(run func(ctx context.Context, teamID string, messageID string, data *dao.UpdateMessageData)) *MockUpdateMessageRepository_UpdateMessage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(*dao.UpdateMessageData))
	})
	return _c
}

func (_c *MockUpdateMessageRepository_UpdateMessage_Call) Return(_a0 *entities.Message, _a1 error) *MockUpdateMessageRepository_UpdateMessage_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUpdateMessageRepository_UpdateMessage_Call) RunAndReturn(run func(context.Context, string, string, *dao.UpdateMessageData) (*entities.Message, error)) *MockUpdateMessageRepository_UpdateMessage_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockUpdateMessageRepository creates a new instance of MockUpdateMessageRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUpdateMessageRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUpdateMessageRepository {
	mock := &MockUpdateMessageRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
