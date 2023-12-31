// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	Analyzer "GitAnalyzer/api/Analyzer"

	git "github.com/go-git/go-git/v5"

	mock "github.com/stretchr/testify/mock"

	object "github.com/go-git/go-git/v5/plumbing/object"

	plumbing "github.com/go-git/go-git/v5/plumbing"
)

// IRepoHelper is an autogenerated mock type for the IRepoHelper type
type IRepoHelper struct {
	mock.Mock
}

type IRepoHelper_Expecter struct {
	mock *mock.Mock
}

func (_m *IRepoHelper) EXPECT() *IRepoHelper_Expecter {
	return &IRepoHelper_Expecter{mock: &_m.Mock}
}

// Checkout provides a mock function with given fields: repo, opts
func (_m *IRepoHelper) Checkout(repo *git.Repository, opts *git.CheckoutOptions) error {
	ret := _m.Called(repo, opts)

	var r0 error
	if rf, ok := ret.Get(0).(func(*git.Repository, *git.CheckoutOptions) error); ok {
		r0 = rf(repo, opts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IRepoHelper_Checkout_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Checkout'
type IRepoHelper_Checkout_Call struct {
	*mock.Call
}

// Checkout is a helper method to define mock.On call
//   - repo *git.Repository
//   - opts *git.CheckoutOptions
func (_e *IRepoHelper_Expecter) Checkout(repo interface{}, opts interface{}) *IRepoHelper_Checkout_Call {
	return &IRepoHelper_Checkout_Call{Call: _e.mock.On("Checkout", repo, opts)}
}

func (_c *IRepoHelper_Checkout_Call) Run(run func(repo *git.Repository, opts *git.CheckoutOptions)) *IRepoHelper_Checkout_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*git.Repository), args[1].(*git.CheckoutOptions))
	})
	return _c
}

func (_c *IRepoHelper_Checkout_Call) Return(_a0 error) *IRepoHelper_Checkout_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IRepoHelper_Checkout_Call) RunAndReturn(run func(*git.Repository, *git.CheckoutOptions) error) *IRepoHelper_Checkout_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteRepository provides a mock function with given fields: repo
func (_m *IRepoHelper) DeleteRepository(repo *git.Repository) {
	_m.Called(repo)
}

// IRepoHelper_DeleteRepository_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteRepository'
type IRepoHelper_DeleteRepository_Call struct {
	*mock.Call
}

// DeleteRepository is a helper method to define mock.On call
//   - repo *git.Repository
func (_e *IRepoHelper_Expecter) DeleteRepository(repo interface{}) *IRepoHelper_DeleteRepository_Call {
	return &IRepoHelper_DeleteRepository_Call{Call: _e.mock.On("DeleteRepository", repo)}
}

func (_c *IRepoHelper_DeleteRepository_Call) Run(run func(repo *git.Repository)) *IRepoHelper_DeleteRepository_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*git.Repository))
	})
	return _c
}

func (_c *IRepoHelper_DeleteRepository_Call) Return() *IRepoHelper_DeleteRepository_Call {
	_c.Call.Return()
	return _c
}

func (_c *IRepoHelper_DeleteRepository_Call) RunAndReturn(run func(*git.Repository)) *IRepoHelper_DeleteRepository_Call {
	_c.Call.Return(run)
	return _c
}

// GetGitHubURLOfRepository provides a mock function with given fields: repo
func (_m *IRepoHelper) GetGitHubURLOfRepository(repo *git.Repository) string {
	ret := _m.Called(repo)

	var r0 string
	if rf, ok := ret.Get(0).(func(*git.Repository) string); ok {
		r0 = rf(repo)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// IRepoHelper_GetGitHubURLOfRepository_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetGitHubURLOfRepository'
type IRepoHelper_GetGitHubURLOfRepository_Call struct {
	*mock.Call
}

// GetGitHubURLOfRepository is a helper method to define mock.On call
//   - repo *git.Repository
func (_e *IRepoHelper_Expecter) GetGitHubURLOfRepository(repo interface{}) *IRepoHelper_GetGitHubURLOfRepository_Call {
	return &IRepoHelper_GetGitHubURLOfRepository_Call{Call: _e.mock.On("GetGitHubURLOfRepository", repo)}
}

func (_c *IRepoHelper_GetGitHubURLOfRepository_Call) Run(run func(repo *git.Repository)) *IRepoHelper_GetGitHubURLOfRepository_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*git.Repository))
	})
	return _c
}

