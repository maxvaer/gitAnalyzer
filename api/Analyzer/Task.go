// Package Analyzer contains all structural components of the application.
package Analyzer

// The Task struct is used to define a working task for the gitAnalyzer
type Task struct {
	// Git URL of the repository to scan
	URL string `csv:"url"`
	// Language of the git repository (optional)
	Language string `csv:"language"`
	// Current State of the task, can be one of: Queued, Cloning, Failed, Running, Finished
	State string `csv:"-"`
	// The Results found for the repository
	Results []Result `csv:"-"`
	// The time it took to run the task
	ElapsedTime string `csv:"-"`
}
