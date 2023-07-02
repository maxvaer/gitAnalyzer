// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	git "github.com/go-git/go-git/v5"
	mock "github.com/stretchr/testify/mock"
)

// IGitHelper is an autogenerated mock type for the IGitHelper type
type IGitHelper struct {
	mock.Mock
}

type IGitHelper_Expecter struct {
	mock *mock.Mock
}

func (_m *IGitHelper) EXPECT() *IGitHelper_Expecter {
	return &IGitHelper_Expecter{mock: &_m.Mock}
}

// Clone provides a mock function with given fields: dir, url
func (_m *IGitHelper) Clone(dir string, url string) (*git.Repository, error) {
	ret := _m.Called(dir, url)

	var r0 *git.Repository
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (*git.Repository, error)); ok {
		return rf(dir, url)
	}
	if rf, ok := ret.Get(0).(func(string, string) *git.Repository); ok {
		r0 = rf(dir, url)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*git.Repository)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(dir, url)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IGitHelper_Clone_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Clone'
type IGitHelper_Clone_Call struct {
	*mock.Call
}

// Clone is a helper method to define mock.On call
//   - dir string
//   - url string
func (_e *IGitHelper_Expecter) Clone(dir interface{}, url interface{}) *IGitHelper_Clone_Call {
	return &IGitHelper_Clone_Call{Call: _e.mock.On("Clone", dir, url)}
}

func (_c *IGitHelper_Clone_Call) Run(run func(dir string, url string)) *IGitHelper_Clone_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *IGitHelper_Clone_Call) Return(_a0 *git.Repository, _a1 error) *IGitHelper_Clone_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IGitHelper_Clone_Call) RunAndReturn(run func(string, string) (*git.Repository, error)) *IGitHelper_Clone_Call {
	_c.Call.Return(run)
	return _c
}

// Open provides a mock function with given fields: dir
func (_m *IGitHelper) Open(dir string) (*git.Repository, error) {
	ret := _m.Called(dir)

	var r0 *git.Repository
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*git.Repository, error)); ok {
		return rf(dir)
	}
	if rf, ok := ret.Get(0).(func(string) *git.Repository); ok {
		r0 = rf(dir)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*git.Repository)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(dir)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IGitHelper_Open_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Open'
type IGitHelper_Open_Call struct {
	*mock.Call
}

// Open is a helper method to define mock.On call
//   - dir string
func (_e *IGitHelper_Expecter) Open(dir interface{}) *IGitHelper_Open_Call {
	return &IGitHelper_Open_Call{Call: _e.mock.On("Open", dir)}
}

func (_c *IGitHelper_Open_Call) Run(run func(dir string)) *IGitHelper_Open_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *IGitHelper_Open_Call) Return(_a0 *git.Repository, _a1 error) *IGitHelper_Open_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IGitHelper_Open_Call) RunAndReturn(run func(string) (*git.Repository, error)) *IGitHelper_Open_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewIGitHelper interface {
	mock.TestingT
	Cleanup(func())
}

// NewIGitHelper creates a new instance of IGitHelper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIGitHelper(t mockConstructorTestingTNewIGitHelper) *IGitHelper {
	mock := &IGitHelper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