func (_c *IRepoHelper_GetGitHubURLOfRepository_Call) Return(_a0 string) *IRepoHelper_GetGitHubURLOfRepository_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IRepoHelper_GetGitHubURLOfRepository_Call) RunAndReturn(run func(*git.Repository) string) *IRepoHelper_GetGitHubURLOfRepository_Call {
	_c.Call.Return(run)
	return _c
}

// GetHeadCommit provides a mock function with given fields: repo
func (_m *IRepoHelper) GetHeadCommit(repo *git.Repository) (*object.Commit, error) {
	ret := _m.Called(repo)

	var r0 *object.Commit
	var r1 error
	if rf, ok := ret.Get(0).(func(*git.Repository) (*object.Commit, error)); ok {
		return rf(repo)
	}
	if rf, ok := ret.Get(0).(func(*git.Repository) *object.Commit); ok {
		r0 = rf(repo)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*object.Commit)
		}
	}

	if rf, ok := ret.Get(1).(func(*git.Repository) error); ok {
		r1 = rf(repo)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IRepoHelper_GetHeadCommit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetHeadCommit'
type IRepoHelper_GetHeadCommit_Call struct {
	*mock.Call
}

// GetHeadCommit is a helper method to define mock.On call
//   - repo *git.Repository
func (_e *IRepoHelper_Expecter) GetHeadCommit(repo interface{}) *IRepoHelper_GetHeadCommit_Call {
	return &IRepoHelper_GetHeadCommit_Call{Call: _e.mock.On("GetHeadCommit", repo)}
}

func (_c *IRepoHelper_GetHeadCommit_Call) Run(run func(repo *git.Repository)) *IRepoHelper_GetHeadCommit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*git.Repository))
	})
	return _c
}

func (_c *IRepoHelper_GetHeadCommit_Call) Return(_a0 *object.Commit, _a1 error) *IRepoHelper_GetHeadCommit_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IRepoHelper_GetHeadCommit_Call) RunAndReturn(run func(*git.Repository) (*object.Commit, error)) *IRepoHelper_GetHeadCommit_Call {
	_c.Call.Return(run)
	return _c
}

// GetPathOfRepository provides a mock function with given fields: repo
func (_m *IRepoHelper) GetPathOfRepository(repo *git.Repository) string {
	ret := _m.Called(repo)

	var r0 string
	if rf, ok := ret.Get(0).(func(*git.Repository) string); ok {
		r0 = rf(repo)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// IRepoHelper_GetPathOfRepository_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPathOfRepository'
type IRepoHelper_GetPathOfRepository_Call struct {
	*mock.Call
}

// GetPathOfRepository is a helper method to define mock.On call
//   - repo *git.Repository
func (_e *IRepoHelper_Expecter) GetPathOfRepository(repo interface{}) *IRepoHelper_GetPathOfRepository_Call {
	return &IRepoHelper_GetPathOfRepository_Call{Call: _e.mock.On("GetPathOfRepository", repo)}
}

func (_c *IRepoHelper_GetPathOfRepository_Call) Run(run func(repo *git.Repository)) *IRepoHelper_GetPathOfRepository_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*git.Repository))
	})
	return _c
}

func (_c *IRepoHelper_GetPathOfRepository_Call) Return(_a0 string) *IRepoHelper_GetPathOfRepository_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IRepoHelper_GetPathOfRepository_Call) RunAndReturn(run func(*git.Repository) string) *IRepoHelper_GetPathOfRepository_Call {
	_c.Call.Return(run)
	return _c
}

// GetRemoteBranches provides a mock function with given fields: repo
func (_m *IRepoHelper) GetRemoteBranches(repo *git.Repository) []*plumbing.Reference {
	ret := _m.Called(repo)

	var r0 []*plumbing.Reference
	if rf, ok := ret.Get(0).(func(*git.Repository) []*plumbing.Reference); ok {
		r0 = rf(repo)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*plumbing.Reference)
		}
	}

	return r0
}

// IRepoHelper_GetRemoteBranches_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRemoteBranches'
type IRepoHelper_GetRemoteBranches_Call struct {
	*mock.Call
}

