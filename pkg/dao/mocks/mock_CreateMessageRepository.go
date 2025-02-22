// Code generated by mockery v2.43.2. DO NOT EDIT.

package daomocks

import (
	context "context"

	dao "github.com/in-rich/uservice-search/pkg/dao"
	entities "github.com/in-rich/uservice-search/pkg/entities"

	mock "github.com/stretchr/testify/mock"
)

// MockCreateMessageRepository is an autogenerated mock type for the CreateMessageRepository type
type MockCreateMessageRepository struct {
	mock.Mock
}

type MockCreateMessageRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCreateMessageRepository) EXPECT() *MockCreateMessageRepository_Expecter {
	return &MockCreateMessageRepository_Expecter{mock: &_m.Mock}
}

// CreateMessage provides a mock function with given fields: ctx, teamID, messageID, data
func (_m *MockCreateMessageRepository) CreateMessage(ctx context.Context, teamID string, messageID string, data *dao.CreateMessageData) (*entities.Message, error) {
	ret := _m.Called(ctx, teamID, messageID, data)

	if len(ret) == 0 {
		panic("no return value specified for CreateMessage")
	}

	var r0 *entities.Message
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *dao.CreateMessageData) (*entities.Message, error)); ok {
		return rf(ctx, teamID, messageID, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *dao.CreateMessageData) *entities.Message); ok {
		r0 = rf(ctx, teamID, messageID, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Message)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, *dao.CreateMessageData) error); ok {
		r1 = rf(ctx, teamID, messageID, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCreateMessageRepository_CreateMessage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateMessage'
type MockCreateMessageRepository_CreateMessage_Call struct {
	*mock.Call
}

// CreateMessage is a helper method to define mock.On call
//   - ctx context.Context
//   - teamID string
//   - messageID string
//   - data *dao.CreateMessageData
func (_e *MockCreateMessageRepository_Expecter) CreateMessage(ctx interface{}, teamID interface{}, messageID interface{}, data interface{}) *MockCreateMessageRepository_CreateMessage_Call {
	return &MockCreateMessageRepository_CreateMessage_Call{Call: _e.mock.On("CreateMessage", ctx, teamID, messageID, data)}
}

func (_c *MockCreateMessageRepository_CreateMessage_Call) Run(run func(ctx context.Context, teamID string, messageID string, data *dao.CreateMessageData)) *MockCreateMessageRepository_CreateMessage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(*dao.CreateMessageData))
	})
	return _c
}

func (_c *MockCreateMessageRepository_CreateMessage_Call) Return(_a0 *entities.Message, _a1 error) *MockCreateMessageRepository_CreateMessage_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCreateMessageRepository_CreateMessage_Call) RunAndReturn(run func(context.Context, string, string, *dao.CreateMessageData) (*entities.Message, error)) *MockCreateMessageRepository_CreateMessage_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockCreateMessageRepository creates a new instance of MockCreateMessageRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCreateMessageRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCreateMessageRepository {
	mock := &MockCreateMessageRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
