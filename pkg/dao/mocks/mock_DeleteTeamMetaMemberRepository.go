// Code generated by mockery v2.43.2. DO NOT EDIT.

package daomocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockDeleteTeamMetaMemberRepository is an autogenerated mock type for the DeleteTeamMetaMemberRepository type
type MockDeleteTeamMetaMemberRepository struct {
	mock.Mock
}

type MockDeleteTeamMetaMemberRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDeleteTeamMetaMemberRepository) EXPECT() *MockDeleteTeamMetaMemberRepository_Expecter {
	return &MockDeleteTeamMetaMemberRepository_Expecter{mock: &_m.Mock}
}

// DeleteTeamMetaMember provides a mock function with given fields: ctx, teamID, userID
func (_m *MockDeleteTeamMetaMemberRepository) DeleteTeamMetaMember(ctx context.Context, teamID string, userID string) error {
	ret := _m.Called(ctx, teamID, userID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteTeamMetaMember")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, teamID, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDeleteTeamMetaMemberRepository_DeleteTeamMetaMember_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteTeamMetaMember'
type MockDeleteTeamMetaMemberRepository_DeleteTeamMetaMember_Call struct {
	*mock.Call
}

// DeleteTeamMetaMember is a helper method to define mock.On call
//   - ctx context.Context
//   - teamID string
//   - userID string
func (_e *MockDeleteTeamMetaMemberRepository_Expecter) DeleteTeamMetaMember(ctx interface{}, teamID interface{}, userID interface{}) *MockDeleteTeamMetaMemberRepository_DeleteTeamMetaMember_Call {
	return &MockDeleteTeamMetaMemberRepository_DeleteTeamMetaMember_Call{Call: _e.mock.On("DeleteTeamMetaMember", ctx, teamID, userID)}
}

func (_c *MockDeleteTeamMetaMemberRepository_DeleteTeamMetaMember_Call) Run(run func(ctx context.Context, teamID string, userID string)) *MockDeleteTeamMetaMemberRepository_DeleteTeamMetaMember_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockDeleteTeamMetaMemberRepository_DeleteTeamMetaMember_Call) Return(_a0 error) *MockDeleteTeamMetaMemberRepository_DeleteTeamMetaMember_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDeleteTeamMetaMemberRepository_DeleteTeamMetaMember_Call) RunAndReturn(run func(context.Context, string, string) error) *MockDeleteTeamMetaMemberRepository_DeleteTeamMetaMember_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDeleteTeamMetaMemberRepository creates a new instance of MockDeleteTeamMetaMemberRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDeleteTeamMetaMemberRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDeleteTeamMetaMemberRepository {
	mock := &MockDeleteTeamMetaMemberRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}