// GetRemoteBranches is a helper method to define mock.On call
//   - repo *git.Repository
func (_e *IRepoHelper_Expecter) GetRemoteBranches(repo interface{}) *IRepoHelper_GetRemoteBranches_Call {
	return &IRepoHelper_GetRemoteBranches_Call{Call: _e.mock.On("GetRemoteBranches", repo)}
}

func (_c *IRepoHelper_GetRemoteBranches_Call) Run(run func(repo *git.Repository)) *IRepoHelper_GetRemoteBranches_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*git.Repository))
	})
	return _c
}

func (_c *IRepoHelper_GetRemoteBranches_Call) Return(_a0 []*plumbing.Reference) *IRepoHelper_GetRemoteBranches_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IRepoHelper_GetRemoteBranches_Call) RunAndReturn(run func(*git.Repository) []*plumbing.Reference) *IRepoHelper_GetRemoteBranches_Call {
	_c.Call.Return(run)
	return _c
}

// GetTemplateTasks provides a mock function with given fields: repo, templates
func (_m *IRepoHelper) GetTemplateTasks(repo *git.Repository, templates []Analyzer.Template) []Analyzer.TemplateTask {
	ret := _m.Called(repo, templates)

	var r0 []Analyzer.TemplateTask
	if rf, ok := ret.Get(0).(func(*git.Repository, []Analyzer.Template) []Analyzer.TemplateTask); ok {
		r0 = rf(repo, templates)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Analyzer.TemplateTask)
		}
	}

	return r0
}

// IRepoHelper_GetTemplateTasks_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTemplateTasks'
type IRepoHelper_GetTemplateTasks_Call struct {
	*mock.Call
}

// GetTemplateTasks is a helper method to define mock.On call
//   - repo *git.Repository
//   - templates []Analyzer.Template
func (_e *IRepoHelper_Expecter) GetTemplateTasks(repo interface{}, templates interface{}) *IRepoHelper_GetTemplateTasks_Call {
	return &IRepoHelper_GetTemplateTasks_Call{Call: _e.mock.On("GetTemplateTasks", repo, templates)}
}

func (_c *IRepoHelper_GetTemplateTasks_Call) Run(run func(repo *git.Repository, templates []Analyzer.Template)) *IRepoHelper_GetTemplateTasks_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*git.Repository), args[1].([]Analyzer.Template))
	})
	return _c
}

func (_c *IRepoHelper_GetTemplateTasks_Call) Return(_a0 []Analyzer.TemplateTask) *IRepoHelper_GetTemplateTasks_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IRepoHelper_GetTemplateTasks_Call) RunAndReturn(run func(*git.Repository, []Analyzer.Template) []Analyzer.TemplateTask) *IRepoHelper_GetTemplateTasks_Call {
	_c.Call.Return(run)
	return _c
}

// ResetRepository provides a mock function with given fields: repo, headCommit
func (_m *IRepoHelper) ResetRepository(repo *git.Repository, headCommit *object.Commit) {
	_m.Called(repo, headCommit)
}

// IRepoHelper_ResetRepository_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ResetRepository'
type IRepoHelper_ResetRepository_Call struct {
	*mock.Call
}

// ResetRepository is a helper method to define mock.On call
//   - repo *git.Repository
//   - headCommit *object.Commit
func (_e *IRepoHelper_Expecter) ResetRepository(repo interface{}, headCommit interface{}) *IRepoHelper_ResetRepository_Call {
	return &IRepoHelper_ResetRepository_Call{Call: _e.mock.On("ResetRepository", repo, headCommit)}
}

func (_c *IRepoHelper_ResetRepository_Call) Run(run func(repo *git.Repository, headCommit *object.Commit)) *IRepoHelper_ResetRepository_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*git.Repository), args[1].(*object.Commit))
	})
	return _c
}

func (_c *IRepoHelper_ResetRepository_Call) Return() *IRepoHelper_ResetRepository_Call {
	_c.Call.Return()
	return _c
}

func (_c *IRepoHelper_ResetRepository_Call) RunAndReturn(run func(*git.Repository, *object.Commit)) *IRepoHelper_ResetRepository_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewIRepoHelper interface {
	mock.TestingT
	Cleanup(func())
}

// NewIRepoHelper creates a new instance of IRepoHelper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIRepoHelper(t mockConstructorTestingTNewIRepoHelper) *IRepoHelper {
	mock := &IRepoHelper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
