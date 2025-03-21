// Code generated by mockery v2.43.2. DO NOT EDIT.

package servicesmocks

import (
	context "context"

	models "github.com/in-rich/uservice-search/pkg/models"
	mock "github.com/stretchr/testify/mock"
)

// MockSearchRemindersService is an autogenerated mock type for the SearchRemindersService type
type MockSearchRemindersService struct {
	mock.Mock
}

type MockSearchRemindersService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSearchRemindersService) EXPECT() *MockSearchRemindersService_Expecter {
	return &MockSearchRemindersService_Expecter{mock: &_m.Mock}
}

// Exec provides a mock function with given fields: ctx, reminder
func (_m *MockSearchRemindersService) Exec(ctx context.Context, reminder *models.SearchReminders) ([]*models.Reminder, error) {
	ret := _m.Called(ctx, reminder)

	if len(ret) == 0 {
		panic("no return value specified for Exec")
	}

	var r0 []*models.Reminder
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.SearchReminders) ([]*models.Reminder, error)); ok {
		return rf(ctx, reminder)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.SearchReminders) []*models.Reminder); ok {
		r0 = rf(ctx, reminder)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Reminder)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.SearchReminders) error); ok {
		r1 = rf(ctx, reminder)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSearchRemindersService_Exec_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Exec'
type MockSearchRemindersService_Exec_Call struct {
	*mock.Call
}

// Exec is a helper method to define mock.On call
//   - ctx context.Context
//   - reminder *models.SearchReminders
func (_e *MockSearchRemindersService_Expecter) Exec(ctx interface{}, reminder interface{}) *MockSearchRemindersService_Exec_Call {
	return &MockSearchRemindersService_Exec_Call{Call: _e.mock.On("Exec", ctx, reminder)}
}

func (_c *MockSearchRemindersService_Exec_Call) Run(run func(ctx context.Context, reminder *models.SearchReminders)) *MockSearchRemindersService_Exec_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*models.SearchReminders))
	})
	return _c
}

func (_c *MockSearchRemindersService_Exec_Call) Return(_a0 []*models.Reminder, _a1 error) *MockSearchRemindersService_Exec_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSearchRemindersService_Exec_Call) RunAndReturn(run func(context.Context, *models.SearchReminders) ([]*models.Reminder, error)) *MockSearchRemindersService_Exec_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockSearchRemindersService creates a new instance of MockSearchRemindersService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSearchRemindersService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSearchRemindersService {
	mock := &MockSearchRemindersService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
