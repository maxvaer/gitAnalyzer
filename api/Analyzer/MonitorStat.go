// Package Analyzer contains all structural components of the application.
package Analyzer

// The MonitorStat struct is used pass information to the webUI frontend
type MonitorStat struct {
	// The number of total tasks to scan
	NumberOfTasks int32 `json:"numberOfTasks,omitempty"`
	// The total number of failed scans
	FailedScans int32 `json:"failedScans,omitempty"`
	// The number of scans which are currently queued
	QueuedScans int32 `json:"queuedScans,omitempty"`
	// The count of currently running scans
	RunningScans int32 `json:"runningScans,omitempty"`
	// The overall number of finished scans
	FinishedScans int32 `json:"finishedScans,omitempty"`
	// The total number of results written to files
	ResultsFound int32 `json:"resultsFound,omitempty"`
	// The count of currently cloning repositories
	Cloning int32 `json:"cloning,omitempty"`
}
