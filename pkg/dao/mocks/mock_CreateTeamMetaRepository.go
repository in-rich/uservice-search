// Code generated by mockery v2.43.2. DO NOT EDIT.

package daomocks

import (
	context "context"

	entities "github.com/in-rich/uservice-search/pkg/entities"

	mock "github.com/stretchr/testify/mock"
)

// MockCreateTeamMetaRepository is an autogenerated mock type for the CreateTeamMetaRepository type
type MockCreateTeamMetaRepository struct {
	mock.Mock
}

type MockCreateTeamMetaRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCreateTeamMetaRepository) EXPECT() *MockCreateTeamMetaRepository_Expecter {
	return &MockCreateTeamMetaRepository_Expecter{mock: &_m.Mock}
}

// CreateTeamMeta provides a mock function with given fields: ctx, teamID, userID
func (_m *MockCreateTeamMetaRepository) CreateTeamMeta(ctx context.Context, teamID string, userID string) (*entities.TeamMeta, error) {
	ret := _m.Called(ctx, teamID, userID)

	if len(ret) == 0 {
		panic("no return value specified for CreateTeamMeta")
	}

	var r0 *entities.TeamMeta
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*entities.TeamMeta, error)); ok {
		return rf(ctx, teamID, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *entities.TeamMeta); ok {
		r0 = rf(ctx, teamID, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.TeamMeta)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, teamID, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCreateTeamMetaRepository_CreateTeamMeta_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateTeamMeta'
type MockCreateTeamMetaRepository_CreateTeamMeta_Call struct {
	*mock.Call
}

// CreateTeamMeta is a helper method to define mock.On call
//   - ctx context.Context
//   - teamID string
//   - userID string
func (_e *MockCreateTeamMetaRepository_Expecter) CreateTeamMeta(ctx interface{}, teamID interface{}, userID interface{}) *MockCreateTeamMetaRepository_CreateTeamMeta_Call {
	return &MockCreateTeamMetaRepository_CreateTeamMeta_Call{Call: _e.mock.On("CreateTeamMeta", ctx, teamID, userID)}
}

func (_c *MockCreateTeamMetaRepository_CreateTeamMeta_Call) Run(run func(ctx context.Context, teamID string, userID string)) *MockCreateTeamMetaRepository_CreateTeamMeta_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockCreateTeamMetaRepository_CreateTeamMeta_Call) Return(_a0 *entities.TeamMeta, _a1 error) *MockCreateTeamMetaRepository_CreateTeamMeta_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCreateTeamMetaRepository_CreateTeamMeta_Call) RunAndReturn(run func(context.Context, string, string) (*entities.TeamMeta, error)) *MockCreateTeamMetaRepository_CreateTeamMeta_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockCreateTeamMetaRepository creates a new instance of MockCreateTeamMetaRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCreateTeamMetaRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCreateTeamMetaRepository {
	mock := &MockCreateTeamMetaRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
