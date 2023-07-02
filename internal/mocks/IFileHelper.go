// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	Analyzer "GitAnalyzer/api/Analyzer"

	mock "github.com/stretchr/testify/mock"
)

// IFileHelper is an autogenerated mock type for the IFileHelper type
type IFileHelper struct {
	mock.Mock
}

type IFileHelper_Expecter struct {
	mock *mock.Mock
}

func (_m *IFileHelper) EXPECT() *IFileHelper_Expecter {
	return &IFileHelper_Expecter{mock: &_m.Mock}
}

// FindFilesForCommands provides a mock function with given fields: rootDir, template
func (_m *IFileHelper) FindFilesForCommands(rootDir string, template Analyzer.Template) map[string][]string {
	ret := _m.Called(rootDir, template)

	var r0 map[string][]string
	if rf, ok := ret.Get(0).(func(string, Analyzer.Template) map[string][]string); ok {
		r0 = rf(rootDir, template)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string][]string)
		}
	}

	return r0
}

// IFileHelper_FindFilesForCommands_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindFilesForCommands'
type IFileHelper_FindFilesForCommands_Call struct {
	*mock.Call
}

// FindFilesForCommands is a helper method to define mock.On call
//   - rootDir string
//   - template Analyzer.Template
func (_e *IFileHelper_Expecter) FindFilesForCommands(rootDir interface{}, template interface{}) *IFileHelper_FindFilesForCommands_Call {
	return &IFileHelper_FindFilesForCommands_Call{Call: _e.mock.On("FindFilesForCommands", rootDir, template)}
}

func (_c *IFileHelper_FindFilesForCommands_Call) Run(run func(rootDir string, template Analyzer.Template)) *IFileHelper_FindFilesForCommands_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(Analyzer.Template))
	})
	return _c
}

func (_c *IFileHelper_FindFilesForCommands_Call) Return(_a0 map[string][]string) *IFileHelper_FindFilesForCommands_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IFileHelper_FindFilesForCommands_Call) RunAndReturn(run func(string, Analyzer.Template) map[string][]string) *IFileHelper_FindFilesForCommands_Call {
	_c.Call.Return(run)
	return _c
}

// GenerateUniqueCSV provides a mock function with given fields: templateName
func (_m *IFileHelper) GenerateUniqueCSV(templateName string) {
	_m.Called(templateName)
}

// IFileHelper_GenerateUniqueCSV_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateUniqueCSV'
type IFileHelper_GenerateUniqueCSV_Call struct {
	*mock.Call
}

// GenerateUniqueCSV is a helper method to define mock.On call
//   - templateName string
func (_e *IFileHelper_Expecter) GenerateUniqueCSV(templateName interface{}) *IFileHelper_GenerateUniqueCSV_Call {
	return &IFileHelper_GenerateUniqueCSV_Call{Call: _e.mock.On("GenerateUniqueCSV", templateName)}
}

func (_c *IFileHelper_GenerateUniqueCSV_Call) Run(run func(templateName string)) *IFileHelper_GenerateUniqueCSV_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *IFileHelper_GenerateUniqueCSV_Call) Return() *IFileHelper_GenerateUniqueCSV_Call {
	_c.Call.Return()
	return _c
}

func (_c *IFileHelper_GenerateUniqueCSV_Call) RunAndReturn(run func(string)) *IFileHelper_GenerateUniqueCSV_Call {
	_c.Call.Return(run)
	return _c
}

// GetResultCSVPath provides a mock function with given fields: templateName
func (_m *IFileHelper) GetResultCSVPath(templateName string) string {
	ret := _m.Called(templateName)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(templateName)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// IFileHelper_GetResultCSVPath_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetResultCSVPath'
type IFileHelper_GetResultCSVPath_Call struct {
	*mock.Call
}

// GetResultCSVPath is a helper method to define mock.On call
//   - templateName string
func (_e *IFileHelper_Expecter) GetResultCSVPath(templateName interface{}) *IFileHelper_GetResultCSVPath_Call {
	return &IFileHelper_GetResultCSVPath_Call{Call: _e.mock.On("GetResultCSVPath", templateName)}
}

func (_c *IFileHelper_GetResultCSVPath_Call) Run(run func(templateName string)) *IFileHelper_GetResultCSVPath_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *IFileHelper_GetResultCSVPath_Call) Return(_a0 string) *IFileHelper_GetResultCSVPath_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IFileHelper_GetResultCSVPath_Call) RunAndReturn(run func(string) string) *IFileHelper_GetResultCSVPath_Call {
	_c.Call.Return(run)
	return _c
}

