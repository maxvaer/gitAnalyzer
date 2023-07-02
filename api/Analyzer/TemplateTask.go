// Package Analyzer contains all structural components of the application.
package Analyzer

// The TemplateTask is used to store the information
// which template are run on a given commit hash
type TemplateTask struct {
	CommitHash string
	// The templates which will be run for the commit hash
	Templates []Template
}