// GetTasks provides a mock function with given fields: urlsCSVPath, checkedCSVPath
func (_m *IFileHelper) GetTasks(urlsCSVPath string, checkedCSVPath string) []Analyzer.Task {
	ret := _m.Called(urlsCSVPath, checkedCSVPath)

	var r0 []Analyzer.Task
	if rf, ok := ret.Get(0).(func(string, string) []Analyzer.Task); ok {
		r0 = rf(urlsCSVPath, checkedCSVPath)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Analyzer.Task)
		}
	}

	return r0
}

// IFileHelper_GetTasks_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTasks'
type IFileHelper_GetTasks_Call struct {
	*mock.Call
}

// GetTasks is a helper method to define mock.On call
//   - urlsCSVPath string
//   - checkedCSVPath string
func (_e *IFileHelper_Expecter) GetTasks(urlsCSVPath interface{}, checkedCSVPath interface{}) *IFileHelper_GetTasks_Call {
	return &IFileHelper_GetTasks_Call{Call: _e.mock.On("GetTasks", urlsCSVPath, checkedCSVPath)}
}

func (_c *IFileHelper_GetTasks_Call) Run(run func(urlsCSVPath string, checkedCSVPath string)) *IFileHelper_GetTasks_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *IFileHelper_GetTasks_Call) Return(tasks []Analyzer.Task) *IFileHelper_GetTasks_Call {
	_c.Call.Return(tasks)
	return _c
}

func (_c *IFileHelper_GetTasks_Call) RunAndReturn(run func(string, string) []Analyzer.Task) *IFileHelper_GetTasks_Call {
	_c.Call.Return(run)
	return _c
}

// MarshalMultipleResults provides a mock function with given fields: results
func (_m *IFileHelper) MarshalMultipleResults(results []Analyzer.Result) {
	_m.Called(results)
}

// IFileHelper_MarshalMultipleResults_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MarshalMultipleResults'
type IFileHelper_MarshalMultipleResults_Call struct {
	*mock.Call
}

// MarshalMultipleResults is a helper method to define mock.On call
//   - results []Analyzer.Result
func (_e *IFileHelper_Expecter) MarshalMultipleResults(results interface{}) *IFileHelper_MarshalMultipleResults_Call {
	return &IFileHelper_MarshalMultipleResults_Call{Call: _e.mock.On("MarshalMultipleResults", results)}
}

func (_c *IFileHelper_MarshalMultipleResults_Call) Run(run func(results []Analyzer.Result)) *IFileHelper_MarshalMultipleResults_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]Analyzer.Result))
	})
	return _c
}

func (_c *IFileHelper_MarshalMultipleResults_Call) Return() *IFileHelper_MarshalMultipleResults_Call {
	_c.Call.Return()
	return _c
}

func (_c *IFileHelper_MarshalMultipleResults_Call) RunAndReturn(run func([]Analyzer.Result)) *IFileHelper_MarshalMultipleResults_Call {
	_c.Call.Return(run)
	return _c
}

// MarshalSingleResult provides a mock function with given fields: result
func (_m *IFileHelper) MarshalSingleResult(result Analyzer.Result) {
	_m.Called(result)
}

// IFileHelper_MarshalSingleResult_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MarshalSingleResult'
type IFileHelper_MarshalSingleResult_Call struct {
	*mock.Call
}

// MarshalSingleResult is a helper method to define mock.On call
//   - result Analyzer.Result
func (_e *IFileHelper_Expecter) MarshalSingleResult(result interface{}) *IFileHelper_MarshalSingleResult_Call {
	return &IFileHelper_MarshalSingleResult_Call{Call: _e.mock.On("MarshalSingleResult", result)}
}

func (_c *IFileHelper_MarshalSingleResult_Call) Run(run func(result Analyzer.Result)) *IFileHelper_MarshalSingleResult_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(Analyzer.Result))
	})
	return _c
}

func (_c *IFileHelper_MarshalSingleResult_Call) Return() *IFileHelper_MarshalSingleResult_Call {
	_c.Call.Return()
	return _c
}

func (_c *IFileHelper_MarshalSingleResult_Call) RunAndReturn(run func(Analyzer.Result)) *IFileHelper_MarshalSingleResult_Call {
	_c.Call.Return(run)
	return _c
}

// MarshalStat provides a mock function with given fields: stat
func (_m *IFileHelper) MarshalStat(stat Analyzer.Stat) {
	_m.Called(stat)
}

// IFileHelper_MarshalStat_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MarshalStat'
type IFileHelper_MarshalStat_Call struct {
	*mock.Call
}

// MarshalStat is a helper method to define mock.On call
//   - stat Analyzer.Stat
func (_e *IFileHelper_Expecter) MarshalStat(stat interface{}) *IFileHelper_MarshalStat_Call {
	return &IFileHelper_MarshalStat_Call{Call: _e.mock.On("MarshalStat", stat)}
}

func (_c *IFileHelper_MarshalStat_Call) Run(run func(stat Analyzer.Stat)) *IFileHelper_MarshalStat_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(Analyzer.Stat))
	})
	return _c
}

func (_c *IFileHelper_MarshalStat_Call) Return() *IFileHelper_MarshalStat_Call {
	_c.Call.Return()
	return _c
}

func (_c *IFileHelper_MarshalStat_Call) RunAndReturn(run func(Analyzer.Stat)) *IFileHelper_MarshalStat_Call {
	_c.Call.Return(run)
	return _c
}

// PrepareResultsFolder provides a mock function with given fields: path
func (_m *IFileHelper) PrepareResultsFolder(path string) {
	_m.Called(path)
}

// IFileHelper_PrepareResultsFolder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PrepareResultsFolder'
type IFileHelper_PrepareResultsFolder_Call struct {
	*mock.Call
}

// PrepareResultsFolder is a helper method to define mock.On call
//   - path string
func (_e *IFileHelper_Expecter) PrepareResultsFolder(path interface{}) *IFileHelper_PrepareResultsFolder_Call {
	return &IFileHelper_PrepareResultsFolder_Call{Call: _e.mock.On("PrepareResultsFolder", path)}
}

func (_c *IFileHelper_PrepareResultsFolder_Call) Run(run func(path string)) *IFileHelper_PrepareResultsFolder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *IFileHelper_PrepareResultsFolder_Call) Return() *IFileHelper_PrepareResultsFolder_Call {
	_c.Call.Return()
	return _c
}

func (_c *IFileHelper_PrepareResultsFolder_Call) RunAndReturn(run func(string)) *IFileHelper_PrepareResultsFolder_Call {
	_c.Call.Return(run)
	return _c
}

// SearchFilesByRegex provides a mock function with given fields: rootDir, template
func (_m *IFileHelper) SearchFilesByRegex(rootDir string, template Analyzer.Template) []Analyzer.Result {
	ret := _m.Called(rootDir, template)

	var r0 []Analyzer.Result
	if rf, ok := ret.Get(0).(func(string, Analyzer.Template) []Analyzer.Result); ok {
		r0 = rf(rootDir, template)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Analyzer.Result)
		}
	}

	return r0
}

// IFileHelper_SearchFilesByRegex_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SearchFilesByRegex'
type IFileHelper_SearchFilesByRegex_Call struct {
	*mock.Call
}

// SearchFilesByRegex is a helper method to define mock.On call
//   - rootDir string
//   - template Analyzer.Template
func (_e *IFileHelper_Expecter) SearchFilesByRegex(rootDir interface{}, template interface{}) *IFileHelper_SearchFilesByRegex_Call {
	return &IFileHelper_SearchFilesByRegex_Call{Call: _e.mock.On("SearchFilesByRegex", rootDir, template)}
}

func (_c *IFileHelper_SearchFilesByRegex_Call) Run(run func(rootDir string, template Analyzer.Template)) *IFileHelper_SearchFilesByRegex_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(Analyzer.Template))
	})
	return _c
}

func (_c *IFileHelper_SearchFilesByRegex_Call) Return(result []Analyzer.Result) *IFileHelper_SearchFilesByRegex_Call {
	_c.Call.Return(result)
	return _c
}

func (_c *IFileHelper_SearchFilesByRegex_Call) RunAndReturn(run func(string, Analyzer.Template) []Analyzer.Result) *IFileHelper_SearchFilesByRegex_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewIFileHelper interface {
	mock.TestingT
	Cleanup(func())
}

// NewIFileHelper creates a new instance of IFileHelper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIFileHelper(t mockConstructorTestingTNewIFileHelper) *IFileHelper {
	mock := &IFileHelper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
